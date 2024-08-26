package service

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/offer"
)

type sell struct {
	notifier notifier
	storage  offer.Repository
}

func (s *sell) Add(ctx context.Context, interaction *discordgo.Interaction, offer offer.Offer) error {
	if err := s.storage.Add(offer); err != nil {
		failMsg := fmt.Sprintf("sell add for item %s, count %d, price %f failed", offer.Product().Name(), offer.Count(), offer.Product().Price())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail item add in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service add error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}
	successMsg := fmt.Sprintf("Item %s successfully added for sell offer! Can I help you with something else?", offer.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return fmt.Errorf("send follow up message for success message add in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
	}
	return nil
}

func (s *sell) Remove(ctx context.Context, interaction *discordgo.Interaction, offer offer.Offer) error {
	if err := s.storage.Remove(offer); err != nil {
		failMsg := fmt.Sprintf("sell remove item %s failed", offer.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail in item removal in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service remove error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}
	successMsg := fmt.Sprintf("Item %s successfully removed from sell offers! Need more help? Ask!", offer.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return fmt.Errorf("send follow up message for success item removal in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
	}
	return nil
}

func (s *sell) Update(ctx context.Context, interaction *discordgo.Interaction, offer offer.Offer) error {
	if err := s.storage.Update(offer); err != nil {
		failMsg := fmt.Sprintf("sell update item %s failed", offer.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail in item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service update error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}
	successMsg := fmt.Sprintf("Item %s successfully updated! Need more help? Ask!", offer.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return fmt.Errorf("send follow up message for success item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
	}
	return nil
}

func (s *sell) List(ctx context.Context, interaction *discordgo.Interaction, offer offer.Offer) error {
	offers, err := s.storage.List(offer.Product().Name())
	if err != nil {
		failMsg := fmt.Sprintf("sell retrieving item %s failed", offer.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail in item listing in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service list error for member %s in guild %s for product %s %w", interaction.Member.User.ID, interaction.GuildID, offer.Product().Name(), err)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, fmt.Sprintf("%v", offers)); err != nil {
		return fmt.Errorf("send follow up message for success item listing in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
	}
	return nil
}
