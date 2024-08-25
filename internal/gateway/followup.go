package gateway

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (g *gateway) SendFollowUpMessage(interaction *discordgo.Interaction, data *discordgo.WebhookParams) error {
	msg, err := g.session.FollowupMessageCreate(interaction, true, data)
	if err != nil {
		log.Printf("followup message err %v", err)
		return err
	}
	log.Printf("sent followup message for interaction %v\nmsg: %v", interaction, msg)
	return nil
}
