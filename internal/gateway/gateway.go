package gateway

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	followup "github.com/madchin/trader-bot/internal/domain/followup_message"
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
	if err := gateway.registerAppCommand(appId, guildId, offerCommand); err != nil {
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

func (g *gateway) SendFollowUpMessage(interaction *discordgo.Interaction, message followup.Message) error {
	data := &discordgo.WebhookParams{Content: message.Randomize().Content()}
	_, err := g.session.FollowupMessageCreate(interaction, true, data)
	if err != nil {
		return fmt.Errorf("gateway send followup message: %w", err)
	}
	return nil
}

func (g *gateway) registerHandlers(handler *handler) {
	g.session.AddHandler(handler.onReady)
	g.session.AddHandler(handler.onMessageSend)
	g.session.AddHandler(handler.onUserInteraction)
}

func (g *gateway) registerAppCommand(appId, guildId string, cmd appCmd) error {
	if _, err := g.session.ApplicationCommandCreate(appId, guildId, cmd(appId, guildId)); err != nil {
		return fmt.Errorf("gateway register app command: %w", err)
	}
	return nil
}
