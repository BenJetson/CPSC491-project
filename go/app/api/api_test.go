package api

import (
	"testing"

	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func newTestAPI(
	t *testing.T,
	db app.DataStore,
	cv app.CommerceVendor,
) (*Server, *logrus.Logger, *logtest.Hook) {

	logger, hook := logtest.NewNullLogger()

	api, err := NewServer(logger, db, cv, Config{
		Tier: TierLocal,
		Port: 8080,
	})
	require.NoError(t, err, "failed to instantiate test api server")

	return api, logger, hook
}
