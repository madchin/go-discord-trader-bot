package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/madchin/trader-bot/internal/gateway"
	"github.com/madchin/trader-bot/internal/scheduler"
	"github.com/madchin/trader-bot/internal/service"
	"github.com/madchin/trader-bot/internal/storage"
	"github.com/madchin/trader-bot/internal/worker"
)

type appEnvs struct {
	botToken string
	appId    string
	guildId  string
}

type envs struct {
	runtimeEnvironment string
	app                appEnvs
}

func main() {
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	envs, err := requiredEnvs()
	if err != nil {
		panic(err)
	}
	dbCreds, err := storage.LoadCredentials()
	if err != nil {
		panic(err)
	}
	conn, err := storage.Connect(dbCreds)
	if err != nil {
		panic(err)
	}
	storage := storage.New(conn)

	gateway, err := gateway.NewGatewaySession(envs.app.botToken, envs.app.appId, envs.app.guildId, scheduler.Scheduler)
	if err != nil {
		panic(err)
	}
	defer gateway.CloseSession()
	if err := gateway.OpenConnection(); err != nil {
		panic(err)
	}

	service := service.New(storage, gateway)
	factoryWorkers := worker.NewFactory(100)
	go worker.Spawner(ctx, service, scheduler.Scheduler, factoryWorkers)
	<-ctx.Done()
}

/*
Retrieve required environment variables for development / production run-time environments.
Function retrieves path to .env file from ENV_FILE environment variable
(path points to /, in order to change it, head to compose.yaml secrets directive)

Needed ENV variables for project to run:
  - POSTGRES_DB (settled in .env) -- POSTGRES database name
  - POSTGRES_USER (settled in .env) -- POSTGRES user name
  - POSTGRES_PASSWORD (settled in .env) -- POSTGRES user password
  - BOT_TOKEN (settled in .env) -- DISCORD bot token
  - APPLICATION_ID (settled in .env) -- DISCORD application id
  - GUILD_ID (settled in .env) [needed only when RUNTIME_ENVIRONMENT=DEV] -- DISCORD guild id
  - ENV_FILE (settled in compose.yaml) -- provides path to .env file
  - RUNTIME_ENVIRONMENT (settled dynamically in docker compose up command [check it in Makefile]) -- determines runtime environment
*/
func requiredEnvs() (envs envs, err error) {
	appEnvsFilePath := os.Getenv("APP_ENV_FILE")
	if appEnvsFilePath == "" {
		panic(errors.New("APP_ENV_FILE environment variable not provided. It needs to be set with path to .app.env file"))
	}
	if err = godotenv.Load(appEnvsFilePath); err != nil {
		err = fmt.Errorf("unable to load environments from .app.env file %v", err)
		return
	}
	envs.app.botToken = os.Getenv("BOT_TOKEN")
	if envs.app.botToken == "" {
		err = errors.New("BOT_TOKEN environment variable not provided")
		return
	}
	envs.app.appId = os.Getenv("APPLICATION_ID")
	if envs.app.appId == "" {
		err = errors.New("APPLICATION_ID environment variable not provided")
		return
	}
	envs.runtimeEnvironment = os.Getenv("RUNTIME_ENVIRONMENT")
	if envs.runtimeEnvironment == "DEV" {
		envs.app.guildId = os.Getenv("GUILD_ID")
		if envs.app.guildId == "" {
			err = errors.New("GUILD_ID environment not provided. Its required in DEV run-time environment")
			return
		}
	}
	return
}
