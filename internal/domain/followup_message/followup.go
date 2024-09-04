package followup

import (
	"fmt"
	"math/rand"
)

type (
	offer struct {
		add        *offerAdd
		update     *offerUpdate
		remove     *offerRemove
		list       *offerList
		listVendor *offerListVendor
	}
	offerAdd struct {
		success         metadata
		fail            metadata
		updateOnAdd     metadata
		failUpdateOnAdd metadata
	}
	offerUpdate struct {
		success                  metadata
		fail                     metadata
		failOnOfferNotExist      metadata
		failOnNotHavingAnyOffers metadata
	}
	offerRemove struct {
		success                  metadata
		fail                     metadata
		failOnOfferNotExist      metadata
		failOnNotHavingAnyOffers metadata
	}
	offerList struct {
		success                  metadata
		fail                     metadata
		failOnNotHavingAnyOffers metadata
	}
	offerListVendor struct {
		success                  metadata
		fail                     metadata
		failOnNotHavingAnyOffers metadata
	}
	item struct {
		add    *itemAdd
		remove *itemRemove
		list   *itemListAll
	}
	itemAdd struct {
		success              metadata
		fail                 metadata
		failItemAlreadyExist metadata
	}
	itemRemove struct {
		success          metadata
		fail             metadata
		failItemNotExist metadata
	}
	itemListAll struct {
		success           metadata
		fail              metadata
		failItemsNotExist metadata
	}
)

type followUp struct {
	offer *offer
	item  *item
}

var followUpMessage = &followUp{
	offer: &offer{
		add: &offerAdd{
			success:         metadata{"offer_success_add"},
			fail:            metadata{"offer_fail_add"},
			updateOnAdd:     metadata{"offer_success_update_on_add"},
			failUpdateOnAdd: metadata{"offer_fail_update_on_add"},
		},
		update: &offerUpdate{
			success:                  metadata{"offer_success_update"},
			fail:                     metadata{"offer_fail_update"},
			failOnOfferNotExist:      metadata{"offer_fail_update_offer_not_exist"},
			failOnNotHavingAnyOffers: metadata{"offer_fail_update_on_not_having_any_offers"},
		},
		remove: &offerRemove{
			success:                  metadata{"offer_success_remove"},
			fail:                     metadata{"offer_fail_remove"},
			failOnOfferNotExist:      metadata{"offer_fail_remove_offer_not_exist"},
			failOnNotHavingAnyOffers: metadata{"offer_fail_remove_on_not_having_any_offers"},
		},
		list: &offerList{
			success:                  metadata{"offer_success_list"},
			fail:                     metadata{"offer_fail_list"},
			failOnNotHavingAnyOffers: metadata{"offer_fail_list_on_not_having_any_offers"},
		},
		listVendor: &offerListVendor{
			success:                  metadata{"offer_success_list_vendor"},
			fail:                     metadata{"offer_fail_list_vendor"},
			failOnNotHavingAnyOffers: metadata{"offer_fail_list_vendor_on_not_having_any_offers"},
		},
	},
	item: &item{
		add: &itemAdd{
			success:              metadata{"item_success_add"},
			fail:                 metadata{"item_fail_add"},
			failItemAlreadyExist: metadata{"item_fail_add_item_already_exist"},
		},
		remove: &itemRemove{
			success:          metadata{"item_success_remove"},
			fail:             metadata{"item_fail_remove"},
			failItemNotExist: metadata{"item_fail_remove_item_not_exist"},
		},
		list: &itemListAll{
			success:           metadata{"item_success_list"},
			fail:              metadata{"item_fail_list"},
			failItemsNotExist: metadata{"item_fail_items_not_exist"},
		},
	},
}

var (
	OfferSuccessAdd = func(name string) Message {
		return Message{followUpMessage.offer.add.success, fmt.Sprintf("Offer with item %s successfully added! Need more help? Just ask!", name), name}
	}
	OfferSuccessUpdateOnAdd = func(name string) Message {
		return Message{followUpMessage.offer.add.updateOnAdd, fmt.Sprintf("Offer with item %s has been updated because you already have same offer. Need more help? Just ask!", name), name}
	}
	OfferFailAdd = func(name string) Message {
		return Message{followUpMessage.offer.add.fail, fmt.Sprintf("Failed to add offer with item %s. Please try again or ask for help.", name), name}
	}
	OfferFailUpdateOnAdd = func(name string) Message {
		return Message{followUpMessage.offer.add.failUpdateOnAdd, fmt.Sprintf("Wanted to update offer %s because you already have it, but unsuccessfully. Please try again.", name), name}
	}
	OfferSuccessRemove = func(name string) Message {
		return Message{followUpMessage.offer.remove.success, fmt.Sprintf("Offer with item %s successfully removed! Need more help? Just ask!", name), name}
	}
	OfferFailRemove = func(name string) Message {
		return Message{followUpMessage.offer.remove.fail, fmt.Sprintf("Failed to remove offer with item %s. Please try again or ask for help.", name), name}
	}
	OfferFailRemoveOnNotHavingAnyOffers = func(name string) Message {
		return Message{followUpMessage.offer.remove.failOnNotHavingAnyOffers, fmt.Sprintf("Can't remove offer with item %s because you dont have any offers. Need more help? Just ask!", name), name}
	}
	OfferFailRemoveOfferNotExists = func(name string) Message {
		return Message{followUpMessage.offer.remove.failOnOfferNotExist, fmt.Sprintf("Can't remove offer with item %s because you dont have offer with this item. Need more help? Just ask!", name), name}
	}
	OfferSuccessUpdate = func(name string) Message {
		return Message{followUpMessage.offer.update.success, fmt.Sprintf("Item %s successfully updated! Need more help? Just ask!", name), name}
	}
	OfferFailUpdate = func(name string) Message {
		return Message{followUpMessage.offer.update.fail, fmt.Sprintf("Failed to update offer with item %s. Please try again or ask for help.", name), name}
	}
	OfferSuccessList = func(items string) Message {
		return Message{followUpMessage.offer.list.success, fmt.Sprintf("List of retrieved offers:\n%s\n Feel free to ask for more!", items), items}
	}
	OfferFailList = func(name string) Message {
		return Message{followUpMessage.offer.list.fail, fmt.Sprintf("Failed to retrieve offers with name %s. Please try again or ask for help.", name), name}
	}
	OfferSuccessListVendor = func(items string) Message {
		return Message{followUpMessage.offer.listVendor.success, fmt.Sprintf("Your offers has been successfully retrieved!\n%s\nNeed more help? Just ask!", items), items}
	}
	OfferFailListVendor = func(name string) Message {
		return Message{followUpMessage.offer.listVendor.fail, "Failed to retrieve your offers. Please try again or ask for help.", name}
	}
	OfferFailUpdateOnNotHavingAnyOffers = func(name string) Message {
		return Message{followUpMessage.offer.update.failOnNotHavingAnyOffers, fmt.Sprintf("Can't update offer with item %s because you dont have any offers. Please try again or ask for help.", name), name}
	}
	OfferFailUpdateOfferNotExists = func(name string) Message {
		return Message{followUpMessage.offer.update.failOnOfferNotExist, fmt.Sprintf("Can't update offer with item %s because you dont have any offers. Please try again or ask for help.", name), name}
	}
	OfferFailListVendorNotHavingAnyOffers = func(name string) Message {
		return Message{followUpMessage.offer.listVendor.failOnNotHavingAnyOffers, "Can't list offers because you do not have any. Need more help? Just ask!", name}
	}
	OfferFailListNotHavingAnyOffers = func(name string) Message {
		return Message{followUpMessage.offer.list.failOnNotHavingAnyOffers, fmt.Sprintf("We do not have offers with name %s. Need more help? Just ask!", name), name}
	}
	ItemSuccessAdd = func(name string) Message {
		return Message{followUpMessage.item.add.success, fmt.Sprintf("Item %s successfully added! Need more help? Just ask!", name), name}
	}
	ItemFailAdd = func(name string) Message {
		return Message{followUpMessage.item.add.fail, fmt.Sprintf("Failed to add item %s. Please try again or ask for help.", name), name}
	}
	ItemFailAddItemAlreadyExist = func(name string) Message {
		return Message{followUpMessage.item.add.failItemAlreadyExist, fmt.Sprintf("Item %s already exists and cannot be added again. Need help with something else?", name), name}
	}
	ItemSuccessRemove = func(name string) Message {
		return Message{followUpMessage.item.remove.success, fmt.Sprintf("Item %s successfully removed! Need more help? Just ask!", name), name}
	}
	ItemFailRemove = func(name string) Message {
		return Message{followUpMessage.item.remove.fail, fmt.Sprintf("Failed to remove item %s. Please try again or ask for help.", name), name}
	}
	ItemFailRemoveItemNotExist = func(name string) Message {
		return Message{followUpMessage.item.remove.failItemNotExist, fmt.Sprintf("Item %s does not exist and cannot be removed. Need help with something else?", name), name}
	}
	ItemSuccessList = func(items string) Message {
		return Message{followUpMessage.item.list.success, fmt.Sprintf("List of retrieved items:\n%s\nFeel free to ask for more!", items), items}
	}
	ItemFailList = func(name string) Message {
		return Message{followUpMessage.item.list.fail, "Failed to retrieve items, something wrong happened. Please try again or ask for another help.", name}
	}
	ItemFailItemsNotExist = func(name string) Message {
		return Message{followUpMessage.item.list.failItemsNotExist, "No items registered actually. Need help with something else?", name}
	}
)

var messageBucket = map[metadata][]messageGenerator{
	followUpMessage.offer.add.success: {
		messageGenerator(OfferSuccessAdd),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.success, fmt.Sprintf("Item %s gracefully added! Need help? Speeding to you up!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.success, fmt.Sprintf("Item %s successfully integrated! Anything else you need?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.success, fmt.Sprintf("Successfully added item %s! Need further assistance? Let us know!", name), name}
		}),
	},
	followUpMessage.offer.add.fail: {
		messageGenerator(OfferFailAdd),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.fail, fmt.Sprintf("Oops, adding item %s failed. Want to try again or need help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.fail, fmt.Sprintf("Adding %s didn't work out. Let's give it another shot or get some help!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.fail, fmt.Sprintf("Failed to add %s. Don't worry, we're here to help!", name), name}
		}),
	},
	followUpMessage.offer.add.updateOnAdd: {
		messageGenerator(OfferSuccessUpdateOnAdd),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.updateOnAdd, fmt.Sprintf("Item %s updated as you already had this offer. Need anything else?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.updateOnAdd, fmt.Sprintf("Your offer for item %s was updated. Any further assistance needed?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.updateOnAdd, fmt.Sprintf("Updated the existing offer with item %s. How else can we help?", name), name}
		}),
	},
	followUpMessage.offer.update.success: {
		messageGenerator(OfferSuccessUpdate),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.success, fmt.Sprintf("Item %s successfully updated! What else can we do for you?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.success, fmt.Sprintf("Update successful for item %s! Need further assistance?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.success, fmt.Sprintf("The update for item %s was a success! How can we assist you next?", name), name}
		}),
	},
	followUpMessage.offer.update.fail: {
		messageGenerator(OfferFailUpdate),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.fail, fmt.Sprintf("Failed to update item %s. Shall we try again or do you need help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.fail, fmt.Sprintf("Updating item %s didn't work. How can we assist?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.fail, fmt.Sprintf("We couldn't update item %s. Let us know how to assist you!", name), name}
		}),
	},
	followUpMessage.offer.update.failOnOfferNotExist: {
		messageGenerator(OfferFailUpdateOfferNotExists),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.failOnOfferNotExist, fmt.Sprintf("Can't update item %s because this offer doesn't exist. Need help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.failOnOfferNotExist, fmt.Sprintf("Offer for item %s not found, hence couldn't update. Assistance required?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.failOnOfferNotExist, fmt.Sprintf("No existing offer found for item %s to update. Need more help?", name), name}
		}),
	},
	followUpMessage.offer.list.success: {
		messageGenerator(OfferSuccessList),
		messageGenerator(func(items string) Message {
			return Message{followUpMessage.offer.list.success, fmt.Sprintf("Here are your retrieved offers:\n%s\nNeed more details?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{followUpMessage.offer.list.success, fmt.Sprintf("Offers listed successfully:\n%s\nAny other help needed?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{followUpMessage.offer.list.success, fmt.Sprintf("Successfully retrieved the following offers:\n%s\nWhat else can we do?", items), items}
		}),
	},
	followUpMessage.offer.listVendor.success: {
		messageGenerator(OfferSuccessListVendor),
		messageGenerator(func(items string) Message {
			return Message{followUpMessage.offer.listVendor.success, fmt.Sprintf("Successfully retrieved your offers:\n%s\nNeed more help?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{followUpMessage.offer.listVendor.success, fmt.Sprintf("Here are your offers:\n%s\nAnything else we can do?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{followUpMessage.offer.listVendor.success, fmt.Sprintf("Your offers are listed:\n%s\nFeel free to ask for more details!", items), items}
		}),
	},
	followUpMessage.offer.list.fail: {
		messageGenerator(OfferFailList),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.list.fail, fmt.Sprintf("Failed to retrieve offers with name %s. Need more assistance?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.list.fail, fmt.Sprintf("Unable to get offers for %s. How can we help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.list.fail, fmt.Sprintf("Retrieval of offers for %s was unsuccessful. Need any help?", name), name}
		}),
	},
	followUpMessage.offer.listVendor.fail: {
		messageGenerator(OfferFailListVendor),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.listVendor.fail, "Failed to retrieve your offers. Need more assistance or want to retry?", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.listVendor.fail, "Couldn’t list your offers. Let us know if you need help!", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.listVendor.fail, "Retrieving your offers didn't work. How can we assist you?", name}
		}),
	},
	followUpMessage.offer.list.failOnNotHavingAnyOffers: {
		messageGenerator(OfferFailListNotHavingAnyOffers),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.list.failOnNotHavingAnyOffers, fmt.Sprintf("We do not have offers with name %s. Need more help? Just ask!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.list.failOnNotHavingAnyOffers, fmt.Sprintf("No offers found with the name %s. How can we assist?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.list.failOnNotHavingAnyOffers, fmt.Sprintf("Can't find offers for %s. Feel free to ask for more details!", name), name}
		}),
	},
	followUpMessage.offer.listVendor.failOnNotHavingAnyOffers: {
		messageGenerator(OfferFailListVendorNotHavingAnyOffers),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.listVendor.failOnNotHavingAnyOffers, "Can't list offers because you do not have any. Need more help? Just ask!", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.listVendor.failOnNotHavingAnyOffers, "No offers available to list. Let us know if you need assistance!", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.listVendor.failOnNotHavingAnyOffers, "There are no offers to display. How else can we help you?", name}
		}),
	},
	followUpMessage.offer.remove.fail: {
		messageGenerator(OfferFailRemove),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.fail, fmt.Sprintf("Failed to remove offer with item %s. Please try again or ask for help.", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.fail, fmt.Sprintf("Oops, removing item %s didn't work. How can we assist?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.fail, fmt.Sprintf("Unable to remove item %s. Need further help?", name), name}
		}),
	},
	followUpMessage.offer.remove.success: {
		messageGenerator(OfferSuccessRemove),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.success, fmt.Sprintf("Offer with item %s successfully removed! Need more help? Just ask!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.success, fmt.Sprintf("Item %s removed successfully! What else can we do for you?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.success, fmt.Sprintf("Successfully removed item %s! Feel free to ask for more assistance.", name), name}
		}),
	},
	followUpMessage.offer.remove.failOnOfferNotExist: {
		messageGenerator(OfferFailRemoveOfferNotExists),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.failOnOfferNotExist, fmt.Sprintf("Can't remove offer with item %s because you don't have an offer with this item. Need more help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.failOnOfferNotExist, fmt.Sprintf("Offer %s not found for removal. Need assistance?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.failOnOfferNotExist, fmt.Sprintf("The offer with item %s does not exist. How can we assist you?", name), name}
		}),
	},
	followUpMessage.offer.remove.failOnNotHavingAnyOffers: {
		messageGenerator(OfferFailRemoveOnNotHavingAnyOffers),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.failOnNotHavingAnyOffers, fmt.Sprintf("Can't remove offer with item %s because you don't have any offers. Need more help? Just ask!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.failOnNotHavingAnyOffers, "No offers available to remove. How can we assist?", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.remove.failOnNotHavingAnyOffers, fmt.Sprintf("Removing %s is not possible as you have no offers. Need more help?", name), name}
		}),
	},
	followUpMessage.offer.update.failOnNotHavingAnyOffers: {
		messageGenerator(OfferFailUpdateOnNotHavingAnyOffers),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.failOnNotHavingAnyOffers, fmt.Sprintf("Can't update offer with item %s because you don't have any offers. Please try again or ask for help.", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.failOnNotHavingAnyOffers, fmt.Sprintf("Updating %s isn't possible due to lack of offers. How can we help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.update.failOnNotHavingAnyOffers, fmt.Sprintf("Update failed for %s as you have no offers. Need assistance?", name), name}
		}),
	},
	followUpMessage.offer.add.failUpdateOnAdd: {
		messageGenerator(OfferFailUpdateOnAdd),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.failUpdateOnAdd, fmt.Sprintf("We are sorry, we wanted to update offer with item %s instead of adding it because you already have one, but we failed. Please try again or ask for another help.", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.failUpdateOnAdd, fmt.Sprintf("We didnt added nor updated offer with item %s, something wrong happened. Can we help in another way?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.offer.add.failUpdateOnAdd, fmt.Sprintf("Adding offer with item %s failed, precisely speaking update did, because you already have same offer. We are sorry. Do you need another help?", name), name}

		}),
	},
	followUpMessage.item.add.success: {
		messageGenerator(ItemSuccessAdd),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.add.success, fmt.Sprintf("Successfully added item %s! How can we assist you further?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.add.success, fmt.Sprintf("Item %s added successfully! Need any more help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.add.success, fmt.Sprintf("Item %s has been added successfully! Anything else you need assistance with?", name), name}
		}),
	},
	followUpMessage.item.add.fail: {
		messageGenerator(ItemFailAdd),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.add.fail, fmt.Sprintf("Failed to add item %s. Would you like to try again or need assistance?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.add.fail, fmt.Sprintf("Adding item %s didn't work, something bad occured. How can we assist?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.add.fail, fmt.Sprintf("Unable to add item %s. Need help with anything else?", name), name}
		}),
	},
	followUpMessage.item.add.failItemAlreadyExist: {
		messageGenerator(ItemFailAddItemAlreadyExist),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.add.failItemAlreadyExist, fmt.Sprintf("Item %s already exists, so it wont be added. How else can we assist you?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.add.failItemAlreadyExist, fmt.Sprintf("The item %s is already registered, cant be added then. Need help with something else?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.add.failItemAlreadyExist, fmt.Sprintf("Item %s already exists, so it can't be added again. What else can we do for you?", name), name}
		}),
	},
	followUpMessage.item.remove.success: {
		messageGenerator(ItemSuccessRemove),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.remove.success, fmt.Sprintf("Item %s successfully removed! How else can we help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.remove.success, fmt.Sprintf("Removed item %s successfully! Any other assistance needed?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.remove.success, fmt.Sprintf("Successfully removed item %s! Let us know if you need more help.", name), name}
		}),
	},
	followUpMessage.item.remove.fail: {
		messageGenerator(ItemFailRemove),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.remove.fail, fmt.Sprintf("Failed to remove item %s, something bad occured. Would you like to try again or want other help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.remove.fail, fmt.Sprintf("Could not remove item %s, something wrong happened. How can we assist?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.remove.fail, fmt.Sprintf("Item %s could not be removed. Need more help?", name), name}
		}),
	},
	followUpMessage.item.remove.failItemNotExist: {
		messageGenerator(ItemFailRemoveItemNotExist),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.remove.failItemNotExist, fmt.Sprintf("Item %s does not exist, so it can't be removed. Need assistance with anything else?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.remove.failItemNotExist, fmt.Sprintf("Cannot remove item %s because it does not exist. How else can we assist?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.remove.failItemNotExist, fmt.Sprintf("Item %s is not in the list, so it can't be removed. Need more help?", name), name}
		}),
	},
	followUpMessage.item.list.success: {
		messageGenerator(ItemSuccessList),
		messageGenerator(func(items string) Message {
			return Message{followUpMessage.item.list.success, fmt.Sprintf("List of items retrieved:\n%s\nNeed more help?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{followUpMessage.item.list.success, fmt.Sprintf("Successfully retrieved the following items:\n%s\nAnything else?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{followUpMessage.item.list.success, fmt.Sprintf("Items listed successfully:\n%s\nWhat else can we do for you?", items), items}
		}),
	},
	followUpMessage.item.list.fail: {
		messageGenerator(ItemFailList),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.list.fail, "Failed to retrieve items, something bad occured. Need more help? Feel free to ask!", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.list.fail, "Unable to retrieve items, something wrong happened. Can we assist you in something else?", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.list.fail, "Retrieving items was unsuccessful. Need any help?", name}
		}),
	},
	followUpMessage.item.list.failItemsNotExist: {
		messageGenerator(ItemFailItemsNotExist),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.list.failItemsNotExist, "No items found. If you want register type command with add", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.list.failItemsNotExist, "You dont have any items registered now, feel free to add!", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{followUpMessage.item.list.failItemsNotExist, "Can't find any items, because you didnt registered any. Need more help? Feel free to ask!", name}
		}),
	},
}

type messageGenerator func(string) Message

type metadata struct {
	msgType string
}

type Message struct {
	metadata metadata
	content  string
	arg      string
}

func (m Message) Content() string {
	return m.content
}

func (m Message) Randomize() Message {
	if msgs, ok := messageBucket[m.metadata]; ok {
		rand := rand.Intn(len(msgs) - 1)
		return msgs[rand](m.arg)
	}
	return m
}
