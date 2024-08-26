package service

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/offer"
)

type Service struct {
	offer *off
}

type off struct {
	buy  *buy
	sell *sell
}

type notifier interface {
	SendFollowUpMessage(interaction *discordgo.Interaction, content string) error
}

func New(storage offer.Repository, notifier notifier) *Service {
	return &Service{
		offer: &off{
			buy: &buy{
				storage:  storage,
				notifier: notifier,
			},
			sell: &sell{
				storage:  storage,
				notifier: notifier,
			},
		},
	}
}

func (s *Service) BuyService() *buy {
	return s.offer.buy
}

func (s *Service) SellService() *sell {
	return s.offer.sell
}
