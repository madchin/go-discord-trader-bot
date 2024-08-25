package gateway

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type data struct {
	name  string
	value interface{}
}

type onUserInteractionHandler func(s *discordgo.Session, interaction *discordgo.InteractionCreate)

type InteractionData map[string][]data

type DataEnqueuer interface {
	Enqueue(data InteractionData) error
}

func onUserInteraction(dataEnqueuer DataEnqueuer) onUserInteractionHandler {
	return onUserInteractionHandler(func(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
		deferredResp := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		}
		if err := s.InteractionRespond(interaction.Interaction, deferredResp); err != nil {
			log.Printf("error occured during responding to interaction %v", err)
		}
		interactionData := getInteractionData(interaction.ApplicationCommandData())
		if err := dataEnqueuer.Enqueue(interactionData); err != nil {
			log.Printf("error occured during enqueueing data for interaction: %v", interactionData)
		}
	})
}

func getInteractionData(appCmdData discordgo.ApplicationCommandInteractionData) InteractionData {
	data := make([]data, 0, 3)
	getInteractionDataRecursive(appCmdData.Options, &data)
	return InteractionData{appCmdData.ID: data}
}

func getInteractionDataRecursive(appCmdData []*discordgo.ApplicationCommandInteractionDataOption, interactionData *[]data) {
	if appCmdData == nil {
		return
	}
	for _, d := range appCmdData {
		if d.Value != nil {
			*interactionData = append(*interactionData, data{d.Name, d.Value})
		}
		getInteractionDataRecursive(d.Options, interactionData)
	}
}

func onReady(s *discordgo.Session, _ *discordgo.Ready) {
	log.Printf("websocket is ready for interaction \nbot user: %s", s.State.User.Username)
}

func onMessageSend(s *discordgo.Session, msg *discordgo.MessageCreate) {
	log.Printf("message received\nAuthor: %s\nContent: %s GID: %s", msg.Author, msg.Content, msg.GuildID)
}

func (g *gateway) registerHandlers(dataEnqueuer DataEnqueuer) {
	g.session.AddHandler(onReady)
	g.session.AddHandler(onMessageSend)
	g.session.AddHandler(onUserInteraction(dataEnqueuer))
}
