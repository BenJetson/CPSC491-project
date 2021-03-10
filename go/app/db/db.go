package db

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register PostgreSQL drivers with sql.DB
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// Config specifies how to connect to the database server.
type Config struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
	DisableTLS bool
}

func (cfg *Config) connectString() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("host='%s' ", cfg.Host))
	buf.WriteString(fmt.Sprintf("port=%d ", cfg.Port))
	buf.WriteString(fmt.Sprintf("user='%s' ", cfg.Username))
	buf.WriteString(fmt.Sprintf("password='%s' ", cfg.Password))
	buf.WriteString(fmt.Sprintf("database='%s' ", cfg.Database))

	tls := "verify-full"
	if cfg.DisableTLS {
		tls = "disable"
	}

	buf.WriteString(fmt.Sprintf("sslmode='%s' ", tls))

	return buf.String()
}

// NewConfigFromEnv attempts to construct a new Config using data from
// environment variables.
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

	// Optional flag to disable TLS for database connections.
	// Necessary for local development databases running in Docker.
	// DO NOT use an unencrypted database connection in production!
	disableTLS := os.Getenv("DB_DANGER_DISABLE_TLS")
	if disableTLS == "accept danger and use unencrypted connection" {
		cfg.DisableTLS = true
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
	handle, err := sqlx.Connect("postgres", cfg.connectString())
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to database")
	}

	// SQLx connection will return nil error when ping was successful, so we
	// can guarantee good database connection here.
	logger.Info("Connected to the database.")

	return &database{
		DB:     handle,
		logger: logger,
		cfg:    cfg,
	}, nil
}
