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

	// Bogus endpoint. Always returns 501.
	router.Path("/todo").Methods("GET").HandlerFunc(svr.handleTODO)

	// Account subroutes.
	accountRouter := router.PathPrefix("/account").Subrouter()
	accountRouter.Path("/register").HandlerFunc(svr.handleRegistration)

	// adminRouter := router.PathPrefix("/admin").Subrouter()
	accountRouter.Path("/forgot").Methods("POST").
		HandlerFunc(svr.handleTODO) // TODO
	accountRouter.Path("/register").Methods("POST").
		HandlerFunc(svr.handleRegistration)

	// My subroutes.
	myRouter := router.PathPrefix("/my").Subrouter()
	myRouter.Use(svr.requireAuthMiddleware(authConfig{
		requireRole: true,
		allowedRoles: []app.Role{
			app.RoleAdmin,
			app.RoleSponsor,
			app.RoleDriver,
		},
	}))

	myProfileRouter := myRouter.PathPrefix("/profile").Subrouter()
	myProfileRouter.Path("/name").Methods("POST").
		HandlerFunc(svr.handleMyProfileUpdateName)
	myProfileRouter.Path("/email").Methods("POST").
		HandlerFunc(svr.handleMyProfileUpdateEmail)
	myProfileRouter.Path("/password").Methods("POST").
		HandlerFunc(svr.handleMyProfileUpdatePassword)
	myProfileRouter.Path("/deactivate").Methods("POST").
		HandlerFunc(svr.handleMyProfileDeactivate)

	// Admin subroutes.
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(svr.requireAuthMiddleware(authConfig{
		requireRole:  true,
		allowedRoles: []app.Role{app.RoleAdmin},
	}))

	adminUserRouter := adminRouter.PathPrefix("/users").Subrouter()
	adminUserRouter.Path("").Methods("GET").
		HandlerFunc(svr.handleAdminGetAllUsers)
	adminUserRouter.Path("/create").Methods("POST").
		HandlerFunc(svr.handleTODO) // TODO
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
		HandlerFunc(svr.handleGetAllOrganizations)
	adminOrgRouter.Path("/create").Methods("POST").
		HandlerFunc(svr.handleAdminCreateOrganization)
	adminOrgRouter.Path("/{orgID}").Methods("GET").
		HandlerFunc(svr.handleAdminGetOrganizationByID)
	adminOrgRouter.Path("/{orgID}/update").Methods("POST").
		HandlerFunc(svr.handleAdminUpdateOrganization)
	adminOrgRouter.Path("/{orgID}/delete").Methods("POST").
		HandlerFunc(svr.handleAdminDeleteOrganization)

	// Sponsor subroutes.
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

	sponsorOrgRouter := sponsorRouter.PathPrefix("/organization").Subrouter()
	sponsorOrgRouter.Path("").Methods("GET").
		HandlerFunc(svr.handleSponsorGetOwnOrganization)
	sponsorOrgRouter.Path("/update").Methods("POST").
		HandlerFunc(svr.handleSponsorUpdateOwnOrganization)

	sponsorDriverRouter := sponsorRouter.PathPrefix("/drivers").Subrouter()
	sponsorDriverRouter.Path("").Methods("GET").
		HandlerFunc(svr.handleTODO) // TODO
	sponsorDriverRouter.Path("/{driverID}").Methods("GET").
		HandlerFunc(svr.handleTODO) // TODO
	sponsorDriverRouter.Path("/{driverID}/points").Methods("POST").
		HandlerFunc(svr.handleTODO) // TODO
	sponsorDriverRouter.Path("/{driverID}/remove").Methods("POST").
		HandlerFunc(svr.handleTODO) // TODO

	sponsorAppRouter := sponsorRouter.PathPrefix("/applications").Subrouter()
	sponsorAppRouter.Path("").Methods("GET").
		HandlerFunc(svr.handleGetApplicationsForOrganization)
	sponsorAppRouter.Path("/{appID}").Methods("GET").
		HandlerFunc(svr.handleGetApplicationByID)
	sponsorAppRouter.Path("{appID}/approve").Methods("POST").
		HandlerFunc(svr.handleApproveApplication)

	driverRouter := router.PathPrefix("/driver").Subrouter()

	driverRouter.Path("/applications/submit").Methods("POST").
		HandlerFunc(svr.handleSubmitApplication)
	driverRouter.Path("/applications/{appID}").Methods("GET").
		HandlerFunc(svr.handleGetApplicationByID)
	driverRouter.Path("/applications").Methods("GET").
		HandlerFunc(svr.handleGetMyApplications)

	driverRouter.Path("/balances").Methods("GET").
		HandlerFunc(svr.handleDriverGetBalances)
	driverRouter.Path("/organizations/all").Methods("GET").
		HandlerFunc(svr.handleGetAllOrganizations)
	driverRouter.Path("/catalog/{orgID}/search").Methods("GET").
		HandlerFunc(svr.handleTODO) // TODO

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
