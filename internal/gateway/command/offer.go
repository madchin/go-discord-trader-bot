package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

var Offer = &offerDescriptor{
	Command: &command{descriptor{"offer"}},
	SubCommand: &offerSubCommand{
		subCommand{descriptor{"buy"}},
		subCommand{descriptor{"sell"}},
	},
	Action: &offerAction{
		Add:               action{descriptor{"add"}},
		Remove:            action{descriptor{"remove"}},
		UpdateCount:       action{descriptor{"update-count"}},
		UpdatePrice:       action{descriptor{"update-price"}},
		ListByProductName: action{descriptor{"list-for-name"}},
		ListByVendor:      action{descriptor{"list-mine"}},
	},
	Option: &offerOption{
		Item:        option{descriptor{"item"}},
		Count:       option{descriptor{"count"}},
		Price:       option{descriptor{"price"}},
		UpdateCount: option{descriptor{"count-update"}},
		UpdatePrice: option{descriptor{"price-update"}},
	},
}

type offerDescriptor struct {
	Command    *command
	SubCommand *offerSubCommand
	Option     *offerOption
	Action     *offerAction
}

type offerOption struct {
	Item        option
	Count       option
	Price       option
	UpdateCount option
	UpdatePrice option
}

type offerSubCommand struct {
	Buy  subCommand
	Sell subCommand
}

type offerAction struct {
	Add               action
	Remove            action
	UpdateCount       action
	UpdatePrice       action
	ListByProductName action
	ListByVendor      action
}

/*
if guildId is empty, cmd is considered as global command instead of guild one
*/
var OfferBuilder offerCommandFunc = func(appId, guildId string, choices *Choices) offerCommand {
	var (
		minCount float64 = float64(offer.MinCount)
	)

	var addOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        Offer.Option.Item.name,
			Description: "item name",
			Required:    true,
			Choices:     choices.c,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        Offer.Option.Count.name,
			Description: "count",
			Required:    true,
			MinValue:    &minCount,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        Offer.Option.Price.name,
			Description: "price per each item",
			Required:    true,
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
		},
	}

	var updatePriceOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        Offer.Option.Item.name,
			Description: "item name",
			Required:    true,
			Choices:     choices.c,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        Offer.Option.Price.name,
			Description: "price per each item",
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        Offer.Option.UpdatePrice.name,
			Description: "update price",
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
			Required:    true,
		},
	}

	var updateCountOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        Offer.Option.Item.name,
			Description: "item name",
			Required:    true,
			Choices:     choices.c,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        Offer.Option.Price.name,
			Description: "price per each item",
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        Offer.Option.UpdateCount.name,
			Description: "update count",
			MinValue:    &minCount,
			Required:    true,
		},
	}

	var removeOptions = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        Offer.Option.Item.name,
			Description: "item name",
			Required:    true,
			Choices:     choices.c,
		},
		{
			Type:        discordgo.ApplicationCommandOptionNumber,
			Name:        Offer.Option.Price.name,
			Description: "price per each item",
			Required:    true,
			MinValue:    &offer.MinPrice,
			MaxValue:    offer.MaxPrice,
		},
	}

	var listByNameOffers = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        Offer.Option.Item.name,
			Description: "item name",
			Required:    true,
			Choices:     choices.c,
		},
	}

	appCmd := offerCommand{
		ApplicationCommand{
			&discordgo.ApplicationCommand{
				ApplicationID: appId,
				GuildID:       guildId,
				Name:          "offer",
				Description:   "Manage sell / buy offers with subcommands as add, update, retrieve, remove. ",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
						Name:        Offer.SubCommand.Sell.name,
						Description: "Add, update, retrieve or remove sell offers",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.Add.name,
								Description: "Add a sell offer. If offer exists (product and price is the same), its count will be updated.",
								Options:     addOptions,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.UpdateCount.name,
								Description: "Update sell offer product count.",
								Options:     updateCountOptions,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.UpdatePrice.name,
								Description: "Update sell offer product price.",
								Options:     updatePriceOptions,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.Remove.name,
								Description: "Remove a sell offer completely.",
								Options:     removeOptions,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.ListByProductName.name,
								Description: "Retrieve specified by name offer.",
								Options:     listByNameOffers,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.ListByVendor.name,
								Description: "List all your offers.",
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
						Name:        Offer.SubCommand.Buy.name,
						Description: "add a buy offer, update or remove existing one, retrieve specified one or all existing",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.Add.name,
								Description: "Add a buy offer. If offer exists (product and price is the same), its count will be updated",
								Options:     addOptions,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.UpdateCount.name,
								Description: "Update buy offer product count.",
								Options:     updateCountOptions,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.UpdatePrice.name,
								Description: "Update buy offer product price.",
								Options:     updatePriceOptions,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.Remove.name,
								Description: "Remove buy offer completely.",
								Options:     removeOptions,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.ListByProductName.name,
								Description: "Retrieve specified by name buy offers.",
								Options:     listByNameOffers,
							},
							{
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Name:        Offer.Action.ListByVendor.name,
								Description: "List all your offers.",
							},
						},
					},
				},
			},
		},
	}
	return appCmd
}
