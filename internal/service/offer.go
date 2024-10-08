package service

import (
	"context"
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/domain/offer"
	followup "github.com/madchin/trader-bot/internal/gateway/followup_message"
)

type offerService struct {
	notifier     messageProducer
	offerStorage offer.Repository
	itemStorage  item.Repository
}

func (s *offerService) Add(ctx context.Context, interaction *discordgo.Interaction, vendorOffer offer.VendorOffer) error {
	item, err := s.itemStorage.ListByName(ctx, item.New(vendorOffer.Product.Name()))
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailAddItemNotRegistered(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item add fail", err))
		}
		return newServiceError(interaction, "item add fail", err)
	}
	if item.IsZero() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailAddItemNotRegistered(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item add fail", err))
		}
		return newServiceError(interaction, "item add fail", errors.New("item not exists"))
	}
	offers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorOffer.VendorIdentity())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailAdd(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item add fail", err))
		}
		return newServiceError(interaction, "item add fail", err)
	}
	if offers.Contains(vendorOffer) {
		vendorOffer = offers.MergeSameOffers(vendorOffer)
		if err := s.offerStorage.UpdateCount(ctx, vendorOffer, offer.OnOfferCountUpdate); err != nil {
			if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailUpdateOnAdd(vendorOffer.Product.Name())); err != nil {
				log.Print(newServiceError(interaction, "item add, fail update existing offer", err))
			}
			return newServiceError(interaction, "item add, fail update existing offer", err)
		}
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferSuccessUpdateOnAdd(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item add, fail update existing offer", err))
		}
		return nil
	}
	if err := s.offerStorage.Add(ctx, vendorOffer, offer.OnVendorOfferAdd); err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailAdd(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item add, fail new add", err))
		}
		return newServiceError(interaction, "item add, fail new add", err)
	}

	if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferSuccessAdd(vendorOffer.Product.Name())); err != nil {
		log.Print(newServiceError(interaction, "item add, success new add", err))
	}
	return nil
}

func (s *offerService) Remove(ctx context.Context, interaction *discordgo.Interaction, vendorOffer offer.VendorOffer) error {
	offers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorOffer.VendorIdentity())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailRemove(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item remove fail", err))
		}
		return newServiceError(interaction, "item remove fail", err)
	}
	if offers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailRemoveOnNotHavingAnyOffers(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item remove, fail user dont have offers", err))
		}
		return newServiceError(interaction, "item remove", errUserHaventOffers)
	}
	if !offers.Contains(vendorOffer) {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailRemoveOfferNotExists(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item remove, fail user dont have requested offer", err))
		}
		return newServiceError(interaction, "item remove", errOfferDoNotExists)
	}
	if err := s.offerStorage.Remove(ctx, vendorOffer); err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailRemove(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item remove failed", err))
		}
		return newServiceError(interaction, "item remove failed", err)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferSuccessRemove(vendorOffer.Product.Name())); err != nil {
		log.Print(newServiceError(interaction, "item remove success", err))
	}
	return nil
}

func (s *offerService) UpdateCount(ctx context.Context, interaction *discordgo.Interaction, vendorOffer offer.VendorOffer) error {
	vendorOffers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorOffer.VendorIdentity())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailUpdate(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item update count fail", err))
		}
		return newServiceError(interaction, "item update count fail", err)
	}
	if vendorOffers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailUpdateOnNotHavingAnyOffers(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item update count fail, user dont have offers", err))
		}
		return newServiceError(interaction, "item update count fail", errUserHaventOffers)
	}
	if !vendorOffers.Contains(vendorOffer) {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailUpdateOfferNotExists(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item update count fail, user dont have requested offer", err))
		}
		return newServiceError(interaction, "item update count fail", errOfferDoNotExists)
	}
	if err := s.offerStorage.UpdateCount(ctx, vendorOffer, offer.OnOfferCountUpdate); err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailUpdate(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item update count fail", err))
		}
		return newServiceError(interaction, "item update count fail", err)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferSuccessUpdate(vendorOffer.Product.Name())); err != nil {
		log.Print(newServiceError(interaction, "item update count success", err))
	}
	return nil
}

func (s *offerService) UpdatePrice(ctx context.Context, interaction *discordgo.Interaction, vendorOffer offer.VendorOffer, updatePrice float64) error {
	vendorOffers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorOffer.VendorIdentity())
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailUpdate(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item update price fail", err))
		}
		return newServiceError(interaction, "item update price fail", err)
	}
	if vendorOffers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailUpdateOnNotHavingAnyOffers(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item update price fail, user do not have offers", err))
		}
		return newServiceError(interaction, "item update price fail", errUserHaventOffers)
	}
	if !vendorOffers.Contains(vendorOffer) {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailUpdateOfferNotExists(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item update price fail, user dont have requested offer", err))
		}
		return newServiceError(interaction, "item update price fail", errUserHaventOffers)
	}
	if err := s.offerStorage.UpdatePrice(ctx, vendorOffer, updatePrice, offer.OnOfferPriceUpdate); err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailUpdate(vendorOffer.Product.Name())); err != nil {
			log.Print(newServiceError(interaction, "item update price fail", err))
		}
		return newServiceError(interaction, "item update price fail", err)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferSuccessUpdate(vendorOffer.Product.Name())); err != nil {
		log.Print(newServiceError(interaction, "item update price success", err))
	}
	return nil
}

func (s *offerService) ListByVendor(ctx context.Context, interaction *discordgo.Interaction, vendorIdentity offer.VendorIdentity) error {
	offers, err := s.offerStorage.ListOffersByIdentity(ctx, vendorIdentity)
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailListVendor("")); err != nil {
			log.Print(newServiceError(interaction, "item list by vendor fail, sending followup messaage", err))
		}
		return newServiceError(interaction, "item list by vendor fail", err)
	}
	if offers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailListVendorNotHavingAnyOffers("")); err != nil {
			log.Print(newServiceError(interaction, "item list by vendor fail", err))
		}
		return newServiceError(interaction, "item list by vendor fail, user dont have any offers", errUserHaventOffers)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferSuccessListVendor(offers.ToReadableMessage())); err != nil {
		log.Print(newServiceError(interaction, "item list by vendor success", err))
	}
	return nil
}

func (s *offerService) ListByProductName(ctx context.Context, interaction *discordgo.Interaction, productName string) error {
	offers, err := s.offerStorage.ListOffersByName(ctx, productName)
	if err != nil {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailList(productName)); err != nil {
			log.Print(newServiceError(interaction, "item list by product name fail", err))
		}
		return newServiceError(interaction, "item list by product name fail", err)
	}
	if offers.NotExists() {
		if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferFailListNotHavingAnyOffers(productName)); err != nil {
			log.Print(newServiceError(interaction, "item list by product name fail", err))
		}
		return newServiceError(interaction, "item list by product name fail", errUserHaventOffers)
	}
	if err := s.notifier.SendFollowUpMessage(interaction, followup.OfferSuccessList(offers.ToReadableMessage())); err != nil {
		log.Print(newServiceError(interaction, "item list by product name success", err))
	}
	return nil
}
