package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const choiceLimit = 25

var errLimitReached = func(limit int) error {
	return fmt.Errorf("choices limit has been reached. Max limit: %d", limit)
}

type Choices struct {
	c []*discordgo.ApplicationCommandOptionChoice
}

func NewChoices(size int) *Choices {
	return &Choices{make([]*discordgo.ApplicationCommandOptionChoice, 0, size)}
}

func (choices *Choices) AddNext(choice *discordgo.ApplicationCommandOptionChoice) error {
	if choices.isLimitReached(choiceLimit) {
		return errLimitReached(choiceLimit)
	}
	choices.c = append(choices.c, choice)
	return nil
}

func (choices *Choices) isLimitReached(limit int) (ok bool) {
	if len(choices.c) == limit {
		ok = true
	}
	return
}
