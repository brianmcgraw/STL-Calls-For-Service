package main

import (
	"CallsForService/CFS-Geocoding/pkg/config"
	"CallsForService/CFS-Geocoding/pkg/geocode"

	"CallsForService/CFS-Geocoding/pkg/pg"
	"CallsForService/CFS-Geocoding/pkg/wards"
	"fmt"
	"log"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func main() {
	file, err := os.OpenFile("cfs-maps-geocoding.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	config, err := config.New()

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	wardsFC, err := config.BuildWardsFC(config.WardsFile)

	if err != nil {
		log.Fatalf("Error building wards FC: %v", err)
	}

	neighborhoodsFC, err := config.BuildNeighborhoodFC()

	if err != nil {
		log.Fatalf("Error building wards FC: %v", err)
	}

	nbrAdded, err := shouldCallMaps(config, wardsFC, neighborhoodsFC)

	if err != nil {
		log.Println("Error from google maps geocoding: ", err)
	}

	log.Printf("Added %v new locations", nbrAdded)

}

func shouldCallMaps(config *config.Config, wardsFC *geojson.FeatureCollection, neighborhoodsFC *geojson.FeatureCollection) (int64, error) {

	// Query locations table for new locations, returns a string array
	locs, err := pg.QueryNewLocations(config.Db)

	if len(locs) == 0 {
		return 0, nil
	}

	if err != nil {
		log.Printf("Error received while querying for new locations: %v", err)
	}
	var improvedLocs []pg.Location

	// Query google maps API with location, should return the modified google maps object (lat/long)
	for _, loc := range locs {
		improvedLoc, err := geocode.CallMaps(config, loc) //add config variable
		if err != nil {
			// TODO, multiple error types. If error in request build, etc. logFatal
			// if error is from a status code, consider ignoring
			log.Printf("Error calling google maps API: %v", err)

		}
		improvedLocs = append(improvedLocs, improvedLoc)
	}

	// Then calculate the neighborhood and ward
	for index, value := range improvedLocs {

		wardNbr, isFound := wards.IsPointInsidePolygon(wardsFC.Features, orb.Point{value.Lng, value.Lat})

		if isFound {
			improvedLocs[index].Ward = fmt.Sprintf("%v", wardNbr)
		}

		neighborHood, isFound := wards.IsPointInsideMultiPolygon(neighborhoodsFC.Features, orb.Point{value.Lng, value.Lat})

		if isFound {
			improvedLocs[index].Neighborhood = neighborHood
		}
	}

	// Upload new locations to location table
	nbrAdded, err := pg.InsertModifiedLocations(config.Db, improvedLocs)

	if err != nil {
		log.Printf("Error inserting records to modified location table: %v", err)
	}

	return nbrAdded, err

}
