package qb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"os"

	"github.com/go-sql-driver/mysql"
)

type CertOptions struct {
	CaCertPath     string
	ClientCertPath string
	ClientKeyPath  string
}

func SetupCustomTLS(cfg *mysql.Config, certOptions *CertOptions) error {
	tlsConfig := &tls.Config{}

	if certOptions.CaCertPath != "" {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(certOptions.CaCertPath)

		if err != nil {
			return err
		}

		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return errors.New("failed to append PEM")
		}

		tlsConfig.RootCAs = rootCertPool
	}

	if certOptions.ClientCertPath != "" && certOptions.ClientKeyPath != "" {
		certs, err := tls.LoadX509KeyPair(certOptions.ClientCertPath, certOptions.ClientKeyPath)

		if err != nil {
			return err
		}

		tlsConfig.Certificates = []tls.Certificate{certs}
	}

	key := "custom"
	if err := mysql.RegisterTLSConfig(key, tlsConfig); err != nil {
		return err
	}

	cfg.TLSConfig = key

	return nil
}

type MysqlConfig struct {
	*mysql.Config
	OnlyPrint bool
}

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

func (myCfg *MysqlConfig) openAndPing(maxIdleConns int) (DB, error) {
	if myCfg.OnlyPrint {
		return &NullDB{}, nil
	}

	dsn := myCfg.FormatDSN()
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(0)

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)

	return db, nil
}
