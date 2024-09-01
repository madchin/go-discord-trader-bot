package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/storage"
)

var (
	errUserHaventOffers = errors.New("user do not have any offers")
	errOfferDoNotExists = errors.New("user dont have this offer")
)

func newServiceError(ctx context.Context, interaction *discordgo.Interaction, description string, previous error) error {
	return fmt.Errorf("in %s for member %s in guild %s with msg %s: %w",
		ctx.Value(storage.CtxBuySellDbTableDescriptorKey),
		interaction.Member.User.ID,
		interaction.GuildID,
		description,
		previous,
	)
}
