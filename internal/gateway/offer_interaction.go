package gateway

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type eventOffer struct {
	item  string
	count int
	price float64
}

type offerEventData struct {
	subCommand, action string
	offer              eventOffer
}

func (o offerEventData) mapToDomainOffer() offer.Offer {
	return offer.NewOffer(offer.NewProduct(o.offer.item, o.offer.price), o.offer.count)
}

func offerData(appCmdData *discordgo.ApplicationCommandInteractionDataOption, data *offerEventData) {
	var isSubCommand = func(candidate string) bool {
		return candidate == buySubCmdDescriptor.name || candidate == sellSubCmdDescriptor.name
	}
	var isAction = func(candidate string) (ok bool) {
		switch candidate {
		case AddActionDescriptor.name:
			fallthrough
		case RemoveActionDescriptor.name:
			fallthrough
		case UpdateCountActionDescriptor.name:
			fallthrough
		case UpdatePriceActionDescriptor.name:
			fallthrough
		case ListByProductNameActionDescriptor.name:
			fallthrough
		case ListByVendorActionDescriptor.name:
			ok = true
		}
		return
	}

	if isSubCommand(appCmdData.Name) {
		data.subCommand = appCmdData.Name
	}
	if isAction(appCmdData.Name) {
		data.action = appCmdData.Name
	}
	switch appCmdData.Name {
	case itemDescriptor.name:
		data.offer.item = appCmdData.StringValue()
	case countDescriptor.name:
		fallthrough
	case updateCountDescriptor.name:
		data.offer.count = int(appCmdData.IntValue())
	case priceDescriptor.name:
		fallthrough
	case updatePriceDescriptor.name:
		data.offer.price = appCmdData.FloatValue()
	}
}
