package command

import "github.com/bwmarrin/discordgo"

type offerCommandFunc func(appId, guildId string, choices *Choices) offerCommand
type itemRegistrarCommandFunc func(appId, guildId string) itemRegistrarCommand

// each descriptor must have unique name in order to be correctly parsed in interaction data retriever
type (
	descriptor struct{ name string }
	option     struct{ descriptor }
	action     struct{ descriptor }
	command    struct{ descriptor }
	subCommand struct{ descriptor }
)

func (d descriptor) Descriptor() string {
	return d.name
}

type ApplicationCommand struct {
	command *discordgo.ApplicationCommand
}

func (appCmd ApplicationCommand) Raw() *discordgo.ApplicationCommand {
	return appCmd.command
}

type offerCommand struct {
	appCmd ApplicationCommand
}

func (offerCmd offerCommand) ApplicationCommand() ApplicationCommand {
	return offerCmd.appCmd
}

type itemRegistrarCommand struct {
	appCmd ApplicationCommand
}

func (itemRegistrarCmd itemRegistrarCommand) ApplicationCommand() ApplicationCommand {
	return itemRegistrarCmd.appCmd
}
