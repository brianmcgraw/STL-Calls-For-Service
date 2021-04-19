package pg

import (
	"CallsForService/CFS-Parsing/pkg/parse"
	"database/sql"
	"log"
)

const DuplicateKeyError = `pq: duplicate key value violates unique constraint "cfs_pkey"`

func InsertToRawCFS(db *sql.DB, cfs []parse.CFS) (realError error) {
	nbrAdded := 0
	nbrDuplicate := 0
	for _, value := range cfs {
		sqlStatement := `
			INSERT INTO cfs (id, location, description, eventtime, raweventtime, eventdate)
			VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := db.Exec(sqlStatement, value.Id, value.Location, value.Description, value.EventTime, value.RawEventTime, value.EventDate)

		if err != nil {
			if err.Error() != DuplicateKeyError {
				log.Printf("Error from inserting into cfs: %v", err.Error())
				realError = err
			} else {
				nbrDuplicate++
			}

		} else {
			nbrAdded++
		}

	}

	log.Printf("Nbr duplicates: %v, number added %v", nbrDuplicate, nbrAdded)

	return realError
}

// SELECT LOCATION FROM cfsraw WHERE LOCATION NOT IN (SELECT LOCATION FROM location);
func QueryNewLocations(db *sql.DB) (loc []string, err error) {
	sqlStatement := `SELECT LOCATION FROM cfsraw WHERE LOCATION NOT IN (SELECT LOCATION FROM location);`
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

func InsertModifiedLocations(db *sql.DB, loc []parse.Location) (updatedRows int64, realError error) {
	for _, value := range loc {
		sqlStatement := `
			INSERT INTO location (location, lat, lng, ward, neighborhood, zipcode)
			VALUES ($1, $2, $3, $4, $5, $6)`
		result, err := db.Exec(sqlStatement, value.Location, value.Lat, value.Lng, value.Ward, value.Neighborhood, value.Zipcode)

		if err != nil {
			if err.Error() != DuplicateKeyError {
				log.Printf("Error from location db insert: %v", err.Error())
				realError = err
			}

		}

		rowAdded, err := result.RowsAffected()

		if err != nil {
			log.Println("Error from location db insert: %v", err.Error())
			realError = err
		}

		updatedRows = updatedRows + rowAdded

	}

	return updatedRows, realError
}
