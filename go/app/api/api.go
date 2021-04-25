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
	cv     app.CommerceVendor
	logger *logrus.Logger
	httpd  *http.Server
	router *mux.Router
}

// NewServer creates a new Server given a logger, data store, and configuration.
func NewServer(logger *logrus.Logger, db app.DataStore, cv app.CommerceVendor,
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
		cv:     cv,
		logger: logger,
		router: router,
	}

	// Register global middleware.
	router.Use(svr.panicRecoveryMiddleware)
	router.Use(svr.authContextMiddleware)

	// Set custom error handlers.
	router.NotFoundHandler = http.
		HandlerFunc(svr.handleNotFound)
	router.MethodNotAllowedHandler = http.
		HandlerFunc(svr.handleMethodNotAllowed)

	// Define routes.
	router.Path("/login").Methods("POST").HandlerFunc(svr.handleLogin)
	router.Path("/logout").Methods("POST").HandlerFunc(svr.handleLogout)
	router.Path("/whoami").Methods("GET").HandlerFunc(svr.handleWhoAmI)

	// Account subroutes.
	accountRouter := router.PathPrefix("/account").Subrouter()
	accountRouter.Path("/forgot").Methods("POST").
		HandlerFunc(svr.handleTODO) // TODO
	accountRouter.Path("/register").Methods("POST").
		HandlerFunc(svr.handleRegistration)

	// Admin subroutes.
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(svr.requireAuthMiddleware(authConfig{
		requireRole:  true,
		allowedRoles: []app.Role{app.RoleAdmin},
	}))

	adminUserRouter := adminRouter.PathPrefix("/users").Subrouter()
	adminUserRouter.Path("").Methods("GET").
		HandlerFunc(svr.handleAdminGetAllUsers)
	adminUserRouter.Path("/{userID}").Methods("GET").
		HandlerFunc(svr.handleAdminGetUserByID)
	adminUserRouter.Path("/{userID}/name").Methods("POST").
		HandlerFunc(svr.handleAdminUpdateUserName)
	adminUserRouter.Path("/{userID}/email").Methods("POST").
		HandlerFunc(svr.handleAdminUpdateUserEmail)
	adminUserRouter.Path("/{userID}/affiliations").Methods("POST").
		HandlerFunc(svr.handleTODO) // TODO
	adminUserRouter.Path("/{userID}/password").Methods("POST").
		HandlerFunc(svr.handleAdminUpdateUserPassword)
	adminUserRouter.Path("/{userID}/activate").Methods("POST").
		HandlerFunc(svr.handleAdminActivateUser)
	adminUserRouter.Path("/{userID}/deactivate").Methods("POST").
		HandlerFunc(svr.handleAdminDeactivateUser)

	adminOrgRouter := adminRouter.PathPrefix("/organizations").Subrouter()
	adminOrgRouter.Path("").Methods("GET").
		HandlerFunc(svr.handleTODO) // TODO
	adminOrgRouter.Path("/{orgID}").Methods("GET").
		HandlerFunc(svr.handleTODO) // TODO
	adminOrgRouter.Path("/{orgID}/update").Methods("POST").
		HandlerFunc(svr.handleTODO) // TODO
	adminOrgRouter.Path("/{orgID}/delete").Methods("POST").
		HandlerFunc(svr.handleTODO) // TODO
	adminOrgRouter.Path("/create").Methods("POST").
		HandlerFunc(svr.handleTODO) // TODO

	sponsorRouter := router.PathPrefix("/sponsor").Subrouter()
	sponsorRouter.Use(svr.requireAuthMiddleware(authConfig{
		requireRole:  true,
		allowedRoles: []app.Role{app.RoleSponsor},
	}))

	sponsorVendorRouter := sponsorRouter.PathPrefix("/vendor").Subrouter()
	sponsorVendorRouter.Path("/search").Methods("GET").
		HandlerFunc(svr.handleSponsorVendorSearch)
	sponsorVendorRouter.Path("/products/{productID}").Methods("GET").
		HandlerFunc(svr.handleSponsorVendorProductByID)
	sponsorVendorRouter.Path("/products/{productID}/add").Methods("POST").
		HandlerFunc(svr.handleSponsorAddVendorProduct)

	sponsorCatalogRouter := sponsorRouter.PathPrefix("/catalog").Subrouter()
	sponsorCatalogRouter.Path("").Methods("GET").
		HandlerFunc(svr.handleGetSponsorCatalog)
	sponsorCatalogRouter.Path("/products/{productID}").Methods("GET").
		HandlerFunc(svr.handleGetSponsorCatalogProduct)
	sponsorCatalogRouter.Path("/products/{productID}/remove").Methods("POST").
		HandlerFunc(svr.handleSponsorRemoveProduct)

	driverRouter := router.PathPrefix("/driver").Subrouter()
	driverRouter.Use(svr.requireAuthMiddleware(authConfig{
		requireRole:  true,
		allowedRoles: []app.Role{app.RoleDriver},
	}))

	driverRouter.Path("/balances").Methods("GET").HandlerFunc(svr.handleTODO)

	appRouter := router.Path("/applications").Subrouter()

	appRouter.Path("/submit").Methods("POST").
		HandlerFunc(svr.handleSubmitApplication)
	appRouter.Path("/id/{appID}").Methods("GET").
		HandlerFunc(svr.handleGetApplicationByID)
	appRouter.Path("/person/{personID}").Methods("GET").
		HandlerFunc(svr.handleGetApplicationsForPerson)
	appRouter.Path("/organization/{orgID}").Methods("GET").
		HandlerFunc(svr.handleGetApplicationsForOrganization)
	appRouter.Path("/approve/{appID}").Methods("POST").
		HandlerFunc(svr.handleApproveApplication)

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

// hostname returns the server hostname.
func (svr *Server) hostname() string {
	switch svr.config.Tier {
	case TierProduction:
		return "app.teamxiv.space"
	case TierDevelopment:
		return "dev.teamxiv.space"
	case TierLocal:
		return "localhost"
	}

	// This should not happen if NewConfigFromEnv is used.
	panic("cannot fetch hostname for unknown tier")
}

// useHTTPS returns true when Nginx is configured for HTTPS.
func (svr *Server) useHTTPS() bool {
	return svr.config.Tier != TierLocal
}

// protocol returns the protocol the app will use for communication.
// nolint: unused // FIXME remove this once used.
func (svr *Server) protocol() string {
	if svr.useHTTPS() {
		return "https"
	}
	return "http"
}

// sendJSONResponse will marshal the given data to JSON and write it to the
// http ResponseWriter.
func (svr *Server) sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	e := json.NewEncoder(w)
	if err := e.Encode(data); err != nil {
		svr.logger.Errorf("could not encode response JSON: %v", err)
	}
}

// An apiError describes a problem that the server encountered while processing
// a request, with an optional user message.
type apiError struct {
	Code        int    `json:"code"`
	Status      string `json:"status"`
	UserMessage string `json:"message,omitempty"`
}

// sendErrorResponse sends a completed apiError back to the user as JSON and
// logs the error.
func (svr *Server) sendErrorResponse(w http.ResponseWriter, err error,
	statusCode int, userMessage string, args ...interface{}) {

	details := apiError{
		Code:        statusCode,
		Status:      http.StatusText(statusCode),
		UserMessage: fmt.Sprintf(userMessage, args...),
	}

	svr.logger.
		WithError(err).
		WithField("details", fmt.Sprintf("%+v", details)).
		Error("encountered error when handling api request")

	w.WriteHeader(statusCode)
	svr.sendJSONResponse(w, details)
}
