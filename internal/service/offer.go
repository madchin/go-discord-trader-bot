package service

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"
	followup "github.com/madchin/trader-bot/internal/domain/followup_message"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type offerService struct {
	notifier     followup.MessageProducer
	offerStorage offer.Repository
}

func (s *offerService) Add(ctx context.Context, interaction *discordgo.Interaction, vendorOffer offer.VendorOffer) error {
	offers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorOffer.VendorIdentity())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailAdd(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item add fail", err))
		}
		return newServiceError(ctx, interaction, "item add fail", err)
	}
	if offers.Contains(vendorOffer) {
		vendorOffer = offers.MergeSameOffers(vendorOffer)
		if err := s.offerStorage.UpdateCount(ctx, vendorOffer, offer.OnOfferCountUpdate); err != nil {
			if err := s.notifier.SendFollowUpMessage(interaction, followup.FailUpdateOnAdd(vendorOffer.Product().Name())); err != nil {
				log.Print(newServiceError(ctx, interaction, "item add, fail update existing offer", err))
			}
			return newServiceError(ctx, interaction, "item add, fail update existing offer", err)
		}
		if err := s.notifier.SendFollowUpMessage(interaction, followup.SuccessUpdateOnAdd(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item add, fail update existing offer", err))
		}
		return nil
	}
	if err := s.offerStorage.Add(ctx, vendorOffer, offer.OnVendorOfferAdd); err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailAdd(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item add, fail new add", err))
		}
		return newServiceError(ctx, interaction, "item add, fail new add", err)
	}

	if err := s.notifier.SendFollowUpMessage(interaction, followup.SuccessAdd(vendorOffer.Product().Name())); err != nil {
		log.Print(newServiceError(ctx, interaction, "item add, success new add", err))
	}
	return nil
}

func (s *offerService) Remove(ctx context.Context, interaction *discordgo.Interaction, vendorOffer offer.VendorOffer) error {
	offers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorOffer.VendorIdentity())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailRemove(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item remove fail", err))
		}
		return newServiceError(ctx, interaction, "item remove fail", err)
	}
	if offers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailRemoveOnNotHavingAnyOffers(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item remove, fail user dont have offers", err))
		}
		return newServiceError(ctx, interaction, "item remove", errUserHaventOffers)
	}
	if !offers.Contains(vendorOffer) {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailRemoveOfferNotExists(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item remove, fail user dont have requested offer", err))
		}
		return newServiceError(ctx, interaction, "item remove", errOfferDoNotExists)
	}
	if err := s.offerStorage.Remove(ctx, vendorOffer); err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailRemove(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item remove failed", err))
		}
		return newServiceError(ctx, interaction, "item remove failed", err)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.SuccessRemove(vendorOffer.Product().Name())); err != nil {
		log.Print(newServiceError(ctx, interaction, "item remove success", err))
	}
	return nil
}

func (s *offerService) UpdateCount(ctx context.Context, interaction *discordgo.Interaction, vendorOffer offer.VendorOffer) error {
	vendorOffers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorOffer.VendorIdentity())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailUpdate(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item update count fail", err))
		}
		return newServiceError(ctx, interaction, "item update count fail", err)
	}
	if vendorOffers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailUpdateOnNotHavingAnyOffers(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item update count fail, user dont have offers", err))
		}
		return newServiceError(ctx, interaction, "item update count fail", errUserHaventOffers)
	}
	if !vendorOffers.Contains(vendorOffer) {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailUpdateOfferNotExists(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item update count fail, user dont have requested offer", err))
		}
		return newServiceError(ctx, interaction, "item update count fail", errOfferDoNotExists)
	}
	if err := s.offerStorage.UpdateCount(ctx, vendorOffer, offer.OnOfferCountUpdate); err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailUpdate(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item update count fail", err))
		}
		return newServiceError(ctx, interaction, "item update count fail", err)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.SuccessUpdate(vendorOffer.Product().Name())); err != nil {
		log.Print(newServiceError(ctx, interaction, "item update count success", err))
	}
	return nil
}

func (s *offerService) UpdatePrice(ctx context.Context, interaction *discordgo.Interaction, vendorOffer offer.VendorOffer) error {
	vendorOffers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorOffer.VendorIdentity())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailUpdate(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item update price fail", err))
		}
		return newServiceError(ctx, interaction, "item update price fail", err)
	}
	if vendorOffers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailUpdateOnNotHavingAnyOffers(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item update price fail, user do not have offers", err))
		}
		return newServiceError(ctx, interaction, "item update price fail", errUserHaventOffers)
	}
	if !vendorOffers.Contains(vendorOffer) {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailUpdateOfferNotExists(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item update price fail, user dont have requested offer", err))
		}
		return newServiceError(ctx, interaction, "item update price fail", errUserHaventOffers)
	}
	if err := s.offerStorage.UpdatePrice(ctx, vendorOffer, offer.OnOfferPriceUpdate); err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailUpdate(vendorOffer.Product().Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item update price fail", err))
		}
		return newServiceError(ctx, interaction, "item update price fail", err)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.SuccessUpdate(vendorOffer.Product().Name())); err != nil {
		log.Print(newServiceError(ctx, interaction, "item update price success", err))
	}
	return nil
}

func (s *offerService) ListByVendor(ctx context.Context, interaction *discordgo.Interaction, vendorIdentity offer.VendorIdentity) error {
	offers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorIdentity)
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailListVendor("")); err != nil {
			log.Print(newServiceError(ctx, interaction, "item list by vendor fail, sending followup messaage", err))
		}
		return newServiceError(ctx, interaction, "item list by vendor fail", err)
	}
	if offers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailListVendorNotHavingAnyOffers("")); err != nil {
			log.Print(newServiceError(ctx, interaction, "item list by vendor fail", err))
		}
		return newServiceError(ctx, interaction, "item list by vendor fail, user dont have any offers", errUserHaventOffers)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.SuccessListVendor(offers.ToReadableMessage())); err != nil {
		log.Print(newServiceError(ctx, interaction, "item list by vendor success", err))
	}
	return nil
}

func (s *offerService) ListByProductName(ctx context.Context, interaction *discordgo.Interaction, productName string) error {
	offers, err := s.offerStorage.ListOffersByName(ctx, productName)
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailList(productName)); err != nil {
			log.Print(newServiceError(ctx, interaction, "item list by product name fail", err))
		}
		return newServiceError(ctx, interaction, "item list by product name fail", err)
	}
	if offers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.FailListNotHavingAnyOffers(productName)); err != nil {
			log.Print(newServiceError(ctx, interaction, "item list by product name fail", err))
		}
		return newServiceError(ctx, interaction, "item list by product name fail", errUserHaventOffers)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.SuccessList(offers.ToReadableMessage())); err != nil {
		log.Print(newServiceError(ctx, interaction, "item list by product name success", err))
	}
	return nil
}
