package db

import (
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// Config specifies how to connect to the database server.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func NewConfigFromEnv() (cfg Config, err error) {
	// First, let us check to ensure the environment variables are set.
	for _, key := range []string{"HOST", "PORT", "USER", "PASS", "DATABASE"} {
		key = "DB_" + key
		if len(os.Getenv(key)) < 1 {
			err = errors.Errorf("must set %s", key)
			return
		}
	}

	// Variables are set. Store values in configuration.
	cfg.Host = os.Getenv("DB_HOST")
	cfg.Username = os.Getenv("DB_USER")
	cfg.Password = os.Getenv("DB_PASS")
	cfg.Database = os.Getenv("DB_DATABASE")

	if cfg.Port, err = strconv.Atoi(os.Getenv("DB_PORT")); err != nil {
		err = errors.New("DB_PORT must be an integer")
		return
	}

	return
}

type database struct {
	*sqlx.DB
	logger *logrus.Logger
	cfg    Config
}

// NewDataStore creates a new database handle.
func NewDataStore(logger *logrus.Logger, cfg Config) (app.DataStore, error) {
	// FIXME must initialize *sqlx.DB
	return &database{
		logger: logger,
		cfg:    cfg,
	}, nil
}
