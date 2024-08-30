package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type offerService struct {
	notifier notifier
	storage  offer.Repository
}

func (s *offerService) Add(ctx context.Context, interaction *discordgo.Interaction, off offer.Offer) error {
	offers, _ := s.storage.ListVendorOffers(ctx, off.Vendor().Name())
	if offers.Contains(off) {
		off = offers.MergeSameOffers(off)
		if err := s.storage.UpdateCount(ctx, off, off.Count(), offer.OnOfferCountUpdate); err != nil {
			return fmt.Errorf("item add error, update count due to same item existing already offer %v %w", off, err)
		}
		successMsg := fmt.Sprintf("Item %s has been updated because you already have one with same price! Can I help you with something else?", off.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
			return fmt.Errorf("send follow up message for success message add in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
	}
	if err := s.storage.Add(ctx, off, offer.OnOfferAdd); err != nil {
		failMsg := fmt.Sprintf("sell add for item %s, count %d, price %f failed", off.Product().Name(), off.Count(), off.Product().Price())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail item add in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service add error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}

	successMsg := fmt.Sprintf("Item %s successfully added for sell offer! Can I help you with something else?", off.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return fmt.Errorf("send follow up message for success message add in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
	}
	return nil
}

func (s *offerService) Remove(ctx context.Context, interaction *discordgo.Interaction, offer offer.Offer) error {
	offers, _ := s.storage.ListVendorOffers(ctx, offer.Vendor().Name())
	if !offers.Contains(offer) {
		failMsg := "unable to remove offer because you do not have one you requested"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return fmt.Errorf("send follow up message for success item removal in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service remove error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, errors.New(failMsg))
	}
	if err := s.storage.Remove(ctx, offer); err != nil {
		failMsg := fmt.Sprintf("hey, we are sorry but sell remove item %s failed", offer.Product().Name())
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

func (s *offerService) UpdateCount(ctx context.Context, interaction *discordgo.Interaction, off offer.Offer, count int) error {
	vendorOffers, err := s.storage.ListVendorOffers(ctx, off.Vendor().Name())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, "hey man, i cant update offer because you dont have any"); err != nil {
			log.Printf("send follow up message for fail in item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service update error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}
	if !vendorOffers.Contains(off) {
		failMsg := "hey man, we cant update offer you requested because you dont have it."
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail in item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service update error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}
	if err := s.storage.UpdateCount(ctx, off, count, offer.OnOfferCountUpdate); err != nil {
		failMsg := fmt.Sprintf("hey, we are sorry but sell update for item %s failed", off.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail in item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service update error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}
	successMsg := fmt.Sprintf("Item %s successfully updated! Need more help? Ask!", off.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return fmt.Errorf("send follow up message for success item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
	}
	return nil
}

func (s *offerService) UpdatePrice(ctx context.Context, interaction *discordgo.Interaction, off offer.Offer, price float64) error {
	vendorOffers, err := s.storage.ListVendorOffers(ctx, off.Vendor().Name())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, "hey man, i cant update offer because you dont have any"); err != nil {
			log.Printf("send follow up message for fail in item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service update error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}
	if !vendorOffers.Contains(off) {
		failMsg := "hey man, we cant update offer you requested because you dont have it."
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail in item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service update error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}
	if err := s.storage.UpdatePrice(ctx, off, price, offer.OnOfferPriceUpdate); err != nil {
		failMsg := fmt.Sprintf("hey, we are sorry but sell update for item %s failed", off.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail in item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service update error for member %s in guild %s %w", interaction.Member.User.ID, interaction.GuildID, err)
	}
	successMsg := fmt.Sprintf("Item %s successfully updated! Need more help? Ask!", off.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return fmt.Errorf("send follow up message for success item update in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
	}
	return nil
}

func (s *offerService) ListByVendor(ctx context.Context, interaction *discordgo.Interaction, vendorName string) error {
	offers, err := s.storage.ListVendorOffers(ctx, vendorName)
	if err != nil {
		failMsg := "Seems that you do not have any items, maybe want to add some? Feel free to ask!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail in item listing in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service list error for member %s in guild %s for vendor %s %w", interaction.Member.User.ID, interaction.GuildID, vendorName, err)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, fmt.Sprintf("%v", offers)); err != nil {
		return fmt.Errorf("send follow up message for success item listing in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
	}
	return nil
}

func (s *offerService) ListByProductName(ctx context.Context, interaction *discordgo.Interaction, productName string) error {
	offers, err := s.storage.ListOffers(ctx, productName)
	if err != nil {
		failMsg := "Seems that items you wanted to retrieve do not exists, maybe want to add some? Feel free to ask!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			log.Printf("send follow up message for fail in item listing in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
		}
		return fmt.Errorf("sell service list error for member %s in guild %s for product %s %w", interaction.Member.User.ID, interaction.GuildID, productName, err)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, fmt.Sprintf("%v", offers)); err != nil {
		return fmt.Errorf("send follow up message for success item listing in sell service failed. For member %s in guild %s \nerr: %v", interaction.Member.User.ID, interaction.GuildID, err)
	}
	return nil
}
