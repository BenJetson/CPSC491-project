package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BenJetson/CPSC491-project/go/app/mock"
)

func TestPanicRecoveryMiddleware(t *testing.T) {
	db := &mock.DB{}
	api, _, hook := newTestAPI(t, db, nil)

	panickingHandler := func(_ http.ResponseWriter, _ *http.Request) {
		panic("aaaaaack!!!")
	}

	require.Panics(t, func() { panickingHandler(nil, nil) },
		"test will not work if our handler does not panic")

	// Set an extra route to force a panic.
	api.router.Path("/panic").Methods("GET").HandlerFunc(panickingHandler)

	r := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	require.NotPanics(t, func() {
		api.router.ServeHTTP(w, r)
	}, "handler should not panic; should be caught by middleware")

	// One log for panic, one log for error.
	assert.Len(t, hook.Entries, 2)

	panicEntry := hook.Entries[0]
	assert.Equal(t, logrus.ErrorLevel, panicEntry.Level)
	assert.Equal(t, "panic when handling api request; recovering",
		panicEntry.Message)
	assert.Contains(t, panicEntry.Data, "stack")
	assert.Contains(t, panicEntry.Data, logrus.ErrorKey)

	errorEntry := hook.Entries[1]
	assert.Equal(t, logrus.ErrorLevel, errorEntry.Level)
	assert.Equal(t, "encountered error when handling api request",
		errorEntry.Message)
	assert.Contains(t, errorEntry.Data, "details")
	assert.Contains(t, errorEntry.Data, logrus.ErrorKey)
}
