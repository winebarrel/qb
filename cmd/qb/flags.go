package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/integrii/flaggy"
	"github.com/winebarrel/qb"
)

var version string

const (
	DefaultTime            = 60
	DefaultDBName          = "qb"
	DefaultNAgents         = 1
	DefaultScale           = 1
	DefaultTransactionType = "tpcb-like"
)

type Flags struct {
	qb.TaskOpts
	qb.RecorderOpts
	initialize bool
}

func parseFlags() (flags *Flags) {
	flaggy.SetVersion(version)
	flaggy.SetDescription("MySQL benchmarking tool using TCP-B(same as pgbench).")
	flags = &Flags{
		TaskOpts: qb.TaskOpts{
			NAgents:         DefaultNAgents,
			TransactionType: DefaultTransactionType,
			Scale:           DefaultScale,
		},
	}
	var dsn string
	flaggy.String(&dsn, "d", "dsn", "Data Source Name, see https://github.com/go-sql-driver/mysql#examples.")
	flaggy.Bool(&flags.initialize, "i", "initialize", "Invokes initialization mode.")
	flaggy.Int(&flags.NAgents, "n", "nagents", "Number of agents.")
	argTime := DefaultTime
	flaggy.Int(&argTime, "t", "time", "Test run time (sec). Zero is infinity.")
	flaggy.Int(&flags.Rate, "r", "rate", "Rate limit for each agent (qps). Zero is unlimited.")
	flaggy.String(&flags.TransactionType, "T", "type", fmt.Sprintf("Transaction type (%s).", strings.Join(qb.ScriptNames(), ",")))
	flaggy.Int(&flags.Scale, "s", "scale", "Scaling factor.")
	flaggy.String(&flags.Engine, "e", "engine", "Engine of the table to be created.")
	hinterval := "0"
	flaggy.String(&hinterval, "", "hinterval", "Histogram interval, e.g. '100ms'.")
	flaggy.Bool(&flags.OnlyPrint, "", "only-print", "Just print SQL without connecting to DB.")
	flaggy.Bool(&flags.NoProgress, "", "no-progress", "Do not show progress.")
	var caCertPath string
	flaggy.String(&caCertPath, "c", "ca-cert", "path to ca cert")
	flaggy.Parse()

	if len(os.Args) <= 1 {
		flaggy.ShowHelpAndExit("")
	}

	// DSN
	if dsn == "" {
		printErrorAndExit("'--dsn(-d)' is required")
	}

	myCfg, err := mysql.ParseDSN(dsn)

	if err != nil {
		printErrorAndExit("DSN parsing error: " + err.Error())
	}

	flags.DSN = dsn

	if myCfg.DBName == "" {
		myCfg.DBName = DefaultDBName
	}

	// Custom CA Cert Path
	if caCertPath != "" {
		err := qb.SetupCustomTLS(myCfg, caCertPath)

		if err != nil {
			printErrorAndExit("Failed to setup custom TLS: " + err.Error())
		}
	}

	flags.MysqlConfig = &qb.MysqlConfig{
		Config:    myCfg,
		OnlyPrint: flags.OnlyPrint,
	}

	// NAgents
	if flags.NAgents < 1 {
		printErrorAndExit("'--nagents(-n)' must be >= 1")
	}

	// Time
	if argTime < 0 {
		printErrorAndExit("'--time(-t)' must be >= 0")
	}

	flags.Time = time.Duration(argTime) * time.Second

	// Rate
	if flags.Rate < 0 {
		printErrorAndExit("'--rate(-r)' must be >= 0")
	}

	// Scaling Factor
	if flags.Scale < 1 {
		printErrorAndExit("'--scale(-1)' must be >= 1")
	}

	// HInterval
	if hi, err := time.ParseDuration(hinterval); err != nil {
		printErrorAndExit("failed to parse hinterval: " + err.Error())
	} else {
		flags.HInterval = hi
	}

	return
}

func printErrorAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
