package followup

import "github.com/bwmarrin/discordgo"

type MessageProducer interface {
	SendFollowUpMessage(interaction *discordgo.Interaction, message Message) error
}
