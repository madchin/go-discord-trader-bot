package service

import (
	followup "github.com/madchin/trader-bot/internal/domain/followup_message"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type Service struct {
	offer *offerService
}

func New(offerStorage offer.Repository, notifier followup.MessageProducer) *Service {
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
