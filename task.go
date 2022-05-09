package qb

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

const (
	ProgressReportPeriod      = 1
	NBranches                 = 1
	NTellers                  = 10
	NAccounts                 = 100000
	InsertChunkSize           = 50000
	GeneratingDataConcurrency = 30
)

var initCreateStmts = []string{
	"CREATE TABLE qb_history (id int PRIMARY KEY NOT NULL AUTO_INCREMENT, tid int, bid int, aid bigint, delta int, mtime timestamp, filler char(22))",
	"CREATE TABLE qb_tellers (tid int PRIMARY KEY NOT NULL, bid int, tbalance int, filler char(84))",
	"CREATE TABLE qb_accounts (aid bigint PRIMARY KEY NOT NULL, bid int, abalance int, filler char(84))",
	"CREATE TABLE qb_branches (bid int PRIMARY KEY NOT NULL, bbalance int, filler char(88))",
}

var initAnalyzeStmts = []string{
	"ANALYZE TABLE qb_history",
	"ANALYZE TABLE qb_tellers",
	"ANALYZE TABLE qb_accounts",
	"ANALYZE TABLE qb_branches",
}

var initInsertStmts = map[string]func(n int) []string{
	"INSERT INTO qb_branches (bid, bbalance) VALUES ": func(scale int) []string {
		values := make([]string, 0, NBranches*scale)

		for i := 1; i <= NBranches*scale; i++ {
			values = append(values, fmt.Sprintf("(%d, 0)", i))
		}

		return values
	},
	"INSERT INTO qb_tellers (tid, bid, tbalance) VALUES ": func(scale int) []string {
		values := make([]string, 0, NTellers*scale)

		for i := 1; i <= NTellers*scale; i++ {
			values = append(values, fmt.Sprintf("(%d, (%d - 1) / %d + 1, 0)", i, i, NTellers))
		}

		return values
	},
	"INSERT INTO qb_accounts (aid, bid, abalance, filler) VALUES ": func(scale int) []string {
		values := make([]string, 0, NAccounts*scale)

		for i := 1; i <= NAccounts*scale; i++ {
			values = append(values, fmt.Sprintf("(%d, (%d - 1) / %d + 1, 0, '')", i, i, NAccounts))
		}

		return values
	},
}

type TaskOpts struct {
	MysqlConfig     *MysqlConfig `json:"-"`
	NAgents         int
	Time            time.Duration `json:"-"`
	Rate            int
	TransactionType string
	Scale           int `json:"-"`
	Engine          string
	OnlyPrint       bool `json:"-"`
	NoProgress      bool `json:"-"`
}

type Task struct {
	*TaskOpts
	agents   []*Agent
	recOpts  *RecorderOpts
	stmtSize int
}

func NewTask(taskOpts *TaskOpts, recOpts *RecorderOpts) (*Task, error) {
	agents := make([]*Agent, taskOpts.NAgents)

	stmts, err := NewScript(taskOpts.TransactionType, taskOpts.NAgents)

	if err != nil {
		return nil, fmt.Errorf("failed to build script: %w", err)
	}

	for i := 0; i < taskOpts.NAgents; i++ {
		agents[i] = newAgent(i, taskOpts.MysqlConfig, taskOpts, stmts)
	}

	return &Task{
		TaskOpts: taskOpts,
		agents:   agents,
		recOpts:  recOpts,
		stmtSize: len(stmts),
	}, nil
}

func (task *Task) Prepare() error {
	for _, agent := range task.agents {
		if err := agent.prepare(task.NAgents); err != nil {
			return fmt.Errorf("failed to prepare Agent: %w", err)
		}
	}

	return nil
}

func (task *Task) Initialize() error {
	// Temporarily empty the DB name
	orgDBName := task.MysqlConfig.DBName
	task.MysqlConfig.DBName = ""

	db, err := task.MysqlConfig.openAndPing(1)

	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	if task.Engine != "" {
		_, err = db.Exec(fmt.Sprintf("SET default_storage_engine = %s", task.Engine))

		if err != nil {
			return fmt.Errorf("set default_storage_engine error: %w", err)
		}
	}

	defer db.Close()
	task.MysqlConfig.DBName = orgDBName

	log.Println("dropping old database...")
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", task.MysqlConfig.DBName))

	if err != nil {
		return fmt.Errorf("drop database error: %w", err)
	}

	log.Println("creating database...")
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s`", task.MysqlConfig.DBName))

	if err != nil {
		return fmt.Errorf("create database error: %w", err)
	}

	_, err = db.Exec(fmt.Sprintf("USE `%s`", task.MysqlConfig.DBName))

	if err != nil {
		return fmt.Errorf("use database error: %w", err)
	}

	log.Println("creating tables...")
	for _, stmt := range initCreateStmts {
		_, err := db.Exec(stmt)

		if err != nil {
			return fmt.Errorf("create table error (query=%s): %w", stmt, err)
		}
	}

	log.Println("generating data...")
	ctxWithoutCancel := context.Background()
	ctx, cancel := context.WithCancel(ctxWithoutCancel)
	eg := task.setupTables(ctx)
	task.trapSigint(ctx, cancel, eg)
	err = eg.Wait()
	cancel()

	if err != nil {
		return fmt.Errorf("generating data error: %w", err)
	}

	log.Println("analyzing tables...")
	for _, stmt := range initAnalyzeStmts {
		_, err := db.Exec(stmt)

		if err != nil {
			return fmt.Errorf("analyze table error (query=%s): %w", stmt, err)
		}
	}

	return nil
}

func (task *Task) setupTables(ctx context.Context) *errgroup.Group {
	sem := make(chan struct{}, GeneratingDataConcurrency)
	eg, ctx := errgroup.WithContext(ctx)

	for prefix, bldr := range initInsertStmts {
		values := bldr(task.Scale)

		for i := 0; i < len(values); i += InsertChunkSize {
			to := i + InsertChunkSize

			if len(values) < to {
				to = len(values)
			}

			stmt := prefix + strings.Join(values[i:to], ",")
			sem <- struct{}{}

			eg.Go(func() error {
				defer func() { <-sem }()
				db, err := task.MysqlConfig.openAndPing(1)

				if err != nil {
					return fmt.Errorf("connection error: %w", err)
				}

				defer db.Close()

				select {
				case <-ctx.Done():
					return nil
				default:
				}

				_, err = db.Exec(stmt)

				if err != nil {
					return fmt.Errorf("insert data error (query=%s): %w", stmt, err)
				}

				return nil
			})
		}
	}

	return eg
}

func (task *Task) Run() (*Recorder, error) {
	uuid, _ := uuid.NewRandom()
	token := uuid.String()
	rec := newRecorder(task.recOpts, task.TaskOpts, token, task.stmtSize)

	defer func() {
		rec.close()

		for _, agent := range task.agents {
			err := agent.close()

			if err != nil {
				fmt.Fprintf(os.Stderr, "[WARN] failed to close agent: %s", err)
			}
		}
	}()

	eg, ctxWithoutCancel := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctxWithoutCancel)
	progressTick := time.NewTicker(ProgressReportPeriod * time.Second)
	rec.start(task.NAgents * 3)
	var numTermAgents int32

	// Variables for progress line
	taskStart := time.Now()
	prevExecCnt := 0

	// Run agents
	for _, v := range task.agents {
		agent := v
		eg.Go(func() error {
			err := agent.run(ctx, rec, token)
			atomic.AddInt32(&numTermAgents, 1)
			return err
		})
	}

	// Periodic report progress
	go func() {
	LOOP:
		for {
			select {
			case <-ctx.Done():
				progressTick.Stop()
				break LOOP
			case <-progressTick.C:
				if !task.NoProgress && !task.OnlyPrint {
					execCnt := rec.Count()
					termAgentCnt := int(atomic.LoadInt32(&numTermAgents))
					task.printProgress(execCnt, prevExecCnt, taskStart, termAgentCnt)
					prevExecCnt = execCnt
				}
			}
		}
	}()

	// Time-out processing
	// NOTE: If it is zero, it will not time out
	if task.Time > 0 {
		go func() {
			select {
			case <-ctx.Done():
				// Nothing to do
			case <-time.After(task.Time):
				cancel()
			}
		}()
	}

	task.trapSigint(ctx, cancel, eg)
	err := eg.Wait()
	cancel()

	// Clear progress line
	if !task.NoProgress || !task.OnlyPrint {
		fmt.Fprintf(os.Stderr, "\r\n\n")
	}

	if err != nil {
		return nil, fmt.Errorf("error during agent running: %w", err)
	}

	return rec, nil
}

func (task *Task) printProgress(execCnt int, prevExecCnt int, taskStart time.Time, numTermAgents int) {
	qps := float64(execCnt-prevExecCnt) / ProgressReportPeriod
	elapsedTime := time.Since(taskStart)
	numRunAgents := task.NAgents - int(numTermAgents)
	termWidth, _, err := term.GetSize(0)

	if err != nil {
		panic("Failed to get terminal width: " + err.Error())
	}

	elapsedTime = elapsedTime.Round(time.Second)
	min := elapsedTime / time.Minute
	sec := (elapsedTime - min*time.Minute) / time.Second
	progressLine := fmt.Sprintf("%02d:%02d | %d agents / run %d queries (%.0f tps)", min, sec, numRunAgents, execCnt, qps/float64(task.stmtSize))
	fmt.Fprintf(os.Stderr, "\r%-*s", termWidth, progressLine)
}

func (task *Task) trapSigint(ctx context.Context, cancel context.CancelFunc, eg *errgroup.Group) {
	// SIGINT
	sgnlCh := make(chan os.Signal, 1)
	signal.Notify(sgnlCh, os.Interrupt)

	go func() {
		select {
		case <-ctx.Done():
			// Nothing to do
		case <-sgnlCh:
			cancel()
			_ = eg.Wait()
			os.Exit(130)
		}
	}()
}
