package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/paulmach/orb/geojson"
	"github.com/spf13/viper"
)

const (
	POSTGRES_DB_ENV  = "POSTGRES_DB_ENV"  // WAS DB_ENV
	POSTGRES_URL_ENV = "POSTGRES_URL_ENV" // Was HOST_ENV
	// TABLE_ENV               = "TABLE_ENV"
	POSTGRES_USER_ENV       = "POSTGRES_USER_ENV" // WAS USER_ENV
	POSTGRES_PORT_ENV       = "POSTGRES_PORT_ENV" // Was PORT_ENV
	POSTGRES_PW_ENV         = "POSTGRES_PW_ENV"   // Was PW_ENV
	CFS_SITE_ENV            = "CFS_SITE_ENV"
	GOOGLE_MAPS_API_KEY_ENV = "GOOGLE_MAPS_API_KEY_ENV"
	DEBUG_ENV               = "DEBUG_ENV"
	GOOGLE_MAPS_URL         = "https://maps.googleapis.com/maps/api/geocode/json"
	WARDS_GEO_FILE_ENV      = "WARDS_GEO_FILE_ENV"
	WARDS_GEO_FILE          = "/wards.geojson"
	SHOULD_CALL_MAPS_ENV    = "SHOULD_CALL_MAPS_ENV"
)

type Config struct {
	CFSWebsite string
	Db         *sql.DB
	// Table          string
	Client      *http.Client
	WardsFile   string
	MapsConfig  MapsConfig
	Debug_Env   string
	Geofile_Loc string
}

type MapsConfig struct {
	URL     string
	API_Key string
}

func New() (*Config, error) {
	c := &Config{}

	host := c.GetHost()
	port := c.GetPort()
	user := c.GetUser()
	password := c.GetPassword()
	c.Geofile_Loc = c.GetWardsGeoFileEnv()
	// c.Table = c.GetTable()
	c.CFSWebsite = c.GetWebsite()
	c.WardsFile = c.Geofile_Loc + WARDS_GEO_FILE
	c.Debug_Env = c.GetDebugEnv()
	c.MapsConfig.URL = GOOGLE_MAPS_URL
	c.MapsConfig.API_Key = c.GetMapsKey()
	c.Client = &http.Client{}

	dbname := c.GetDB()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return c, err
	}

	err = db.Ping()

	if err != nil {
		return c, err
	}

	c.Db = db

	c.Client = &http.Client{
		Timeout: 10 * time.Second,
	}

	return c, err
}

func (c Config) GetWardsGeoFileEnv() string {
	viper.BindEnv(WARDS_GEO_FILE_ENV)
	if !viper.IsSet(WARDS_GEO_FILE_ENV) {
		log.Fatalf("Unable to retrieve environment variable: %v", WARDS_GEO_FILE_ENV)
	}

	return viper.GetString(WARDS_GEO_FILE_ENV)
}

func (c Config) GetMapsKey() string {
	viper.BindEnv(GOOGLE_MAPS_API_KEY_ENV)
	if !viper.IsSet(GOOGLE_MAPS_API_KEY_ENV) {
		log.Fatalf("Unable to retrieve environment variable: %v", GOOGLE_MAPS_API_KEY_ENV)
	}

	return viper.GetString(GOOGLE_MAPS_API_KEY_ENV)
}

func (c Config) GetDebugEnv() string {
	viper.BindEnv(DEBUG_ENV)
	if !viper.IsSet(DEBUG_ENV) {
		log.Fatalf("Unable to retrieve environment variable: %v", DEBUG_ENV)
	}

	env := viper.GetString(DEBUG_ENV)

	if env != "debug" && env != "prod" && env != "test" {
		log.Fatalf("Invalid env value, should be dev test or prod: %v", DEBUG_ENV)
	}

	return env
}

func (c Config) GetHost() string {
	viper.BindEnv(POSTGRES_URL_ENV)
	if !viper.IsSet(POSTGRES_URL_ENV) {
		log.Fatalf("Unable to retrieve environment variable: %v", POSTGRES_URL_ENV)
	}

	return viper.GetString(POSTGRES_URL_ENV)
}

func (c Config) GetPort() string {
	viper.BindEnv(POSTGRES_PORT_ENV)
	if !viper.IsSet(POSTGRES_PORT_ENV) {
		log.Fatalf("Unable to retrieve environment variable: %v", POSTGRES_PORT_ENV)
	}

	return viper.GetString(POSTGRES_PORT_ENV)
}

func (c Config) GetUser() string {
	viper.BindEnv(POSTGRES_USER_ENV)
	if !viper.IsSet(POSTGRES_USER_ENV) {
		log.Fatalf("Unable to retrieve environment variable: %v", POSTGRES_USER_ENV)
	}

	return viper.GetString(POSTGRES_USER_ENV)
}

func (c Config) GetPassword() string {
	viper.BindEnv(POSTGRES_PW_ENV)
	if !viper.IsSet(POSTGRES_PW_ENV) {
		log.Fatalf("Unable to retrieve environment variable: %v", POSTGRES_PW_ENV)
	}

	return viper.GetString(POSTGRES_PW_ENV)
}

func (c Config) GetDB() string {
	viper.BindEnv(POSTGRES_DB_ENV)
	if !viper.IsSet(POSTGRES_DB_ENV) {
		log.Fatalf("Unable to retrieve environment variable: %v", POSTGRES_DB_ENV)
	}

	return viper.GetString(POSTGRES_DB_ENV)
}

func (c Config) GetWebsite() string {
	viper.BindEnv(CFS_SITE_ENV)
	if !viper.IsSet(CFS_SITE_ENV) {
		log.Fatalf("Unable to retrieve environment variable: %v", CFS_SITE_ENV)
	}

	return viper.GetString(CFS_SITE_ENV)
}

func (c Config) BuildWardsFC(filename string) (*geojson.FeatureCollection, error) {

	b, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Unable to open geojson file: %v", err)
	}

	featureCollection, err := geojson.UnmarshalFeatureCollection(b)

	if err != nil {
		log.Fatalf("Unable to decode geojson file: %v", err)
	}

	return featureCollection, err
}
