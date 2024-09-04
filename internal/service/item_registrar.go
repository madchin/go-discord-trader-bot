package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	followup "github.com/madchin/trader-bot/internal/domain/followup_message"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/gateway/command"
)

type commandRegistrar interface {
	RegisterAppCommand(appId, guildId string, cmd command.ApplicationCommand) error
}

type itemRegistrar struct {
	itemStorage      item.Repository
	notifier         messageProducer
	commandRegistrar commandRegistrar
}

func (itemRegistrar *itemRegistrar) Add(ctx context.Context, interaction *discordgo.Interaction, incomingItem item.Item) error {
	items, err := itemRegistrar.itemStorage.List(ctx)
	if err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailAdd(incomingItem.Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar add", err))
		}
		return newServiceError(ctx, interaction, "item registrar add", err)
	}
	if items.Contains(incomingItem) {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailAddItemAlreadyExist(incomingItem.Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar add", err))
		}
		return newServiceError(ctx, interaction, "item registrar add", errors.New("item already registered"))
	}
	choices, err := mapDomainItemsToChoices(items)
	//FIXME
	if err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailAdd(fmt.Sprintf("%s, MAX ITEM LIMIT REACHED", incomingItem.Name()))); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar add", err))
		}
		return newServiceError(ctx, interaction, "item registrar add, map domain items to command choices", err)
	}
	if err := itemRegistrar.itemStorage.Add(ctx, incomingItem); err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailAdd(incomingItem.Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar add", err))
		}
		return newServiceError(ctx, interaction, "item registrar add", err)
	}
	appId, guildId := os.Getenv("APPLICATION_ID"), os.Getenv("GUILD_ID")
	cmd := command.OfferBuilder(appId, guildId, choices).ApplicationCommand()
	if err := itemRegistrar.commandRegistrar.RegisterAppCommand(appId, guildId, cmd); err != nil {
		//FIXME
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailAdd("SOmething wrong happened during reigstering offer")); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar add", err))
		}
		return newServiceError(ctx, interaction, "item registrar add, register command with added item", err)
	}
	if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemSuccessAdd(incomingItem.Name())); err != nil {
		log.Print(newServiceError(ctx, interaction, "item registrar add", err))
	}
	return nil
}

func (itemRegistrar *itemRegistrar) Remove(ctx context.Context, interaction *discordgo.Interaction, incomingItem item.Item) error {
	item, err := itemRegistrar.itemStorage.ListByName(ctx, incomingItem)
	if err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailRemove(incomingItem.Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar remove", err))
		}
		return newServiceError(ctx, interaction, "item registrar remove", err)
	}
	if item.IsZero() {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailRemoveItemNotExist(incomingItem.Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar remove", err))
		}
		return newServiceError(ctx, interaction, "item registrar remove", errors.New("item not exists"))
	}
	if err := itemRegistrar.itemStorage.Remove(ctx, item); err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailRemove(incomingItem.Name())); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar remove", err))
		}
		return newServiceError(ctx, interaction, "item registrar remove", err)
	}
	items, err := itemRegistrar.itemStorage.List(ctx)
	if err != nil {
		return newServiceError(ctx, interaction, "item registrar remove, retrieve items for command registrar", err)
	}
	choices, err := mapDomainItemsToChoices(items)
	if err != nil {
		//FIXME
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailRemove("ERROR HAPPENED DURING REGISTERING COMMAND")); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar add", err))
		}
		return newServiceError(ctx, interaction, "item registrar remove, map items to choices for command registrar", err)
	}
	appId, guildId := os.Getenv("APPLICATION_ID"), os.Getenv("GUILD_ID")
	cmd := command.OfferBuilder(appId, guildId, choices).ApplicationCommand()
	if err := itemRegistrar.commandRegistrar.RegisterAppCommand(appId, guildId, cmd); err != nil {
		//FIXME
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailRemove("ERROR HAPPENED DURING REGISTERING COMMAND")); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar add", err))
		}
		return newServiceError(ctx, interaction, "item registrar remove, register new command without removed item", err)
	}
	if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemSuccessRemove(incomingItem.Name())); err != nil {
		log.Print(newServiceError(ctx, interaction, "item registrar add", err))
	}
	return nil
}

func (itemRegistrar *itemRegistrar) List(ctx context.Context, interaction *discordgo.Interaction) error {
	items, err := itemRegistrar.itemStorage.List(ctx)
	if err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailList("")); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar list", err))
		}
		return newServiceError(ctx, interaction, "item registrar list", err)
	}
	if items.AreEmpty() {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailItemsNotExist("")); err != nil {
			log.Print(newServiceError(ctx, interaction, "item registrar list", err))
		}
		return newServiceError(ctx, interaction, "item registrar list", errors.New("there is no items registered"))
	}
	if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemSuccessList(items.ToReadableMessage())); err != nil {
		log.Print(newServiceError(ctx, interaction, "item registrar list", err))
	}
	return nil
}

func mapDomainItemsToChoices(domainItems item.Items) (*command.Choices, error) {
	choices := command.NewChoices(len(domainItems) + 1)
	for i := 0; i < len(domainItems); i++ {
		if err := choices.AddNext(&discordgo.ApplicationCommandOptionChoice{
			Name:  domainItems[i].Name(),
			Value: domainItems[i].Name(),
		}); err != nil {
			return nil, fmt.Errorf("during addition next choice: %w", err)
		}
	}
	return choices, nil
}
