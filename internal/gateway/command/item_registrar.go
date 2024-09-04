package command

import "github.com/bwmarrin/discordgo"

var ItemRegistrar = &itemRegistrarDescriptor{
	Command: &command{descriptor{"item-register"}},
	SubCommand: &itemRegistrarSubCommand{
		Add:    subCommand{descriptor{"add"}},
		Remove: subCommand{descriptor{"remove"}},
		List:   subCommand{descriptor{"list"}},
	},
	Option: &itemRegistrarOption{
		Name: option{descriptor{"name"}},
	},
}

type itemRegistrarDescriptor struct {
	Command    *command
	SubCommand *itemRegistrarSubCommand
	Option     *itemRegistrarOption
}

type itemRegistrarSubCommand struct {
	Add    subCommand
	Remove subCommand
	List   subCommand
}

type itemRegistrarOption struct {
	Name option
}

var ItemRegistrarBuilder itemRegistrarCommandFunc = func(appId, guildId string) itemRegistrarCommand {
	// disabled for all users except guild administrator
	var permissions int64 = 0
	var options = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        ItemRegistrar.Option.Name.name,
			Description: "item name used for register",
			Required:    true,
		},
	}
	var actions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        ItemRegistrar.SubCommand.Add.name,
			Description: "Register new market item",
			Options:     options,
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        ItemRegistrar.SubCommand.Remove.name,
			Description: "Remove registered market item with name",
			Options:     options,
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        ItemRegistrar.SubCommand.List.name,
			Description: "List all registered market items",
		},
	}
	itemRegistrar := itemRegistrarCommand{
		ApplicationCommand{
			&discordgo.ApplicationCommand{
				ApplicationID:            appId,
				GuildID:                  guildId,
				Name:                     ItemRegistrar.Command.name,
				Description:              "Bunch of commands used to register items for offer buy/sell commands",
				DefaultMemberPermissions: &permissions,
				Options:                  actions,
			},
		},
	}
	return itemRegistrar
}
