package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/gateway/command"
)

type itemRegistrarEventData struct {
	name string
}

func (o itemRegistrarEventData) mapToItemRegistrarEvent() itemRegistrarEvent {
	return itemRegistrarEvent{item.New(o.name)}
}

func itemRegistrarData(appCmdData *discordgo.ApplicationCommandInteractionDataOption, itemRegistrar *itemRegistrarEventData) {
	if appCmdData.Name == command.ItemRegistrar.Option.Name.Descriptor() {
		itemRegistrar.name = appCmdData.StringValue()
	}
}

func isItemRegistrarSubCommand(candidate string) bool {
	isAdd := candidate == command.ItemRegistrar.SubCommand.Add.Descriptor()
	isRemove := candidate == command.ItemRegistrar.SubCommand.Remove.Descriptor()
	isList := candidate == command.ItemRegistrar.SubCommand.List.Descriptor()
	return isAdd || isRemove || isList
}
