package gateway

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type gateway struct {
	session *discordgo.Session
}

func NewGatewaySession(botToken, appId, guildId string, dataEnqueuer DataEnqueuer) (*gateway, error) {
	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, fmt.Errorf("error during creating gateway session, %v", err)
	}
	gateway := &gateway{session}
	gateway.registerHandlers(dataEnqueuer)
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
