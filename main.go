package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/madchin/trader-bot/internal/gateway"
	"github.com/madchin/trader-bot/internal/scheduler"
	"github.com/madchin/trader-bot/internal/service"
	"github.com/madchin/trader-bot/internal/storage"
	"github.com/madchin/trader-bot/internal/worker"
)

func main() {
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	botToken, applicationId, guildId := requiredEnvs()
	gateway, err := gateway.NewGatewaySession(botToken, applicationId, guildId, scheduler.Scheduler)
	if err != nil {
		panic(err)
	}
	defer gateway.CloseSession()
	if err := gateway.OpenConnection(); err != nil {
		panic(err)
	}
	storage := storage.New()
	factoryWorkers := worker.NewFactory(100)
	service := service.New(storage, gateway)
	go worker.Spawner(ctx, service, scheduler.Scheduler, factoryWorkers)
	<-ctx.Done()
}

/*
Retrieve required environment variables for development / production run-time environments.
Function retrieves path to .env file from ENV_FILE environment variable
(path points to /, in order to change it, head to compose.yaml secrets directive)

Needed ENV variables for project to run:
  - BOT_TOKEN (settled in .env) -- DISCORD bot token
  - APPLICATION_ID (settled in .env) -- DISCORD application id
  - GUILD_ID (settled in .env) [needed only when RUNTIME_ENVIRONMENT=DEV] -- DISCORD guild id
  - ENV_FILE (settled in compose.yaml) -- provides path to .env file
  - RUNTIME_ENVIRONMENT (settled dynamically in docker compose up command [check it in Makefile]) -- determines runtime environment
*/
func requiredEnvs() (botToken, applicationId, guildId string) {
	path := os.Getenv("ENV_FILE")
	if path == "" {
		panic(errors.New("ENV_FILE environment variable not provided. It needs to be set with path to .env file"))
	}
	err := godotenv.Load(path)
	if err != nil {
		panic(fmt.Errorf("unable to load environments from env file %v", err))
	}
	botToken = os.Getenv("BOT_TOKEN")
	if botToken == "" {
		panic(errors.New("BOT_TOKEN environment variable not provided"))
	}
	applicationId = os.Getenv("APPLICATION_ID")
	if applicationId == "" {
		panic(errors.New("APPLICATION_ID environment variable not provided"))
	}
	if os.Getenv("RUNTIME_ENVIRONMENT") == "DEV" {
		log.Printf("RUNTIME_ENV %v", os.Getenv("RUNTIME_ENVIRONMENT"))
		guildId = os.Getenv("GUILD_ID")
		if guildId == "" {
			panic(errors.New("GUILD_ID environment not provided. Its required in DEV run-time environment"))
		}
	}
	return
}
