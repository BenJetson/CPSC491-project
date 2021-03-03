package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// Server is a wrapper for http.Server that exposes our app's endpoints.
type Server struct {
	config Config
	db     app.DataStore
	logger *logrus.Logger
	httpd  *http.Server
}

// NewServer creates a new Server given a logger, data store, and configuration.
func NewServer(logger *logrus.Logger, db app.DataStore,
	cfg Config) (*Server, error) {

	if logger == nil {
		return nil, errors.New("must specify a logger for the server")
	}

	router := mux.NewRouter()

	httpd := &http.Server{
		Addr:     fmt.Sprintf(":%d", cfg.Port),
		ErrorLog: log.New(logger.WriterLevel(logrus.ErrorLevel), "", 0),
		Handler:  router,
	}

	svr := &Server{
		httpd:  httpd,
		config: cfg,
		db:     db,
		logger: logger,
	}

	// Register global middleware.
	router.Use(svr.authContextMiddleware)

	// Define routes.
	router.Path("/login").Methods("POST").HandlerFunc(svr.handleLogin)
	router.Path("/logout").Methods("POST").HandlerFunc(svr.handleLogout)
	router.Path("/")

	return svr, nil
}

// Start will call ListenAndServe on the internal http.Server and log the
// current configuration.
func (svr *Server) Start() error {
	svr.logger.Infof(
		"Starting API server for tier %s on port %d.\n",
		svr.config.Tier,
		svr.config.Port,
	)

	return svr.httpd.ListenAndServe()
}

func (svr *Server) hostname() string {
	switch svr.config.Tier {
	case TierProduction:
		return "app.teamxiv.space"
	case TierDevelopment:
		return "dev.teamxiv.space"
	case TierLocal:
		return "localhost:8000"
	}

	// This should not happen if NewConfigFromEnv is used.
	panic("cannot fetch hostname for unknown tier")
}

func (svr *Server) useHTTPS() bool {
	return svr.config.Tier != TierLocal
}

func (svr *Server) protocol() string {
	if svr.useHTTPS() {
		return "https"
	}
	return "http"
}

func (svr *Server) sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	e := json.NewEncoder(w)
	if err := e.Encode(data); err != nil {
		svr.logger.Errorf("could not encode response JSON: %v", err)
	}
}

type apiError struct {
	Code        int    `json:"code"`
	Status      string `json:"status"`
	UserMessage string `json:"message,omitempty"`
}

// nolint: unparam // FIXME remove this later once user message gets used.
func (svr *Server) sendErrorResponse(w http.ResponseWriter, err error,
	statusCode int, userMessage string, args ...interface{}) {

	details := apiError{
		Code:        statusCode,
		Status:      http.StatusText(statusCode),
		UserMessage: fmt.Sprintf(userMessage, args...),
	}

	svr.logger.
		WithError(err).
		WithField("details", details).
		Error("encountered error when handling api request")

	w.WriteHeader(statusCode)
	svr.sendJSONResponse(w, details)
}
