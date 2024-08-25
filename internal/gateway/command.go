package gateway

import (
	"github.com/bwmarrin/discordgo"
)

type appCmdBuilder func(appId, guildId string) *discordgo.ApplicationCommand

/*
if guildId is empty, cmd is considered as global command instead of guild one
*/
var offerCmdBuilder appCmdBuilder = func(appId, guildId string) *discordgo.ApplicationCommand {
	var (
		minCount float64 = 1
		minPrice float64 = 0
	)

	var sellBuySubCommandOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "item",
			Description: "item name",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "count",
			Description: "count",
			MinValue:    &minCount,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        "price",
			Description: "price per each item",
			MinValue:    &minPrice,
		},
	}

	appCmd := &discordgo.ApplicationCommand{
		ApplicationID: appId,
		GuildID:       guildId,
		Name:          "offer",
		Description:   "Add, Remove, Update, Retrieve specified one or all Sell / Buy offers",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "sell",
				Description: "add a sell offer, update or remove existing one, retrieve specified one or all existing",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "add",
						Description: "Add a sell offer. Example: /offer sell add itemName, priceForEach, count",
						Options:     sellBuySubCommandOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "update",
						Description: "Update sell offer. Example: /offer sell update itemName, priceForEach, count",
						Options:     sellBuySubCommandOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "remove",
						Description: "Remove a sell offer Example: /offer sell remove itemName, priceForEach, count",
						Options:     sellBuySubCommandOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "list",
						Description: "Retrieve specified or all sell offers Example: /offer list itemName",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "item",
								Description: "specified item for listing",
							},
						},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "buy",
				Description: "add a buy offer, update or remove existing one, retrieve specified one or all existing",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "add",
						Description: "Add a buy offer. Example: /offer buy add itemName, priceForEach, count",
						Options:     sellBuySubCommandOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "update",
						Description: "Update buy offer. Example: /offer buy update itemName, priceForEach, count",
						Options:     sellBuySubCommandOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "remove",
						Description: "Remove a buy offer Example: /offer buy remove itemName, priceForEach, count",
						Options:     sellBuySubCommandOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "list",
						Description: "Retrieve specified or all buy offers Example: /offer list itemName",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "item",
								Description: "specified item for listing",
							},
						},
					},
				},
			},
		},
	}
	return appCmd
}

func (g *gateway) registerAppCommand(appId, guildId string, cmd appCmdBuilder) error {
	if _, err := g.session.ApplicationCommandCreate(appId, guildId, cmd(appId, guildId)); err != nil {
		return err
	}
	return nil
}
