package service

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/domain/offer"
	"github.com/madchin/trader-bot/internal/gateway/command"
	followup "github.com/madchin/trader-bot/internal/gateway/followup_message"
)

type messageProducer interface {
	SendFollowUpMessage(interaction *discordgo.Interaction, message followup.Message) error
}

type commandRegistrar interface {
	RegisterAppCommand(cmd command.ApplicationCommand) error
}

type botService interface {
	messageProducer
	commandRegistrar
}

type Service struct {
	Offer         *offerService
	ItemRegistrar *itemRegistrar
}

func New(offerStorage offer.Repository, itemStorage item.Repository, botService botService) *Service {
	return &Service{
		Offer: &offerService{
			offerStorage: offerStorage,
			itemStorage:  itemStorage,
			notifier:     botService.(messageProducer),
		},
		ItemRegistrar: &itemRegistrar{
			itemStorage:      itemStorage,
			notifier:         botService.(messageProducer),
			commandRegistrar: botService.(commandRegistrar),
		},
	}
}
