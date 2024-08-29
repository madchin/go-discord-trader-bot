package gateway

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type eventData struct {
	authorId        string
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

func (e eventData) mapToOffer() offer.Offer {
	product := offer.NewProduct(e.itemName, e.itemPrice)
	return offer.NewOffer(e.authorId, product, e.itemCount)
}

func (e eventData) mapToUpdateOffer() offer.Offer {
	product := offer.NewProduct(e.itemName, e.updateItemPrice)
	return offer.NewOffer(e.authorId, product, e.updateItemCount)
}

func immediateInteractionRespond(s *discordgo.Session, interaction *discordgo.Interaction, responseContent string) error {
	immediateResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:         responseContent,
			AllowedMentions: &discordgo.MessageAllowedMentions{RepliedUser: true},
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
	eventData := &eventData{authorId: interactionEvent.Member.User.ID}
	getInteractionDataRecursive(cmdData.Options, eventData)
	fmt.Printf(":event data is %v", cmdData)
	off := eventData.mapToOffer()
	updateOff := eventData.mapToUpdateOffer()
	return &InteractionData{interactionEvent.Interaction, eventData.command, eventData.subCommand, off, updateOff}, nil
}

func getInteractionDataRecursive(appCmdData []*discordgo.ApplicationCommandInteractionDataOption, data *eventData) {
	if appCmdData == nil {
		return
	}

	for _, d := range appCmdData {
		if d != nil {
			fmt.Printf("app cmd data is %v", *d)
		}
		switch d.Name {
		case BuyCmdDescriptor.name:
			fallthrough
		case SellCmdDescriptor.name:
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
