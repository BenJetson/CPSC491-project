package api

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Config specifies the configuration for an API server instance.
type Config struct {
	Port int
	Tier string
}

// NewConfigFromEnv attempts to construct a new Config using data from
// environment variables.
func NewConfigFromEnv() (c Config, err error) {
	c.Tier = os.Getenv("TIER")
	if len(c.Tier) < 1 {
		err = errors.New("must set TIER")
		return
	}

	port := os.Getenv("PORT")
	if len(port) < 1 {
		err = errors.New("must set PORT")
		return
	} else if c.Port, err = strconv.Atoi(port); err != nil {
		err = errors.New("PORT must be an integer")
		return
	}

	return
}
