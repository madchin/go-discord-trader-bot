package gateway

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type gateway struct {
	session *discordgo.Session
}

func NewGatewaySession(botToken, appId, guildId string, scheduler scheduler) (*gateway, error) {
	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, fmt.Errorf("error during creating gateway session, %v", err)
	}
	gateway := &gateway{session}
	handler := &handler{scheduler}
	gateway.registerHandlers(handler)
	if err := gateway.registerAppCommand(appId, guildId, offerCmdBuilder); err != nil {
		return nil, err
	}
	return gateway, nil
}

func (g *gateway) CloseSession() {
	log.Printf("closing websocket")
	if err := g.session.Close(); err != nil {
		panic(fmt.Errorf("error during closing session, %v", err))
	}
}

func (g *gateway) OpenConnection() error {
	if err := g.session.Open(); err != nil {
		return fmt.Errorf("error during opening connection to gateway session, %v", err)
	}
	return nil
}

func (g *gateway) SendFollowUpMessage(interaction *discordgo.Interaction, content string) error {
	data := &discordgo.WebhookParams{Content: content}
	msg, err := g.session.FollowupMessageCreate(interaction, true, data)
	if err != nil {
		log.Printf("followup message err %v", err)
		return err
	}
	log.Printf("sent followup message for interaction %v\nmsg: %v", interaction, msg)
	return nil
}

func (g *gateway) registerHandlers(handler *handler) {
	g.session.AddHandler(handler.onReady)
	g.session.AddHandler(handler.onMessageSend)
	g.session.AddHandler(handler.onUserInteraction)
}

func (g *gateway) registerAppCommand(appId, guildId string, cmd appCmdBuilder) error {
	if _, err := g.session.ApplicationCommandCreate(appId, guildId, cmd(appId, guildId)); err != nil {
		return err
	}
	return nil
}
