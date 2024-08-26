package gateway

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/offer"
)

type eventData struct {
	authorId   string
	command    string
	subCommand string
	itemName   string
	itemCount  int
	itemPrice  float64
}

type InteractionData struct {
	interaction *discordgo.Interaction
	offer       offer.Offer
}

func (i *InteractionData) Offer() offer.Offer {
	return i.offer
}

func (i *InteractionData) Interaction() *discordgo.Interaction {
	return i.interaction
}

func (e eventData) mapToOffer() (offer.Offer, error) {
	offerType, err := offer.OfferTypeFromCandidate(e.command)
	if err != nil {
		return offer.Offer{}, err
	}
	actionType, err := offer.ActionTypeFromCandidate(e.subCommand)
	if err != nil {
		return offer.Offer{}, err
	}
	// item for listing dont need validation
	if e.subCommand == offer.List.Action() {
		return offer.NewListingOffer(offerType, actionType, e.itemName), nil
	}
	prod, err := offer.NewProduct(e.itemName, e.itemPrice)
	if err != nil {
		return offer.Offer{}, err
	}
	off, err := offer.NewOffer(offerType, actionType, e.authorId, prod, e.itemCount)
	if err != nil {
		return offer.Offer{}, err
	}
	return off, nil
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
	off, err := eventData.mapToOffer()
	if err != nil {
		return nil, err
	}
	log.Printf("offer is %v", off)
	log.Printf("interaction data is %s", cmdData.ID)
	return &InteractionData{interactionEvent.Interaction, off}, nil
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
		case addSubCmdDescriptor.name:
			fallthrough
		case removeSubCmdDescriptor.name:
			fallthrough
		case updateSubCmdDescriptor.name:
			fallthrough
		case listSubCmdDescriptor.name:
			data.subCommand = d.Name
		case itemDescriptor.name:
			data.itemName = d.StringValue()
		case countDescriptor.name:
			data.itemCount = int(d.IntValue())
		case priceDescriptor.name:
			data.itemPrice = d.FloatValue()
		}
		getInteractionDataRecursive(d.Options, data)
	}
}
