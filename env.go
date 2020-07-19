package main

import (
	"fmt"
	"os"
)

const (
	DiscordBotTokenEnv = "DISCORD_BOT_TOKEN"
	EnvironmentEnv     = "ENVIRONMENT"
	GiphyAPIKeyEnv     = "GIPHY_API_KEY"
)

var Environment string

func assertGetEnv(envVar string) string {
	value := os.Getenv(envVar)
	if value == "" {
		panic(fmt.Sprintf("%s was not found in env.", envVar))
	}
	return value
}

func GetDiscordBotToken() string {
	return assertGetEnv(DiscordBotTokenEnv)
}

func GetEnvironment() string {
	if Environment == "" {
		Environment = assertGetEnv(EnvironmentEnv)
	}
	return Environment
}

func SetEnvironment(environment string) {
	Environment = environment
}

func GetGiphyAPIKey() string {
	return assertGetEnv(GiphyAPIKeyEnv)
}

func IsProduction(environment string) bool {
	return environment == "production"
}
