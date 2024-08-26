package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"
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
Retrieve required environment variables for development / production run-time environments

Function panics when:
  - BOT_TOKEN is empty / not provided
  - APPLICATION_ID is empty / not provided
  - GUILD_ID is empty (only in "DEV" run-time environment, which is determined by RUNTIME_ENVIRONMENT env variable)
*/
func requiredEnvs() (botToken, applicationId, guildId string) {
	botToken = os.Getenv("BOT_TOKEN")
	if botToken == "" {
		panic(errors.New("BOT_TOKEN environment variable not provided"))
	}
	applicationId = os.Getenv("APPLICATION_ID")
	if applicationId == "" {
		panic(errors.New("APPLICATION_ID environment variable not provided"))
	}
	if os.Getenv("RUNTIME_ENVIRONMENT") == "DEV" {
		guildId = os.Getenv("GUILD_ID")
		log.Printf("guild is %s", guildId)
		if guildId == "" {
			panic(errors.New("GUILD_ID environment not provided. Its required in DEV run-time environment"))
		}
	}
	return
}
