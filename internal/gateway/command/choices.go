package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/item"
)

const ChoicesLimit = 25

var errLimitReached = func(limit int) error {
	return fmt.Errorf("choices limit has been reached. Max limit: %d", limit)
}

type Choices struct {
	c []*discordgo.ApplicationCommandOptionChoice
}

func NewChoices(items item.Items) (*Choices, error) {
	if isLimitExceeded(len(items)) {
		return nil, errLimitReached(ChoicesLimit)
	}
	choices := &Choices{make([]*discordgo.ApplicationCommandOptionChoice, len(items))}
	choices.seedWithItems(items)
	return choices, nil
}

func (choices *Choices) seedWithItems(items item.Items) {
	for i := 0; i < len(choices.c); i++ {
		choice := &discordgo.ApplicationCommandOptionChoice{
			Name:  items[i].Name(),
			Value: items[i].Name(),
		}
		choices.c[i] = choice
	}
}

func isLimitExceeded(size int) (ok bool) {
	if size > ChoicesLimit {
		ok = true
	}
	return
}
