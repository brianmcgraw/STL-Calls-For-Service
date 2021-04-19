package cfs

import "time"

type RawCFS struct {
	EventTime   time.Time `json:"eventtime"`
	Id          string    `json:"Id"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}

// A call for service represents information about a call made requesting service.
//
// swagger:model cfs
type CallForService struct {
	// The identification of the call.
	Id           string    `json:"Id"`
	Description  string    `json:"description"`
	EventTime    time.Time `json:"eventTime"`
	Location     string    `json:"location"`
	Lat          float64   `json:"latitude"`
	Lng          float64   `json:"longitude"`
	Ward         string    `json:"ward"`
	Neighborhood string    `json:"neighborhood"`
}

type CallForServiceQueryResponse struct {
	Count int `json:"count"`
	CallForService
}

type Location struct {
	Location           string  `json:"location"`
	NormalizedLocation string  `json:"normalizedLocation"`
	Lat                float64 `json:"lat"`
	Lng                float64 `json:"lng"`
	Ward               string  `json:"ward"`
	Neighborhood       string  `json:"neighborhood"`
	Zipcode            string  `json:"zipcode"`
	HasIssue           bool    `json:"hasIssue"`
}

// A UserParams model.
//
// swagger:parameters getCFS
type SearchQuery struct {
	// The starting cutoff for your query. Defaults to 24 hours prior to the current time. Expressed as ISO 8601 string.
	//
	// in: request parameter
	// required: false
	StartDate string `json:"startDate" schema:"startDate"`
	// The end date cutoff for your query. Defaults to the current time. Expressed as ISO 8601 string.
	//
	// in: request parameter
	// required: false
	EndDate string `json:"endDate" schema:"endDate"`
	// Limits your query to calls within a specific ward. Must be an integer between 1 and 28.
	//
	// in: request parameter
	// required: false
	Ward string `json:"ward" schema:"ward"`
	// Limits your query to calls within a neighborhood.
	// List of neighborhoods: https://github.com/slu-openGIS/STL_BOUNDARY_Nhood/tree/master/data
	// in: request parameter
	// required: false
	Neighborhood string `json:"neighborhood" schema:"neighborhood"`

	// Limits your query to a given call description.
	// TODO --- map description to a description code.
	// in: request parameter
	// required: false
	Description string `json:"description" schema:"description"`
	// Used in paginating requests. An offset of 10 would skip the first 10 results. Defaults to zero.
	//
	// in: request parameter
	// required: false
	Offset *int `json:"offset" schema:"offset"`
	// Limits the number of calls returned. Default is 100. Maximum is 500.
	//
	// in: request parameter
	// required: false
	Limit *int `json:"limit" schema:"limit"`
}

// swagger:response getCFS
type SearchQueryResponse struct {
	CFS  []CallForService `json:"cfs"`
	Meta Meta             `json:"meta"`
}

type Meta struct {
	Offset       string `json:"offset"`
	PageNumber   int    `json:"pageNumber"`
	PageSize     int    `json:"pageSize"`
	TotalRecords int    `json:"totalRecords"`
	Next         string `json:"next"`
	Prev         string `json:"prev"`
}

// swagger:response ErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
}
