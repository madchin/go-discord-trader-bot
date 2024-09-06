package service

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	errUserHaventOffers = errors.New("user do not have any offers")
	errOfferDoNotExists = errors.New("user dont have this offer")
)

func newServiceError(interaction *discordgo.Interaction, description string, previous error) error {
	return fmt.Errorf("for member %s in guild %s with msg %s: %w",
		interaction.Member.User.ID,
		interaction.GuildID,
		description,
		previous,
	)
}
