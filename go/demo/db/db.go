package db

import (
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/BenJetson/go-api-demo/go/demo"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type database struct {
	*sqlx.DB
	logger *logrus.Entry
	cfg    Config
}

func NewDataStore(logger *logrus.Entry, cfg Config) (demo.DataStore, error) {
	// FIXME must initialize *sqlx.DB
	return &database{
		logger: logger,
		cfg:    cfg,
	}, nil
}

func NewDataStoreFromEnv(logger *logrus.Entry) (demo.DataStore, error) {
	cfg := Config{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_DATABASE"),
	}

	var err error
	if cfg.Port, err = strconv.Atoi(os.Getenv("DB_PORT")); err != nil {
		return nil, errors.Wrap(err, "failed to parse DB port from env")
	}

	return NewDataStore(logger, cfg)
}
