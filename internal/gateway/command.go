package gateway

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

//FIXME list-mine to not accept any param and display all data for requesting user
//FIXME displaying list data to user

var (
	itemDescriptor        = option{descriptor{"item"}}
	countDescriptor       = option{descriptor{"count"}}
	priceDescriptor       = option{descriptor{"price"}}
	updateCountDescriptor = option{descriptor{"count-update"}}
	updatePriceDescriptor = option{descriptor{"price-update"}}

	buySubCmdDescriptor  = subCommand{descriptor{"buy"}}
	sellSubCmdDescriptor = subCommand{descriptor{"sell"}}

	OfferCommandDescriptor = command{descriptor{"offer"}}

	AddActionDescriptor               = action{descriptor{"add"}}
	RemoveActionDescriptor            = action{descriptor{"remove"}}
	UpdateCountActionDescriptor       = action{descriptor{"update-count"}}
	UpdatePriceActionDescriptor       = action{descriptor{"update-price"}}
	ListByProductNameActionDescriptor = action{descriptor{"list-for-name"}}
	ListByVendorActionDescriptor      = action{descriptor{"list-mine"}}
)

type appCmd func(appId, guildId string) *discordgo.ApplicationCommand

// each descriptor must have unique name in order to be correctly parsed in interaction data retriever
type descriptor struct {
	name string
}

type option struct {
	descriptor
}

type action struct {
	descriptor
}

type command struct {
	descriptor
}

type subCommand struct {
	descriptor
}

func (d descriptor) Descriptor() string {
	return d.name
}

/*
if guildId is empty, cmd is considered as global command instead of guild one
*/
var offerCommand appCmd = func(appId, guildId string) *discordgo.ApplicationCommand {
	var (
		minCount float64 = float64(offer.MinCount)
	)

	var addOptions = []*discordgo.ApplicationCommandOption{
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
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
		},
	}

	var updatePriceOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        itemDescriptor.name,
			Description: "item name",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        priceDescriptor.name,
			Description: "price per each item",
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        updatePriceDescriptor.name,
			Description: "update price",
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
			Required:    true,
		},
	}

	var updateCountOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        itemDescriptor.name,
			Description: "item name",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        priceDescriptor.name,
			Description: "price per each item",
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        updateCountDescriptor.name,
			Description: "update count",
			MinValue:    &minCount,
			Required:    true,
		},
	}

	var removeOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        itemDescriptor.name,
			Description: "item name",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        priceDescriptor.name,
			Description: "price per each item",
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
		},
	}

	var listByNameOffers = []*discordgo.ApplicationCommandOption{
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
		Description:   "Manage sell / buy offers with subcommands as add, update, retrieve, remove. ",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        sellSubCmdDescriptor.name,
				Description: "Add, update, retrieve or remove sell offers",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        AddActionDescriptor.name,
						Description: "Add a sell offer. If offer exists (product and price is the same), its count will be updated.",
						Options:     addOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        UpdateCountActionDescriptor.name,
						Description: "Update sell offer product count.",
						Options:     updateCountOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        UpdatePriceActionDescriptor.name,
						Description: "Update sell offer product price.",
						Options:     updatePriceOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        RemoveActionDescriptor.name,
						Description: "Remove a sell offer completely.",
						Options:     removeOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        ListByProductNameActionDescriptor.name,
						Description: "Retrieve specified by name offer.",
						Options:     listByNameOffers,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        ListByVendorActionDescriptor.name,
						Description: "List all your offers.",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        buySubCmdDescriptor.name,
				Description: "add a buy offer, update or remove existing one, retrieve specified one or all existing",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        AddActionDescriptor.name,
						Description: "Add a buy offer. If offer exists (product and price is the same), its count will be updated",
						Options:     addOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        UpdateCountActionDescriptor.name,
						Description: "Update buy offer product count.",
						Options:     updateCountOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        UpdatePriceActionDescriptor.name,
						Description: "Update buy offer product price.",
						Options:     updatePriceOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        RemoveActionDescriptor.name,
						Description: "Remove buy offer completely.",
						Options:     removeOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        ListByProductNameActionDescriptor.name,
						Description: "Retrieve specified by name buy offers.",
						Options:     listByNameOffers,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        ListByVendorActionDescriptor.name,
						Description: "List all your offers.",
					},
				},
			},
		},
	}
	return appCmd
}
