package qb

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	RecordPeriod = 1 * time.Second
)

type Agent struct {
	id          int
	mysqlConfig *MysqlConfig
	db          DB
	taskOps     *TaskOpts
	stmts       []string
}

func newAgent(id int, myCfg *MysqlConfig, taskOps *TaskOpts, stmts []string) *Agent {
	return &Agent{
		id:          id,
		mysqlConfig: myCfg,
		taskOps:     taskOps,
		stmts:       stmts,
	}
}

func (agent *Agent) prepare(maxIdleConns int) error {
	db, err := agent.mysqlConfig.openAndPing(maxIdleConns)

	if err != nil {
		dsn := agent.mysqlConfig.FormatDSN()
		return fmt.Errorf("failed to open/ping DB (agent id=%d, dsn=%s): %w", agent.id, dsn, err)
	}

	_, err = db.Exec("SET autocommit = 0")

	if err != nil {
		return fmt.Errorf("disable autocommit error: %w", err)
	}

	agent.db = db

	return nil
}

func (agent *Agent) run(ctx context.Context, recorder *Recorder, token string) error {
	_, err := agent.db.Exec(fmt.Sprintf("SELECT 'agent(%d) start: token=%s'", agent.id, token))

	if err != nil {
		return fmt.Errorf("failed to execute start query (agent id=%d): %w", agent.id, err)
	}

	recordTick := time.NewTicker(RecordPeriod)
	defer recordTick.Stop()
	recDps := []recorderDataPoint{}

	err = loopWithThrottle(agent.stmts, agent.taskOps.Rate, func(i int, q string) (bool, error) {
		select {
		case <-ctx.Done():
			return false, nil
		case <-recordTick.C:
			recorder.add(recDps)
			recDps = recDps[:0]
		default:
			// Nothing to do
		}

		rt, err := agent.query(ctx, q)

		if err != nil {
			return false, fmt.Errorf("execute query error (query=%s): %w", q, err)
		}

		recDps = append(recDps, recorderDataPoint{
			timestamp: time.Now(),
			resTime:   rt,
		})

		return true, nil
	})

	if err != nil {
		return fmt.Errorf("failed to transact (agent id=%d): %w", agent.id, err)
	}

	_, err = agent.db.Exec(fmt.Sprintf("SELECT 'agent(%d) end: token=%s'", agent.id, token))

	if err != nil {
		return fmt.Errorf("failed to execute exit query (agent id=%d): %w", agent.id, err)
	}

	return nil
}

func (agent *Agent) close() error {
	err := agent.db.Close()

	if err != nil {
		return fmt.Errorf("failed to close DB (agent id=%d): %w", agent.id, err)
	}

	return nil
}

func (agent *Agent) query(ctx context.Context, q string) (time.Duration, error) {
	start := time.Now()
	_, err := agent.db.ExecContext(ctx, q)
	end := time.Now()

	if err != nil && !errors.Is(err, context.Canceled) {
		return 0, err
	}

	return end.Sub(start), nil
}
