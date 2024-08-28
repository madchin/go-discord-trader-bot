package gateway

import (
	"github.com/bwmarrin/discordgo"
)

type appCmdBuilder func(appId, guildId string) *discordgo.ApplicationCommand

type descriptor struct {
	name string
}

var (
	itemDescriptor         = descriptor{"item"}
	countDescriptor        = descriptor{"count"}
	priceDescriptor        = descriptor{"price"}
	buyCmdDescriptor       = descriptor{"buy"}
	sellCmdDescriptor      = descriptor{"sell"}
	addSubCmdDescriptor    = descriptor{"add"}
	removeSubCmdDescriptor = descriptor{"remove"}
	updateSubCmdDescriptor = descriptor{"update"}
	listSubCmdDescriptor   = descriptor{"list"}
)

/*
if guildId is empty, cmd is considered as global command instead of guild one
*/
var offerCmdBuilder appCmdBuilder = func(appId, guildId string) *discordgo.ApplicationCommand {
	var (
		minCount float64 = 1
		minPrice float64 = 0
	)

	var removeUpdateAddOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        itemDescriptor.name,
			Description: "item name",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        countDescriptor.name,
			Description: "count",
			MinValue:    &minCount,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        priceDescriptor.name,
			Description: "price per each item",
			MinValue:    &minPrice,
		},
	}

	var listOffersByNameOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        itemDescriptor.name,
			Description: "item name",
			Required:    true,
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
				Name:        sellCmdDescriptor.name,
				Description: "add a sell offer, update or remove existing one, retrieve specified one or all existing",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        addSubCmdDescriptor.name,
						Description: "Add a sell offer. Example: /offer sell add itemName, priceForEach, count",
						Options:     removeUpdateAddOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        updateSubCmdDescriptor.name,
						Description: "Update sell offer. Example: /offer sell update itemName, priceForEach, count",
						Options:     removeUpdateAddOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        removeSubCmdDescriptor.name,
						Description: "Remove a sell offer Example: /offer sell remove itemName, priceForEach, count",
						Options:     removeUpdateAddOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        listSubCmdDescriptor.name,
						Description: "Retrieve specified or all sell offers Example: /offer list itemName",
						Options:     listOffersByNameOptions,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        buyCmdDescriptor.name,
				Description: "add a buy offer, update or remove existing one, retrieve specified one or all existing",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        addSubCmdDescriptor.name,
						Description: "Add a buy offer. Example: /offer buy add itemName, priceForEach, count",
						Options:     removeUpdateAddOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        updateSubCmdDescriptor.name,
						Description: "Update buy offer. Example: /offer buy update itemName, priceForEach, count",
						Options:     removeUpdateAddOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        removeSubCmdDescriptor.name,
						Description: "Remove a buy offer Example: /offer buy remove itemName, priceForEach, count",
						Options:     removeUpdateAddOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        listSubCmdDescriptor.name,
						Description: "Retrieve specified or all buy offers Example: /offer list itemName",
						Options:     listOffersByNameOptions,
					},
				},
			},
		},
	}
	return appCmd
}
