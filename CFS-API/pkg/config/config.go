package config

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	DB_ENV            = "DB_ENV"
	HOST_ENV          = "HOST_ENV"
	USER_ENV          = "USER_ENV"
	POSTGRES_PORT_ENV = "POSTGRES_PORT_ENV"
	PW_ENV            = "PW_ENV"
	API_PORT_ENV      = "API_PORT_ENV"
	DEFAULT_API_PORT  = "8000"
	BUILD_COMMIT_HASH = "BUILD_COMMIT_HASH"
	BUILD_USER        = "BUILD_USER"
	BUILD_TIME        = "BUILD_TIME"
)

type Config struct {
	Port             string
	DBConnectionInfo DBConnectionInfo
	BuildInfo        BuildInfo
}

type DBConnectionInfo struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
}

// swagger:response info
type BuildInfoResponse struct {
	BuildInfo BuildInfo `json:"buildInfo"`
}

// Info represents the current state of the API.
//
// swagger:model info
type BuildInfo struct {
	BuildUser  string `json:"user"`
	BuildTime  string `json:"buildTime"`
	CommitHash string `json:"commitHash"`
}

func New() (*Config, error) {

	c := &Config{}

	c.BuildInfo = BuildInfo{}
	c.BuildInfo.BuildUser = c.GetOptionalConfig(BUILD_USER, "missing")
	c.BuildInfo.BuildTime = c.GetOptionalConfig(BUILD_TIME, "missing")
	c.BuildInfo.CommitHash = c.GetOptionalConfig(BUILD_COMMIT_HASH, "missing")

	c.Port = c.GetOptionalConfig(API_PORT_ENV, DEFAULT_API_PORT)

	c.DBConnectionInfo = DBConnectionInfo{}
	c.DBConnectionInfo.Host = c.GetRequiredConfig(HOST_ENV)
	c.DBConnectionInfo.Port = c.GetRequiredConfig(POSTGRES_PORT_ENV)

	c.DBConnectionInfo.User = c.GetRequiredConfig(USER_ENV)
	c.DBConnectionInfo.Password = c.GetRequiredConfig(PW_ENV)

	c.DBConnectionInfo.DatabaseName = c.GetRequiredConfig(DB_ENV)

	return c, nil
}

func (c Config) GetRequiredConfig(env string) string {
	viper.BindEnv(env)
	if !viper.IsSet(env) {
		log.Fatalf("Unable to retrieve environment variable: %v", env)
	}

	return viper.GetString(env)
}

func (c Config) GetOptionalConfig(environment, default_environment string) string {
	viper.BindEnv(environment)
	if !viper.IsSet(environment) {
		log.Printf("Unable to retrieve environment variable: %v, using default value: %v", environment, default_environment)
		return default_environment
	}

	return viper.GetString(environment)
}
