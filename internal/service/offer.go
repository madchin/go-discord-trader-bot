package service

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type offerService struct {
	notifier     notifier
	offerStorage offer.Repository
}

func (s *offerService) Add(ctx context.Context, interaction *discordgo.Interaction, off offer.Offer) error {
	offers, err := s.offerStorage.ListVendorOffers(ctx, off.VendorIdentity())
	if err != nil {
		failMsg := "Oops, Something went wrong!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item add fail, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item add fail", err)
	}
	if offers.Contains(off) {
		off = offers.MergeSameOffers(off)
		if err := s.offerStorage.UpdateCount(ctx, off, off.Count(), offer.OnOfferCountUpdate); err != nil {
			return newServiceError(ctx, interaction, "item add, fail update existing offer", err)
		}
		successMsg := fmt.Sprintf("Item %s has been updated because you already have one with same price! Can I help you with something else?", off.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
			return newServiceError(ctx, interaction, "item add, fail update existing offer, sending followup message", err)
		}
		return nil
	}
	if err := s.offerStorage.Add(ctx, off, offer.OnOfferAdd); err != nil {
		failMsg := fmt.Sprintf("sell add for item %s, count %d, price %f failed", off.Product().Name(), off.Count(), off.Product().Price())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item add, fail new add", err)
		}
		return newServiceError(ctx, interaction, "item add, fail new add, sending followup message", err)
	}

	successMsg := fmt.Sprintf("Item %s successfully added for sell offer! Can I help you with something else?", off.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return newServiceError(ctx, interaction, "item add, success new add, sending followup message", err)
	}
	return nil
}

func (s *offerService) Remove(ctx context.Context, interaction *discordgo.Interaction, off offer.Offer) error {
	offers, err := s.offerStorage.ListVendorOffers(ctx, off.VendorIdentity())
	if err != nil {
		failMsg := "Oops, Something went wrong!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item remove fail, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item remove fail", err)
	}
	if offers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, "hey man, i cant update offer because you dont have any"); err != nil {
			return newServiceError(ctx, interaction, "item remove, fail user dont have offers, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item remove, sending followup message", errUserHaventOffers)
	}
	if !offers.Contains(off) {
		failMsg := "unable to remove offer because you do not have one you requested"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item remove, fail user dont have requested offer, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item remove", errOfferDoNotExists)
	}
	if err := s.offerStorage.Remove(ctx, off); err != nil {
		failMsg := fmt.Sprintf("hey, we are sorry but sell remove item %s failed", off.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item remove failed, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item remove failed", err)
	}
	successMsg := fmt.Sprintf("Item %s successfully removed from sell offers! Need more help? Ask!", off.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return newServiceError(ctx, interaction, "item remove success, sending followup message", err)
	}
	return nil
}

func (s *offerService) UpdateCount(ctx context.Context, interaction *discordgo.Interaction, off offer.Offer, count int) error {
	vendorOffers, err := s.offerStorage.ListVendorOffers(ctx, off.VendorIdentity())
	if err != nil {
		failMsg := "Oops, Something went wrong!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item update count fail, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item update count fail", err)
	}
	if vendorOffers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, "hey man, i cant update offer because you dont have any"); err != nil {
			return newServiceError(ctx, interaction, "item update count fail, user dont have offers, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item update count fail", errUserHaventOffers)
	}
	if !vendorOffers.Contains(off) {
		failMsg := "hey man, we cant update offer you requested because you dont have it."
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item update count fail, user dont have requested offer", err)
		}
		return newServiceError(ctx, interaction, "item update count fail", errOfferDoNotExists)
	}
	if err := s.offerStorage.UpdateCount(ctx, off, count, offer.OnOfferCountUpdate); err != nil {
		failMsg := fmt.Sprintf("hey, we are sorry but sell update for item %s failed", off.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item update count fail, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item update count fail", err)
	}
	successMsg := fmt.Sprintf("Item %s successfully updated! Need more help? Ask!", off.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return newServiceError(ctx, interaction, "item update count success, sending followup message", err)
	}
	return nil
}

func (s *offerService) UpdatePrice(ctx context.Context, interaction *discordgo.Interaction, off offer.Offer, price float64) error {
	vendorOffers, err := s.offerStorage.ListVendorOffers(ctx, off.VendorIdentity())
	if err != nil {
		failMsg := "Oops, Something went wrong!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item update price fail, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item update price fail", err)
	}
	if vendorOffers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, "hey man, i cant update offer because you dont have any"); err != nil {
			return newServiceError(ctx, interaction, "item update price fail, user do not have offers, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item update price fail", errUserHaventOffers)
	}
	if !vendorOffers.Contains(off) {
		failMsg := "hey man, we cant update offer you requested because you dont have it."
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item update price fail, user dont have requested offer, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item update price fail", errUserHaventOffers)
	}
	if err := s.offerStorage.UpdatePrice(ctx, off, price, offer.OnOfferPriceUpdate); err != nil {
		failMsg := fmt.Sprintf("hey, we are sorry but sell update for item %s failed", off.Product().Name())
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item update price fail, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item update price fail", err)
	}
	successMsg := fmt.Sprintf("Item %s successfully updated! Need more help? Ask!", off.Product().Name())
	if err := s.notifier.SendFollowUpMessage(interaction, successMsg); err != nil {
		return newServiceError(ctx, interaction, "item update price success, sending followup message", err)
	}
	return nil
}

func (s *offerService) ListByVendor(ctx context.Context, interaction *discordgo.Interaction, vendorIdentity offer.VendorIdentity) error {
	offers, err := s.offerStorage.ListVendorOffers(ctx, vendorIdentity)
	if err != nil {
		failMsg := "Oops, Something went wrong!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item list by vendor fail, sending followup messaage", err)
		}
		return newServiceError(ctx, interaction, "item list by vendor fail", err)
	}
	if offers.NotExists() {
		failMsg := "Seems that you do not have any items, maybe want to add some? Feel free to ask!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item list by vendor fail, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item list by vendor fail, user dont have any offers", errUserHaventOffers)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, offers.ToReadableMessage()); err != nil {
		return newServiceError(ctx, interaction, "item list by vendor success, sending followup message", err)
	}
	return nil
}

func (s *offerService) ListByProductName(ctx context.Context, interaction *discordgo.Interaction, productName string) error {
	offers, err := s.offerStorage.ListOffers(ctx, productName)
	if err != nil {
		failMsg := "Oops, Something went wrong!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item list by product name fail, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item list by product name fail", err)
	}
	if offers.NotExists() {
		failMsg := "Seems that you do not have any items, maybe want to add some? Feel free to ask!"
		if err := s.notifier.SendFollowUpMessage(interaction, failMsg); err != nil {
			return newServiceError(ctx, interaction, "item list by product name fail, sending followup message", err)
		}
		return newServiceError(ctx, interaction, "item list by product name fail", errUserHaventOffers)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, offers.ToReadableMessage()); err != nil {
		return newServiceError(ctx, interaction, "item list by product name success, sending followup message", err)
	}
	return nil
}
