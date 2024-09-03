package service

import (
	followup "github.com/madchin/trader-bot/internal/domain/followup_message"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type Service struct {
	offer         *offerService
	itemRegistrar *itemRegistrar
}

func New(offerStorage offer.Repository, itemStorage item.Repository, notifier followup.MessageProducer) *Service {
	return &Service{
		offer: &offerService{
			offerStorage: offerStorage,
			notifier:     notifier,
		},
		itemRegistrar: &itemRegistrar{
			itemStorage: itemStorage,
		},
	}
}

func (s *Service) Offer() *offerService {
	return s.offer
}
