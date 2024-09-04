package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
	"github.com/madchin/trader-bot/internal/gateway/command"
)

type offerEventData struct {
	item        string
	count       int
	price       float64
	updatePrice float64
}

func (o offerEventData) mapToOfferEvent(vendorIdentity string) offerEvent {
	return offerEvent{
		vendorOffer: offer.NewVendorOffer(offer.NewVendorIdentity(vendorIdentity), offer.NewOffer(offer.NewProduct(o.item, o.price), o.count)),
		updatePrice: o.updatePrice,
	}
}

func offerData(appCmdData *discordgo.ApplicationCommandInteractionDataOption, offer *offerEventData) {
	switch appCmdData.Name {
	case command.Offer.Option.Item.Descriptor():
		offer.item = appCmdData.StringValue()
	case command.Offer.Option.Count.Descriptor():
		fallthrough
	case command.Offer.Option.UpdateCount.Descriptor():
		offer.count = int(appCmdData.IntValue())
	case command.Offer.Option.Price.Descriptor():
		offer.price = appCmdData.FloatValue()
	case command.Offer.Option.UpdatePrice.Descriptor():
		offer.updatePrice = appCmdData.FloatValue()
	}
}

func isOfferSubCommand(candidate string) bool {
	return candidate == command.Offer.SubCommand.Buy.Descriptor() || candidate == command.Offer.SubCommand.Sell.Descriptor()
}
func isOfferAction(candidate string) (ok bool) {
	switch candidate {
	case command.Offer.Action.Add.Descriptor():
		fallthrough
	case command.Offer.Action.Remove.Descriptor():
		fallthrough
	case command.Offer.Action.UpdateCount.Descriptor():
		fallthrough
	case command.Offer.Action.ListByProductName.Descriptor():
		fallthrough
	case command.Offer.Action.ListByVendor.Descriptor():
		fallthrough
	case command.Offer.Action.UpdatePrice.Descriptor():
		ok = true
	}
	return
}
