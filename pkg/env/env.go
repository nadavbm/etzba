package env

import (
	"os"
	"strconv"
)

// Etzba configuration
var ApiRequestTimeout = getEnvVar("", "ETZ_API_REQUEST_TIMEOUT")
var SqlQueryTimeout = getEnvVar("", "ETZ_SQL_QUERY_TIMEOUT")

// Database connection variables
var DatabaseUser = getEnvVar("", "ETZ_POSTGRES_USER")
var DatabasePass = getEnvVar("", "ETZ_POSTGRES_PASSWORD")
var DatabaseDB = getEnvVar("", "ETZ_POSTGRES_DB")
var DatabaseHost = getEnvVar("localhost", "ETZ_POSTGRES_HOST")
var DatabasePort = getEnvInt(5432, "ETZ_POSTGRES_PORT")

// API server authentication
var ApiAuthMethod = getEnvVar("", "ETZ_API_AUTH_METHOD")
var ApiToken = getEnvVar("", "ETZ_API_TOKEN")

// Prometheus configuration
var PrometheusPushGateway = getEnvVar("", "PROMETHEUS_PUSH_GATEWAY")

func getEnvVar(defaultValue, key string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(defaultValue int, key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value == 0 {
		return defaultValue
	}
	return value
}

func getBoolEnv(defaultValue bool, key string) bool {
	val := os.Getenv(key)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return ret
}
