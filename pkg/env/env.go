package env

import (
	"os"
	"strconv"
)

// Etzba configuration
var ApiRequestTimeout = GetEnvVar("", "ETZ_API_REQUEST_TIMEOUT")
var SqlQueryTimeout = GetEnvVar("", "ETZ_SQL_QUERY_TIMEOUT")

// Database connection variables
var DatabaseUser = GetEnvVar("", "ETZ_POSTGRES_USER")
var DatabasePass = GetEnvVar("", "ETZ_POSTGRES_PASSWORD")
var DatabaseDB = GetEnvVar("", "ETZ_POSTGRES_DB")
var DatabaseHost = GetEnvVar("localhost", "ETZ_POSTGRES_HOST")
var DatabasePort = GetEnvInt(5432, "ETZ_POSTGRES_PORT")

// API server authentication
var ApiAuthMethod = GetEnvVar("", "ETZ_API_AUTH_METHOD")
var ApiToken = GetEnvVar("", "ETZ_API_TOKEN")

func GetEnvVar(defaultValue, key string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(defaultValue int, key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value == 0 {
		return defaultValue
	}
	return value
}
