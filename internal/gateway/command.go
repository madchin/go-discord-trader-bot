package gateway

import (
	"github.com/bwmarrin/discordgo"
)

//FIXME list-mine to not accept any param and display all data for requesting user
//FIXME displaying list data to user

var (
	itemDescriptor                    = descriptor{"item"}
	countDescriptor                   = descriptor{"count"}
	priceDescriptor                   = descriptor{"price"}
	updateCountDescriptor             = descriptor{"count-update"}
	updatePriceDescriptor             = descriptor{"price-update"}
	BuyCmdDescriptor                  = descriptor{"buy"}
	SellCmdDescriptor                 = descriptor{"sell"}
	AddSubCmdDescriptor               = descriptor{"add"}
	RemoveSubCmdDescriptor            = descriptor{"remove"}
	UpdateCountSubCmdDescriptor       = descriptor{"update-count"}
	UpdatePriceSubCmdDescriptor       = descriptor{"update-price"}
	ListByProductNameSubCmdDescriptor = descriptor{"list-for-name"}
	ListByVendorSubCmdDescriptor      = descriptor{"list-mine"}
)

type appCmdBuilder func(appId, guildId string) *discordgo.ApplicationCommand

type descriptor struct {
	name string
}

func (d descriptor) Descriptor() string {
	return d.name
}

/*
if guildId is empty, cmd is considered as global command instead of guild one
*/
var offerCmdBuilder appCmdBuilder = func(appId, guildId string) *discordgo.ApplicationCommand {
	var (
		minCount float64 = 1
		minPrice float64 = 0
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
			MinValue:    &minPrice,
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
			MinValue:    &minPrice,
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        updatePriceDescriptor.name,
			Description: "update price",
			MinValue:    &minPrice,
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
			MinValue:    &minPrice,
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
			MinValue:    &minPrice,
		},
	}

	var listOffersOptions = []*discordgo.ApplicationCommandOption{
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
				Name:        SellCmdDescriptor.name,
				Description: "add a sell offer, update or remove existing one, retrieve specified one or all existing",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        AddSubCmdDescriptor.name,
						Description: "Add a sell offer. Example: /offer sell add itemName, priceForEach, count",
						Options:     addOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        UpdateCountSubCmdDescriptor.name,
						Description: "Update sell offer product count. Example: /offer sell update itemName, priceForEach, count",
						Options:     updateCountOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        UpdatePriceSubCmdDescriptor.name,
						Description: "Update sell offer product price. Example: /offer sell update itemName, priceForEach, count",
						Options:     updatePriceOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        RemoveSubCmdDescriptor.name,
						Description: "Remove a sell offer Example: /offer sell remove itemName, priceForEach, count",
						Options:     removeOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        ListByProductNameSubCmdDescriptor.name,
						Description: "Retrieve specified or all sell offers Example: /offer list itemName",
						Options:     listOffersOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        ListByVendorSubCmdDescriptor.name,
						Description: "Retrieve specified or all your sell offers Example: /offer list itemName",
						Options:     listOffersOptions,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        BuyCmdDescriptor.name,
				Description: "add a buy offer, update or remove existing one, retrieve specified one or all existing",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        AddSubCmdDescriptor.name,
						Description: "Add a buy offer. Example: /offer buy add itemName, priceForEach, count",
						Options:     addOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        UpdateCountSubCmdDescriptor.name,
						Description: "Update buy offer product count. Example: /offer sell update itemName, priceForEach, count",
						Options:     updateCountOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        UpdatePriceSubCmdDescriptor.name,
						Description: "Update buy offer product price. Example: /offer sell update itemName, priceForEach, count",
						Options:     updatePriceOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        RemoveSubCmdDescriptor.name,
						Description: "Remove a buy offer Example: /offer buy remove itemName, priceForEach, count",
						Options:     removeOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        ListByProductNameSubCmdDescriptor.name,
						Description: "Retrieve specified or all buy offers Example: /offer list itemName",
						Options:     listOffersOptions,
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        ListByVendorSubCmdDescriptor.name,
						Description: "Retrieve specified or all your buy offers Example: /offer list itemName",
						Options:     listOffersOptions,
					},
				},
			},
		},
	}
	return appCmd
}
