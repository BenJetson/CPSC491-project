package api

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Tier represents a different instance of this application.
type Tier string

const (
	// TierProduction represents the production tier of this app running
	// at app.teamxiv.space.
	TierProduction Tier = "prod"
	// TierDevelopment represents the development tier of this app running
	// at dev.teamxiv.space.
	TierDevelopment Tier = "dev"
	// TierLocal represents the local tier of this app running on your local
	// computer inside of Docker.
	TierLocal Tier = "local"
)

// Config specifies the configuration for an API server instance.
type Config struct {
	Port int
	Tier Tier
}

// NewConfigFromEnv attempts to construct a new Config using data from
// environment variables.
func NewConfigFromEnv() (c Config, err error) {
	c.Tier = Tier(os.Getenv("TIER"))
	if len(c.Tier) < 1 {
		err = errors.New("must set TIER")
		return
	}

	validTier := false
	for _, t := range []Tier{TierProduction, TierDevelopment, TierLocal} {
		if c.Tier == t {
			validTier = true
			break
		}
	}

	if !validTier {
		err = errors.New("unknown value for TIER")
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
