package cfs

import (
	"CallsForService/CFS-API/pkg/pg"
	"CallsForService/CFS-API/pkg/util"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/schema"
)

func GetCFS(db pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var search SearchQuery

		if err := schema.NewDecoder().Decode(&search, r.URL.Query()); err != nil {
			util.RespondJSON(http.StatusBadRequest, fmt.Sprintf("Error: %v", err), w)
			return

		}

		search, err := validateLimitOffset(search)

		if err != nil {
			util.RespondJSON(http.StatusBadRequest, fmt.Sprintf("Error: %v", err), w)
			return
		}

		query := BaseCFSQuery

		query, queryArguments, err := CFSQueryBuilder(BaseCFSQuery, search)

		if err != nil {
			util.RespondJSON(http.StatusBadRequest, fmt.Sprintf("Error: %v", err), w)
			return
		}

		log.Println("Printing query: ", query)
		log.Println("Printing queryArgs: ", queryArguments)
		result, err := db.Query(query, queryArguments...)

		if err != nil {
			log.Println(err)
			util.RespondJSON(http.StatusBadRequest, fmt.Sprintf("Error: %v", err), w)
			return
		}

		cfsResult, count, err := ParseCFSResult(result)

		if err != nil {
			log.Println("Error querying the DB...")
			util.RespondJSON(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err), w)
			return
		}

		resp := SearchQueryResponse{CFS: cfsResult, Meta: Meta{
			TotalRecords: count,
			PageSize:     *search.Limit,
			PageNumber:   *search.Offset,
		}}

		util.RespondJSON(http.StatusOK, resp, w)
		return
	}
}

func CFSQueryBuilder(query string, search SearchQuery) (string, []interface{}, error) {
	updatedQuery := query
	var queryArguments []interface{}

	// At minimum we have two arguments, a start and end date
	totalArguments := 1

	if search.StartDate != "" {
		time, err := time.Parse(time.RFC3339, search.StartDate)
		if err != nil {
			return "", queryArguments, errors.New("Invalid value for startDate query parameter, must pass ISO 8601 string.")
		}

		queryArguments = append(queryArguments, time)
		updatedQuery = updatedQuery + fmt.Sprintf("AND eventtime >= $%v", totalArguments)
		totalArguments++
	}

	if search.EndDate != "" {
		time, err := time.Parse(time.RFC3339, search.EndDate)
		if err != nil {
			return "", queryArguments, err
		}
		queryArguments = append(queryArguments, time)
		updatedQuery = updatedQuery + fmt.Sprintf("AND eventtime <= $%v", totalArguments)
		totalArguments++
	}

	if search.Ward != "" {
		queryArguments = append(queryArguments, search.Ward)
		updatedQuery = updatedQuery + fmt.Sprintf("AND ward=$%v", totalArguments)
		totalArguments++
	}

	if search.Neighborhood != "" {
		queryArguments = append(queryArguments, strings.ToLower(search.Neighborhood))
		updatedQuery = updatedQuery + fmt.Sprintf("AND lower(neighborhood)=$%v", totalArguments)
		totalArguments++
	}

	if search.Description != "" {
		queryArguments = append(queryArguments, search.Description)
		updatedQuery = updatedQuery + fmt.Sprintf("AND description=$%v", totalArguments)
		totalArguments++
	}

	updatedQuery = updatedQuery + fmt.Sprint("\nORDER BY cfs.eventTime desc")
	updatedQuery = updatedQuery + fmt.Sprintf("\nLIMIT %v", *search.Limit)
	updatedQuery = updatedQuery + fmt.Sprintf("\nOFFSET %v", *search.Offset)

	return updatedQuery, queryArguments, nil

}

func ParseCFSResult(result *sql.Rows) ([]CallForService, int, error) {
	cfsResult := []CallForService{}
	var count int
	for result.Next() {
		var r CallForServiceQueryResponse
		if err := result.Scan(&r.Id, &r.Description, &r.EventTime, &r.Location, &r.Lat, &r.Lng, &r.Ward, &r.Neighborhood, &r.Count); err != nil {
			if err.Error() == "ErrNoRows" {
				return cfsResult, 0, nil
			}
			return cfsResult, 0, err
		}
		count = r.Count
		cfsResult = append(cfsResult, CallForService{
			Id:           r.Id,
			Description:  r.Description,
			EventTime:    r.EventTime,
			Location:     r.Location,
			Lat:          r.Lat,
			Lng:          r.Lng,
			Ward:         r.Ward,
			Neighborhood: r.Neighborhood,
		})
	}
	return cfsResult, count, nil
}
