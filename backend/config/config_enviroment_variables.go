package config

import (
	"log"
	"os"
	"strings"
)

func CheckEnvironmentVariables() {
	requiredEnvVariables := []string{
		"GIN_MODE",
		"DOMAIN",
		"HOST",
		"PORT",
		"JWT_KEY",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSL_MODE",
		"DB_TIMEZONE",
		"LOGIN",
		"PASSWORD",
	}
	var notPresentedVariables []string
	for _, envName := range requiredEnvVariables {
		_, present := os.LookupEnv(envName)
		if !present {
			notPresentedVariables = append(notPresentedVariables, envName)
		}
	}

	if len(notPresentedVariables) != 0 {
		log.Fatalf("\nNeeded %d more environment variables:\n%s", len(notPresentedVariables), strings.Join(notPresentedVariables, "\n"))
	}
}
