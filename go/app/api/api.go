package api

import (
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
	// router.Use(mwf ...mux.MiddlewareFunc)

	// Define routes.
	router.Path("/login").Methods("POST").HandlerFunc(svr.handleLogin)
	router.Path("/logout").Methods("POST").HandlerFunc(svr.handleLogout)

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
