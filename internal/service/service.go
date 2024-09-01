package service

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type Service struct {
	offer *offerService
}

type notifier interface {
	SendFollowUpMessage(interaction *discordgo.Interaction, content string) error
}

func New(offerStorage offer.Repository, notifier notifier) *Service {
	return &Service{
		offer: &offerService{
			offerStorage: offerStorage,
			notifier:     notifier,
		},
	}
}

func (s *Service) Offer() *offerService {
	return s.offer
}
