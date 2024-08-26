package gateway

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type scheduler interface {
	Schedule(data *InteractionData)
}

type handler struct {
	scheduler scheduler
}

func (h *handler) onReady(s *discordgo.Session, _ *discordgo.Ready) {
	log.Printf("websocket is ready for interaction \nbot user: %s", s.State.User.Username)
}

func (h *handler) onMessageSend(s *discordgo.Session, msg *discordgo.MessageCreate) {
	log.Printf("message received\nAuthor: %s\nContent: %s GID: %s", msg.Author, msg.Content, msg.GuildID)
}

func (h *handler) onUserInteraction(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	interactionData, err := getInteractionEventData(interaction)
	if err != nil {
		log.Printf("error occured during mapping item to offer. %v", err)
		if err := immediateInteractionRespond(s, interaction.Interaction, err.Error()); err != nil {
			log.Printf("error during immediate responding to command invoker. %v", err)
		}
	}
	if err := deferredInteractionRespond(s, interaction.Interaction); err != nil {
		log.Printf("error occured during responding to interaction %v", err)
	}
	log.Printf("interaction event %v", *interaction.Interaction)

	h.scheduler.Schedule(interactionData)
}
