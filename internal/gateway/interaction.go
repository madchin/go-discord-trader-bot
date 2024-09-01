package gateway

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type eventData struct {
	command         string
	subCommand      string
	itemName        string
	itemCount       int
	itemPrice       float64
	updateItemPrice float64
	updateItemCount int
}

type InteractionData struct {
	interaction *discordgo.Interaction
	command     string
	subCommand  string
	offer       offer.Offer
	updateOffer offer.Offer
}

type Job interface {
	Command() string
	Subcommand() string
	Interaction() *discordgo.Interaction
	Offer() offer.Offer
	UpdateOffer() offer.Offer
}

func (i *InteractionData) Offer() offer.Offer {
	return i.offer
}

func (i *InteractionData) UpdateOffer() offer.Offer {
	return i.updateOffer
}

func (i *InteractionData) Interaction() *discordgo.Interaction {
	return i.interaction
}

func (i *InteractionData) Command() string {
	return i.command
}

func (i *InteractionData) Subcommand() string {
	return i.subCommand
}

func (e eventData) mapToOffer(vendorIdentity offer.VendorIdentity) offer.Offer {
	product := offer.NewProduct(e.itemName, e.itemPrice)
	return offer.NewOffer(vendorIdentity, product, e.itemCount)
}

func (e eventData) mapToUpdateOffer(vendorIdentity offer.VendorIdentity) offer.Offer {
	product := offer.NewProduct(e.itemName, e.updateItemPrice)
	return offer.NewOffer(vendorIdentity, product, e.updateItemCount)
}

func immediateInteractionRespond(s *discordgo.Session, interaction *discordgo.Interaction, responseContent string) error {
	immediateResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseContent,
		},
	}
	return s.InteractionRespond(interaction, immediateResponse)
}

func deferredInteractionRespond(s *discordgo.Session, interaction *discordgo.Interaction) error {
	deferredResp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}
	return s.InteractionRespond(interaction, deferredResp)
}

func getInteractionEventData(interactionEvent *discordgo.InteractionCreate) (*InteractionData, error) {
	cmdData := interactionEvent.ApplicationCommandData()
	eventData := &eventData{}
	getInteractionDataRecursive(cmdData.Options, eventData)
	vendorIdentity := offer.NewVendorIdentity(interactionEvent.Member.User.ID)
	off := eventData.mapToOffer(vendorIdentity)
	updateOff := eventData.mapToUpdateOffer(vendorIdentity)
	return &InteractionData{
		interactionEvent.Interaction,
		eventData.command,
		eventData.subCommand,
		off,
		updateOff,
	}, nil
}

func getInteractionDataRecursive(appCmdData []*discordgo.ApplicationCommandInteractionDataOption, data *eventData) {
	if appCmdData == nil {
		return
	}
	for _, d := range appCmdData {
		switch d.Name {
		case buyCmdDescriptor.name:
			fallthrough
		case sellCmdDescriptor.name:
			data.command = d.Name
		case AddSubCmdDescriptor.name:
			fallthrough
		case RemoveSubCmdDescriptor.name:
			fallthrough
		case UpdateCountSubCmdDescriptor.name:
			fallthrough
		case UpdatePriceSubCmdDescriptor.name:
			fallthrough
		case ListByProductNameSubCmdDescriptor.name:
			fallthrough
		case ListByVendorSubCmdDescriptor.name:
			data.subCommand = d.Name
		case itemDescriptor.name:
			data.itemName = d.StringValue()
		case countDescriptor.name:
			data.itemCount = int(d.IntValue())
		case priceDescriptor.name:
			data.itemPrice = d.FloatValue()
		case updatePriceDescriptor.name:
			data.updateItemPrice = d.FloatValue()
		case updateCountDescriptor.name:
			data.updateItemCount = int(d.IntValue())
		}
		getInteractionDataRecursive(d.Options, data)
	}
}
