package gateway

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type metadata struct {
	command, subCommand, action string
}

type event struct {
	metadata metadata
	offer    offer.VendorOffer
}

type InteractionData struct {
	interaction *discordgo.Interaction
	event       event
}

type Job interface {
	Metadata() metadata
	Interaction() *discordgo.Interaction
	VendorOffer() offer.VendorOffer
}

func (i *InteractionData) VendorOffer() offer.VendorOffer {
	return i.event.offer
}

func (i *InteractionData) Interaction() *discordgo.Interaction {
	return i.interaction
}

func (i *InteractionData) Metadata() metadata {
	return i.event.metadata
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
	commandData := interactionEvent.ApplicationCommandData()
	offerEvent, meta := &offerEventData{}, metadata{command: commandData.Name}
	switch commandData.Name {
	case OfferCommandDescriptor.name:
		getInteractionDataRecursive(commandData.Options, offerEvent)
		meta.subCommand, meta.action = offerEvent.subCommand, offerEvent.action
	default:
		return nil, fmt.Errorf("command data is unknown, need to be registered")
	}

	return &InteractionData{
		interactionEvent.Interaction,
		event{
			meta,
			offer.NewVendorOffer(offer.NewVendorIdentity(interactionEvent.Member.User.ID), offerEvent.mapToDomainOffer()),
		},
	}, nil
}

func getInteractionDataRecursive(appCmdData []*discordgo.ApplicationCommandInteractionDataOption, data any) {
	if appCmdData == nil {
		return
	}
	for _, d := range appCmdData {
		switch t := data.(type) {
		case *offerEventData:
			offerData(d, t)
		}
		getInteractionDataRecursive(d.Options, data)
	}
}
