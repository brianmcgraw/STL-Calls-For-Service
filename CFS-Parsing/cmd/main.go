package main

import (
	"CallsForService/CFS-Parsing/pkg/config"
	"CallsForService/CFS-Parsing/pkg/parse"
	"CallsForService/CFS-Parsing/pkg/pg"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("cfs-parsing.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	config, err := config.New()

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Get new CFS updates
	cfs, err := parse.CallCFS(config)
	if err != nil {
		log.Fatalf("Error calling CFS: %v", err)
	}

	// Insert to raw cfs table
	err = pg.InsertToRawCFS(config.Db, cfs)

	if err != nil {
		log.Fatalf("Error inserting data to raw CFS table: %v", err)
	}

}
