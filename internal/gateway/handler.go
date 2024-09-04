package gateway

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/gateway/interaction"
)

type scheduler interface {
	Schedule(data Job)
}

type Job interface {
	Data() interaction.Event
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

func (h *handler) onUserInteraction(s *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	event, err := interaction.EventData(interactionCreate)
	if err != nil {
		log.Printf("error occured during mapping item to offer. %v", err)
		if err := interaction.ImmediateRespond(s, interactionCreate.Interaction, err.Error()); err != nil {
			log.Printf("error during immediate responding to command invoker. %v", err)
		}
	}
	if err := interaction.DeferredRespond(s, interactionCreate.Interaction); err != nil {
		log.Printf("error occured during responding to interaction %v", err)
	}
	log.Printf("interaction event %v", *interactionCreate.Interaction)

	h.scheduler.Schedule(event)
}
