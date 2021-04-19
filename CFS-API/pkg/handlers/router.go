package handlers

import (
	"CallsForService/CFS-API/pkg/cfs"
	"CallsForService/CFS-API/pkg/config"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	cache "github.com/victorspringer/http-cache"
)

// cfsraw
// /api/cfs/{id}
// /api/cfs?startDate,endDate,ward,neighborhood (ward OR neighborhood)
// query by eventid, query by date, query by ward, neighborhood

// query by locations, wards, neighborhoods (to preopulate client side query options)

func NewRouter(db *sql.DB, buildInfo config.BuildInfo, cacheClient *cache.Client) *mux.Router {

	r := mux.NewRouter()
	router := r.PathPrefix("/api").Subrouter()

	router.HandleFunc("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("ok"))
	})).Methods(http.MethodGet)

	// GetInfo swagger:route GET /api/info info listInfo
	//
	// Lists information about the state of the api
	//
	//     Produces:
	//     - application/json
	//
	// Deprecated: false
	// Responses:
	// 		    200: info

	router.HandleFunc("/info", Info(buildInfo)).Methods(http.MethodGet)

	// GetCFS swagger:route GET /api/cfs cfs getCFS
	//
	//  Returns calls for service. Defaults to the 100 most recent calls.
	//
	//     Produces:
	//     - application/json
	//
	// Deprecated: false
	// Responses:
	// 		    200: []cfs
	//          404: ErrorResponse
	//          401: ErrorResponse
	//			403: ErrorResponse
	router.Handle("/cfs", cacheClient.Middleware(handleGetCFS(db))).Methods(http.MethodGet)

	// TODO: Implement these routes
	// router.HandleFunc("/cfs/aggregate", handleGetCFSAggregate(config)).Methods(http.MethodGet)
	// router.HandleFunc("/cfs/{id}", handleGetCFSByID()).Methods(http.MethodPatch)

	return router

}

func handleGetCFS(db *sql.DB) http.HandlerFunc {
	return cfs.GetCFS(db)
}

func Info(buildInfo config.BuildInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, _ := json.Marshal(buildInfo)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}
