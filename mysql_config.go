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

func SetupCustomTLS(cfg *mysql.Config, caCertPath string) error {
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(caCertPath)

	if err != nil {
		return err
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return errors.New("Failed to append PEM.")
	}

	key := "custom"

	err = mysql.RegisterTLSConfig(key, &tls.Config{
		RootCAs: rootCertPool,
	})

	if err != nil {
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
