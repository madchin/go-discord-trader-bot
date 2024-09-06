package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/gateway/command"
	followup "github.com/madchin/trader-bot/internal/gateway/followup_message"
)

type itemRegistrar struct {
	itemStorage      item.Repository
	notifier         messageProducer
	commandRegistrar commandRegistrar
}

func (itemRegistrar *itemRegistrar) Add(ctx context.Context, interaction *discordgo.Interaction, incomingItem item.Item) error {
	items, err := itemRegistrar.itemStorage.List(ctx)
	if err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailAdd(incomingItem.Name())); err != nil {
			log.Print(newServiceError(interaction, "item registrar add", err))
		}
		return newServiceError(interaction, "item registrar add", err)
	}
	if items.Contains(incomingItem) {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailAddItemAlreadyExist(incomingItem.Name())); err != nil {
			log.Print(newServiceError(interaction, "item registrar add", err))
		}
		return newServiceError(interaction, "item registrar add", errors.New("item already registered"))
	}
	choices, err := command.NewChoices(len(items) + 1)
	if err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemRegisterFailLimitExceeded(fmt.Sprintf("%d", command.ChoicesLimit))); err != nil {
			log.Print(newServiceError(interaction, "item registrar add", err))
		}
		return newServiceError(interaction, "item registrar add", err)
	}
	if err := itemRegistrar.itemStorage.Add(ctx, incomingItem); err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailAdd(incomingItem.Name())); err != nil {
			log.Print(newServiceError(interaction, "item registrar add", err))
		}
		return newServiceError(interaction, "item registrar add", err)
	}
	appId, guildId := interaction.AppID, interaction.GuildID
	cmd := command.OfferBuilder(appId, guildId, choices).ApplicationCommand()
	if err := itemRegistrar.commandRegistrar.RegisterAppCommand(appId, guildId, cmd); err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemRegisterFail(incomingItem.Name())); err != nil {
			log.Print(newServiceError(interaction, "item registrar add", err))
		}
		return newServiceError(interaction, "item registrar add", err)
	}
	if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemRegisterSuccess(incomingItem.Name())); err != nil {
		log.Print(newServiceError(interaction, "item registrar add", err))
	}
	return nil
}

func (itemRegistrar *itemRegistrar) Remove(ctx context.Context, interaction *discordgo.Interaction, incomingItem item.Item) error {
	items, err := itemRegistrar.itemStorage.List(ctx)
	if err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailRemove(incomingItem.Name())); err != nil {
			log.Print(newServiceError(interaction, "item registrar remove", err))
		}
		return newServiceError(interaction, "item registrar remove", err)
	}
	if !items.Contains(incomingItem) {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailRemoveItemNotExist(incomingItem.Name())); err != nil {
			log.Print(newServiceError(interaction, "item registrar remove", err))
		}
		return newServiceError(interaction, "item registrar remove", errors.New("item not exists"))
	}
	if err := itemRegistrar.itemStorage.Remove(ctx, incomingItem); err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailRemove(incomingItem.Name())); err != nil {
			log.Print(newServiceError(interaction, "item registrar remove", err))
		}
		return newServiceError(interaction, "item registrar remove", err)
	}
	choices, _ := command.NewChoices(len(items) - 1)
	appId, guildId := interaction.AppID, interaction.GuildID
	cmd := command.OfferBuilder(appId, guildId, choices).ApplicationCommand()
	if err := itemRegistrar.commandRegistrar.RegisterAppCommand(appId, guildId, cmd); err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemRegisterFail(incomingItem.Name())); err != nil {
			log.Print(newServiceError(interaction, "item registrar remove", err))
		}
		return newServiceError(interaction, "item registrar remove", err)
	}
	if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemRegisterSuccess(incomingItem.Name())); err != nil {
		log.Print(newServiceError(interaction, "item registrar add", err))
	}
	return nil
}

func (itemRegistrar *itemRegistrar) List(ctx context.Context, interaction *discordgo.Interaction) error {
	items, err := itemRegistrar.itemStorage.List(ctx)
	if err != nil {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailList("")); err != nil {
			log.Print(newServiceError(interaction, "item registrar list", err))
		}
		return newServiceError(interaction, "item registrar list", err)
	}
	if items.AreEmpty() {
		if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemFailItemsNotExist("")); err != nil {
			log.Print(newServiceError(interaction, "item registrar list", err))
		}
		return newServiceError(interaction, "item registrar list", errors.New("there is no items registered"))
	}
	if err := itemRegistrar.notifier.SendFollowUpMessage(interaction, followup.ItemSuccessList(items.String())); err != nil {
		log.Print(newServiceError(interaction, "item registrar list", err))
	}
	return nil
}
