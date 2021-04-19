package pg

import (
	"database/sql"
	"log"
	"time"
)

type CFS struct {
	EventTime    time.Time `json:"eventtime"`
	RawEventTime string    `json:"rawEventTime"`
	EventDate    string    `json:"eventDay"`
	Id           string    `json:"Id"`
	Location     string    `json:"location"`
	Description  string    `json:"description"`
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

// SELECT LOCATION FROM cfsraw WHERE LOCATION NOT IN (SELECT LOCATION FROM location);
const DuplicateKeyError = `pq: duplicate key value violates unique constraint "location_pkey"`

func QueryNewLocations(db *sql.DB) (loc []string, err error) {
	sqlStatement := `SELECT LOCATION FROM cfs WHERE LOCATION NOT IN (SELECT LOCATION FROM location);`
	result, err := db.Query(sqlStatement)
	if err != nil {
		return loc, err
	}
	defer result.Close()

	for result.Next() {
		var location string
		if err := result.Scan(&location); err != nil {
			log.Fatal(err)
		}
		loc = append(loc, location)
	}
	// Check for errors from iterating over rows.
	if err := result.Err(); err != nil {
		log.Fatal(err)
	}

	return loc, err

}

func InsertModifiedLocations(db *sql.DB, loc []Location) (updatedRows int64, realError error) {
	for _, value := range loc {
		var rowAdded int64
		sqlStatement := `
			INSERT INTO location (location, latitude, longitude, ward, neighborhood, zipcode, hasissue, normalizedlocation)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		result, err := db.Exec(sqlStatement, value.Location, value.Lat, value.Lng, value.Ward, value.Neighborhood, value.Zipcode, value.HasIssue, value.NormalizedLocation)

		if err != nil {
			if err.Error() != DuplicateKeyError {
				log.Printf("Error other than duplicate key error from db insert: %v", err.Error())
				realError = err
			}
			err = nil

		} else {
			rowAdded, err = result.RowsAffected()
		}

		if err != nil {
			log.Printf("Error from retrieving number of effected rows: %v", err.Error())
			realError = err
		}

		updatedRows = updatedRows + rowAdded

	}

	return updatedRows, realError
}
