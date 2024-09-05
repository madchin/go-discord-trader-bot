package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const ChoicesLimit = 25

var errLimitReached = func(limit int) error {
	return fmt.Errorf("choices limit has been reached. Max limit: %d", limit)
}

type Choices struct {
	c []*discordgo.ApplicationCommandOptionChoice
}

func NewChoices(size int) (*Choices, error) {
	if isLimitReached(size) {
		return nil, errLimitReached(ChoicesLimit)
	}
	return &Choices{make([]*discordgo.ApplicationCommandOptionChoice, 0, size)}, nil
}

func (choices *Choices) AddNext(choice *discordgo.ApplicationCommandOptionChoice) error {
	choices.c = append(choices.c, choice)
	return nil
}

func isLimitReached(size int) (ok bool) {
	if size == ChoicesLimit {
		ok = true
	}
	return
}
