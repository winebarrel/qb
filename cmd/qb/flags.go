package main

import (
	"crypto/tls"
	"crypto/x509"
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
	flaggy.String(&caCertPath, "", "ca-cert", "Path to ca cert. Requires 'tls=custom' in DSN")
	var clientCertPath string
	flaggy.String(&clientCertPath, "", "client-cert", "Path to client certificate. Requires 'tls=custom' in DSN.")
	var clientKeyPath string
	flaggy.String(&clientKeyPath, "", "client-key", "Path to client key. Requires 'tls=custom' in DSN.")
	flaggy.Parse()

	if len(os.Args) <= 1 {
		flaggy.ShowHelpAndExit("")
	}

	// Custom TLS Configuration
	tlsConfig := &tls.Config{}
	var customTLS = false
	if clientCertPath != "" && clientKeyPath != "" {
		customTLS = true
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)

		if err != nil {
			printErrorAndExit("could not load key pair: " + err.Error())
		}

		clientCert = append(clientCert, certs)
		tlsConfig.Certificates = clientCert
	}

	if caCertPath != "" {
		customTLS = true
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(caCertPath)

		if err != nil {
			printErrorAndExit("could not read ca cert: " + err.Error())
		}

		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			printErrorAndExit("failed to append PEM")
		}

		tlsConfig.RootCAs = rootCertPool
	}

	if customTLS {
		err := mysql.RegisterTLSConfig("custom", tlsConfig)

		if err != nil {
			printErrorAndExit("failed to register custom tls config: " + err.Error())
		}
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
