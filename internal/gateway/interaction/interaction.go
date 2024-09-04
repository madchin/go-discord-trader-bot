package interaction

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/gateway/command"
)

func ImmediateRespond(s *discordgo.Session, interaction *discordgo.Interaction, responseContent string) error {
	immediateResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseContent,
		},
	}
	return s.InteractionRespond(interaction, immediateResponse)
}

func DeferredRespond(s *discordgo.Session, interaction *discordgo.Interaction) error {
	deferredResp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}
	return s.InteractionRespond(interaction, deferredResp)
}

func EventData(interactionEvent *discordgo.InteractionCreate) (*Event, error) {
	incomingCommand := interactionEvent.ApplicationCommandData()
	metadata := &metadata{command: incomingCommand.Name, interaction: interactionEvent.Interaction}
	event := &Event{Metadata: metadata}

	switch incomingCommand.Name {
	case command.Offer.Command.Descriptor():
		offerEventData := &offerEventData{}
		eventDataRecursive(incomingCommand.Options, metadata, offerEventData)
		event.OfferEvent = offerEventData.mapToOfferEvent(interactionEvent.Member.User.ID)
	case command.ItemRegistrar.Command.Descriptor():
		itemRegistrarEventData := &itemRegistrarEventData{}
		eventDataRecursive(incomingCommand.Options, metadata, itemRegistrarEventData)
		event.ItemRegistrarEvent = itemRegistrarEventData.mapToItemRegistrarEvent()
	default:
		return nil, fmt.Errorf("command data is unknown, need to be registered")
	}

	return event, nil
}

func eventDataRecursive(appCmdData []*discordgo.ApplicationCommandInteractionDataOption, metadata *metadata, data any) {
	if appCmdData == nil {
		return
	}
	for _, commandData := range appCmdData {
		switch typedData := data.(type) {
		case *offerEventData:
			if isOfferAction(commandData.Name) {
				metadata.action = commandData.Name
			}
			if isOfferSubCommand(commandData.Name) {
				metadata.subCommand = commandData.Name
			}
			offerData(commandData, typedData)
		case *itemRegistrarEventData:
			if isItemRegistrarSubCommand(commandData.Name) {
				metadata.subCommand = commandData.Name
			}
			itemRegistrarData(commandData, typedData)
		}
		eventDataRecursive(commandData.Options, metadata, data)
	}
}
