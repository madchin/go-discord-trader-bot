package service

import (
	"github.com/bwmarrin/discordgo"
	followup "github.com/madchin/trader-bot/internal/domain/followup_message"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type messageProducer interface {
	SendFollowUpMessage(interaction *discordgo.Interaction, message followup.Message) error
}

type Service struct {
	offer         *offerService
	itemRegistrar *itemRegistrar
}

func New(offerStorage offer.Repository, itemStorage item.Repository, commandRegistrar commandRegistrar, notifier messageProducer) *Service {
	return &Service{
		offer: &offerService{
			offerStorage: offerStorage,
			notifier:     notifier,
		},
		itemRegistrar: &itemRegistrar{
			itemStorage:      itemStorage,
			notifier:         notifier,
			commandRegistrar: commandRegistrar,
		},
	}
}

func (s *Service) Offer() *offerService {
	return s.offer
}

func (s *Service) ItemRegistrar() *itemRegistrar {
	return s.itemRegistrar
}
