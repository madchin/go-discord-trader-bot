package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type metadata struct {
	command, subCommand, action string
	interaction                 *discordgo.Interaction
}

func (m metadata) Action() string {
	return m.action
}

func (m metadata) Command() string {
	return m.command
}

func (m metadata) Subcommand() string {
	return m.subCommand
}

func (m metadata) Interaction() *discordgo.Interaction {
	return m.interaction
}

type offerEvent struct {
	vendorOffer offer.VendorOffer
	updatePrice float64
}

func (o offerEvent) UpdatePrice() float64 {
	return o.updatePrice
}

func (o offerEvent) VendorOffer() offer.VendorOffer {
	return o.vendorOffer
}

type itemRegistrarEvent struct {
	item item.Item
}

func (i itemRegistrarEvent) Item() item.Item {
	return i.item
}

type Event struct {
	Metadata           *metadata
	OfferEvent         offerEvent
	ItemRegistrarEvent itemRegistrarEvent
}

// xD
func (e Event) Data() Event {
	return e
}
