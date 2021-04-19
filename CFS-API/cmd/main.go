package main

import (
	"CallsForService/CFS-API/pkg/config"
	"CallsForService/CFS-API/pkg/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TV4/graceful"
	// log "github.com/sirupsen/logrus"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

// TODO lower() https://makandracards.com/makandra/41768-case-sensitivity-in-postgresql

// Calls For Service API.
//
//
//
//     Schemes: http, https
//     BasePath: /docs
//     Version: 0.0.1
//     Contact: Brian McGraw <mcgraw.bm@gmail.com> https://www.briandavidmcgraw.com
//
//     Produces:
//     - application/json
//
//     Security:
//     - None
//
//
// swagger:meta

func main() {

	config, err := config.New()

	if err != nil {
		log.Fatalf("Error creating configuration: %v", err)
	}

	db, err := ConnectPostgres(config.DBConnectionInfo)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(1000),
	)
	if err != nil {
		log.Fatalf("Error setting up memcache: %v", err)
	}

	// log.SetLevel(log.WarnLevel)

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		log.Fatalf("Error setting up memcache: %v", err)
	}

	router := handlers.NewRouter(db, config.BuildInfo, cacheClient)

	h := &http.Server{Addr: ":" + config.Port, Handler: router}
	graceful.Timeout = time.Duration(25) * time.Second
	log.Printf("API listening on port %v", config.Port)
	graceful.ListenAndServe(h)

}

func ConnectPostgres(connectionInfo config.DBConnectionInfo) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", connectionInfo.Host, connectionInfo.Port, connectionInfo.User, connectionInfo.Password, connectionInfo.DatabaseName))

	if err != nil {
		return db, err
	}

	err = db.Ping()

	if err != nil {
		return db, err
	}

	return db, nil
}
