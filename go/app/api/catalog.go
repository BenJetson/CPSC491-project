package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func getOrganizationIDOfSponsor(r *http.Request) (orgID int, err error) {
	s := getSessionFromContext(r.Context())
	if s == nil {
		err = errors.New("no session")
		return
	}

	if s.Role != app.RoleSponsor {
		err = errors.New("not a sponsor")
		return
	}

	if len(s.Affiliations) != 1 {
		err = errors.Errorf("expected 1 affiliation but found %d",
			len(s.Affiliations))
		return
	}

	orgID = s.Affiliations[0]
	return
}

func (svr *Server) handleSponsorVendorSearch(
	w http.ResponseWriter,
	r *http.Request,
) {

	keywords := r.URL.Query().Get("q")
	if len(keywords) < 1 {
		svr.sendErrorResponse(w, errors.New("missing vendor search keywords"),
			http.StatusBadRequest, "Must supply keywords.")
		return
	}

	ps, err := svr.cv.Search(r.Context(),
		app.CommerceQuery{Keywords: keywords})
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "vendor search failed"),
			http.StatusInternalServerError, "")
		return
	}

	svr.sendJSONResponse(w, ps)
}

func (svr *Server) fetchVendorProductFromURL(
	w http.ResponseWriter,
	r *http.Request,
) (app.CommerceProduct, bool) {

	pathParams := mux.Vars(r)

	productID, err := strconv.Atoi(pathParams["productID"])
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "productID must be an integer"),
			http.StatusBadRequest, "Product ID must be an integer.")
		return app.CommerceProduct{}, false
	}

	p, err := svr.cv.GetProductByID(r.Context(), productID)
	if errors.Is(err, app.ErrNotFound) {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "no such vendor product was found"),
			http.StatusNotFound, "No product with ID of %d.", productID)
		return app.CommerceProduct{}, false
	} else if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "vendor product fetch failed"),
			http.StatusInternalServerError, "")
		return app.CommerceProduct{}, false
	}

	return p, true
}

func (svr *Server) handleSponsorVendorProductByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	vp, ok := svr.fetchVendorProductFromURL(w, r)
	if !ok {
		return
	}

	svr.sendJSONResponse(w, vp)
}

func (svr *Server) handleSponsorAddVendorProduct(
	w http.ResponseWriter,
	r *http.Request,
) {

	orgID, err := getOrganizationIDOfSponsor(r)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "cannot determine sponsor organization identity"),
			http.StatusInternalServerError, "")
		return
	}

	vp, ok := svr.fetchVendorProductFromURL(w, r)
	if !ok {
		return
	}

	p := vp.ToProduct(orgID)

	_, err = svr.db.AddProduct(r.Context(), p)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to save new product"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (svr *Server) handleGetSponsorCatalog(
	w http.ResponseWriter,
	r *http.Request,
) {

	orgID, err := getOrganizationIDOfSponsor(r)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "cannot determine sponsor organization identity"),
			http.StatusInternalServerError, "")
		return
	}

	cps, err := svr.db.GetProductsForOrganization(r.Context(), orgID)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to retrieve catalog products"),
			http.StatusInternalServerError, "")
		return
	}

	svr.sendJSONResponse(w, cps)
}

func (svr *Server) handleGetSponsorCatalogProduct(
	w http.ResponseWriter,
	r *http.Request,
) {

	orgID, err := getOrganizationIDOfSponsor(r)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "cannot determine sponsor organization identity"),
			http.StatusInternalServerError, "")
		return
	}

	pathParams := mux.Vars(r)

	productID, err := strconv.Atoi(pathParams["productID"])
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "productID must be an integer"),
			http.StatusBadRequest, "Product ID must be an integer.")
		return
	}

	cp, err := svr.db.GetProductByID(r.Context(), productID, orgID)
	if errors.Is(err, app.ErrNotFound) {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "no such catalog product was found"),
			http.StatusNotFound,
			"No product with ID of %d in organization %d.", productID, orgID)
		return
	} else if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "catalog product fetch failed"),
			http.StatusInternalServerError, "")
		return
	}

	svr.sendJSONResponse(w, cp)
}

func (svr *Server) handleSponsorRemoveProduct(
	w http.ResponseWriter,
	r *http.Request,
) {

	orgID, err := getOrganizationIDOfSponsor(r)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "cannot determine sponsor organization identity"),
			http.StatusInternalServerError, "")
		return
	}

	pathParams := mux.Vars(r)

	productID, err := strconv.Atoi(pathParams["productID"])
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "productID must be an integer"),
			http.StatusBadRequest, "Product ID must be an integer.")
		return
	}

	err = svr.db.MakeProductUnavailable(r.Context(), productID, orgID)
	if errors.Is(err, app.ErrNotFound) {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "no such catalog product was found"),
			http.StatusNotFound,
			"No product with ID of %d in organization %d.", productID, orgID)
		return
	} else if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "catalog product fetch failed"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
