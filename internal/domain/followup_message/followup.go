package followup

import (
	"fmt"
	"math/rand"
)

var (
	successAddType      = add{metadata{"success_add"}}
	failAddType         = add{metadata{"fail_add"}}
	updateOnAddType     = add{metadata{"success_update_on_add"}}
	failUpdateOnAddType = add{metadata{"fail_update_on_add"}}

	successUpdateType                  = update{metadata{"success_update"}}
	failUpdateType                     = update{metadata{"fail_update"}}
	failUpdateOfferNotExistType        = update{metadata{"fail_update_offer_not_exist"}}
	failUpdateOnNotHavingAnyOffersType = update{metadata{"fail_update_on_not_having_any_offers"}}

	successRemoveType              = remove{metadata{"success_remove"}}
	failRemoveType                 = remove{metadata{"fail_remove"}}
	failRemoveOnOfferNotExistType  = remove{metadata{"fail_remove_offer_not_exist"}}
	failRemoveOnNotHavingAnyOffers = remove{metadata{"fail_remove_on_not_having_any_offers"}}

	successListType                  = list{metadata{"success_list"}}
	failListType                     = list{metadata{"fail_list"}}
	failListOnNotHavingAnyOffersType = list{metadata{"fail_list_on_not_having_any_offers"}}

	failListVendorType                     = listVendor{metadata{"fail_list_vendor"}}
	failListVendorOnNotHavingAnyOffersType = listVendor{metadata{"fail_list_vendor_on_not_having_any_offers"}}
	successListVendorType                  = listVendor{metadata{"success_list_vendor"}}
)

var (
	SuccessAdd = func(name string) Message {
		return Message{successAddType.metadata, fmt.Sprintf("Offer with item %s successfully added! Need more help? Just ask!", name), name}
	}
	SuccessUpdateOnAdd = func(name string) Message {
		return Message{updateOnAddType.metadata, fmt.Sprintf("Offer with item %s has been updated because you already have same offer. Need more help? Just ask!", name), name}
	}
	FailAdd = func(name string) Message {
		return Message{failAddType.metadata, fmt.Sprintf("Failed to add offer with item %s. Please try again or ask for help.", name), name}
	}
	FailUpdateOnAdd = func(name string) Message {
		return Message{failUpdateOnAddType.metadata, fmt.Sprintf("Wanted to update offer %s because you already have it, but unsuccessfully. Please try again.", name), name}
	}
	SuccessRemove = func(name string) Message {
		return Message{successRemoveType.metadata, fmt.Sprintf("Offer with item %s successfully removed! Need more help? Just ask!", name), name}
	}
	FailRemove = func(name string) Message {
		return Message{failRemoveType.metadata, fmt.Sprintf("Failed to remove offer with item %s. Please try again or ask for help.", name), name}
	}
	FailRemoveOnNotHavingAnyOffers = func(name string) Message {
		return Message{failRemoveOnNotHavingAnyOffers.metadata, fmt.Sprintf("Can't remove offer with item %s because you dont have any offers. Need more help? Just ask!", name), name}
	}
	FailRemoveOfferNotExists = func(name string) Message {
		return Message{failRemoveOnOfferNotExistType.metadata, fmt.Sprintf("Can't remove offer with item %s because you dont have offer with this item. Need more help? Just ask!", name), name}
	}
	SuccessUpdate = func(name string) Message {
		return Message{successUpdateType.metadata, fmt.Sprintf("Item %s successfully updated! Need more help? Just ask!", name), name}
	}
	FailUpdate = func(name string) Message {
		return Message{failUpdateType.metadata, fmt.Sprintf("Failed to update offer with item %s. Please try again or ask for help.", name), name}
	}
	SuccessList = func(items string) Message {
		return Message{successListType.metadata, fmt.Sprintf("List of retrieved offers:\n%s\n Feel free to ask for more!", items), items}
	}
	FailList = func(name string) Message {
		return Message{failListType.metadata, fmt.Sprintf("Failed to retrieve offers with name %s. Please try again or ask for help.", name), name}
	}
	SuccessListVendor = func(items string) Message {
		return Message{successListVendorType.metadata, fmt.Sprintf("Your offers has been successfully retrieved!\n%s\nNeed more help? Just ask!", items), items}
	}
	FailListVendor = func(name string) Message {
		return Message{failListVendorType.metadata, "Failed to retrieve your offers. Please try again or ask for help.", name}
	}
	FailUpdateOnNotHavingAnyOffers = func(name string) Message {
		return Message{failUpdateOnNotHavingAnyOffersType.metadata, fmt.Sprintf("Can't update offer with item %s because you dont have any offers. Please try again or ask for help.", name), name}
	}
	FailUpdateOfferNotExists = func(name string) Message {
		return Message{failUpdateOfferNotExistType.metadata, fmt.Sprintf("Can't update offer with item %s because you dont have any offers. Please try again or ask for help.", name), name}
	}
	FailListVendorNotHavingAnyOffers = func(name string) Message {
		return Message{failListVendorOnNotHavingAnyOffersType.metadata, "Can't list offers because you do not have any. Need more help? Just ask!", name}
	}
	FailListNotHavingAnyOffers = func(name string) Message {
		return Message{failListOnNotHavingAnyOffersType.metadata, fmt.Sprintf("We do not have offers with name %s. Need more help? Just ask!", name), name}
	}
)

var messageBucket = map[metadata][]messageGenerator{
	successAddType.metadata: {
		messageGenerator(SuccessAdd),
		messageGenerator(func(name string) Message {
			return Message{successAddType.metadata, fmt.Sprintf("Item %s gracefully added! Need help? Speeding to you up!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{successAddType.metadata, fmt.Sprintf("Item %s successfully integrated! Anything else you need?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{successAddType.metadata, fmt.Sprintf("Successfully added item %s! Need further assistance? Let us know!", name), name}
		}),
	},
	failAddType.metadata: {
		messageGenerator(FailAdd),
		messageGenerator(func(name string) Message {
			return Message{failAddType.metadata, fmt.Sprintf("Oops, adding item %s failed. Want to try again or need help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failAddType.metadata, fmt.Sprintf("Adding %s didn't work out. Let's give it another shot or get some help!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failAddType.metadata, fmt.Sprintf("Failed to add %s. Don't worry, we're here to help!", name), name}
		}),
	},
	updateOnAddType.metadata: {
		messageGenerator(SuccessUpdateOnAdd),
		messageGenerator(func(name string) Message {
			return Message{updateOnAddType.metadata, fmt.Sprintf("Item %s updated as you already had this offer. Need anything else?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{updateOnAddType.metadata, fmt.Sprintf("Your offer for item %s was updated. Any further assistance needed?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{updateOnAddType.metadata, fmt.Sprintf("Updated the existing offer with item %s. How else can we help?", name), name}
		}),
	},
	successUpdateType.metadata: {
		messageGenerator(SuccessUpdate),
		messageGenerator(func(name string) Message {
			return Message{successUpdateType.metadata, fmt.Sprintf("Item %s successfully updated! What else can we do for you?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{successUpdateType.metadata, fmt.Sprintf("Update successful for item %s! Need further assistance?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{successUpdateType.metadata, fmt.Sprintf("The update for item %s was a success! How can we assist you next?", name), name}
		}),
	},
	failUpdateType.metadata: {
		messageGenerator(FailUpdate),
		messageGenerator(func(name string) Message {
			return Message{failUpdateType.metadata, fmt.Sprintf("Failed to update item %s. Shall we try again or do you need help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failUpdateType.metadata, fmt.Sprintf("Updating item %s didn't work. How can we assist?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failUpdateType.metadata, fmt.Sprintf("We couldn't update item %s. Let us know how to assist you!", name), name}
		}),
	},
	failUpdateOfferNotExistType.metadata: {
		messageGenerator(FailUpdateOfferNotExists),
		messageGenerator(func(name string) Message {
			return Message{failUpdateOfferNotExistType.metadata, fmt.Sprintf("Can't update item %s because this offer doesn't exist. Need help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failUpdateOfferNotExistType.metadata, fmt.Sprintf("Offer for item %s not found, hence couldn't update. Assistance required?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failUpdateOfferNotExistType.metadata, fmt.Sprintf("No existing offer found for item %s to update. Need more help?", name), name}
		}),
	},
	successListType.metadata: {
		messageGenerator(SuccessList),
		messageGenerator(func(items string) Message {
			return Message{successListType.metadata, fmt.Sprintf("Here are your retrieved offers:\n%s\nNeed more details?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{successListType.metadata, fmt.Sprintf("Offers listed successfully:\n%s\nAny other help needed?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{successListType.metadata, fmt.Sprintf("Successfully retrieved the following offers:\n%s\nWhat else can we do?", items), items}
		}),
	},
	successListVendorType.metadata: {
		messageGenerator(SuccessListVendor),
		messageGenerator(func(items string) Message {
			return Message{successListVendorType.metadata, fmt.Sprintf("Successfully retrieved your offers:\n%s\nNeed more help?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{successListVendorType.metadata, fmt.Sprintf("Here are your offers:\n%s\nAnything else we can do?", items), items}
		}),
		messageGenerator(func(items string) Message {
			return Message{successListVendorType.metadata, fmt.Sprintf("Your offers are listed:\n%s\nFeel free to ask for more details!", items), items}
		}),
	},
	failListType.metadata: {
		messageGenerator(FailList),
		messageGenerator(func(name string) Message {
			return Message{failListType.metadata, fmt.Sprintf("Failed to retrieve offers with name %s. Need more assistance?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failListType.metadata, fmt.Sprintf("Unable to get offers for %s. How can we help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failListType.metadata, fmt.Sprintf("Retrieval of offers for %s was unsuccessful. Need any help?", name), name}
		}),
	},
	failListVendorType.metadata: {
		messageGenerator(FailListVendor),
		messageGenerator(func(name string) Message {
			return Message{failListVendorType.metadata, "Failed to retrieve your offers. Need more assistance or want to retry?", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failListVendorType.metadata, "Couldn’t list your offers. Let us know if you need help!", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failListVendorType.metadata, "Retrieving your offers didn't work. How can we assist you?", name}
		}),
	},
	failListOnNotHavingAnyOffersType.metadata: {
		messageGenerator(FailListNotHavingAnyOffers),
		messageGenerator(func(name string) Message {
			return Message{failListOnNotHavingAnyOffersType.metadata, fmt.Sprintf("We do not have offers with name %s. Need more help? Just ask!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failListOnNotHavingAnyOffersType.metadata, fmt.Sprintf("No offers found with the name %s. How can we assist?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failListOnNotHavingAnyOffersType.metadata, fmt.Sprintf("Can't find offers for %s. Feel free to ask for more details!", name), name}
		}),
	},
	failListVendorOnNotHavingAnyOffersType.metadata: {
		messageGenerator(FailListVendorNotHavingAnyOffers),
		messageGenerator(func(name string) Message {
			return Message{failListVendorOnNotHavingAnyOffersType.metadata, "Can't list offers because you do not have any. Need more help? Just ask!", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failListVendorOnNotHavingAnyOffersType.metadata, "No offers available to list. Let us know if you need assistance!", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failListVendorOnNotHavingAnyOffersType.metadata, "There are no offers to display. How else can we help you?", name}
		}),
	},
	failRemoveType.metadata: {
		messageGenerator(FailRemove),
		messageGenerator(func(name string) Message {
			return Message{failRemoveType.metadata, fmt.Sprintf("Failed to remove offer with item %s. Please try again or ask for help.", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failRemoveType.metadata, fmt.Sprintf("Oops, removing item %s didn't work. How can we assist?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failRemoveType.metadata, fmt.Sprintf("Unable to remove item %s. Need further help?", name), name}
		}),
	},
	successRemoveType.metadata: {
		messageGenerator(SuccessRemove),
		messageGenerator(func(name string) Message {
			return Message{successRemoveType.metadata, fmt.Sprintf("Offer with item %s successfully removed! Need more help? Just ask!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{successRemoveType.metadata, fmt.Sprintf("Item %s removed successfully! What else can we do for you?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{successRemoveType.metadata, fmt.Sprintf("Successfully removed item %s! Feel free to ask for more assistance.", name), name}
		}),
	},
	failRemoveOnOfferNotExistType.metadata: {
		messageGenerator(FailRemoveOfferNotExists),
		messageGenerator(func(name string) Message {
			return Message{failRemoveOnOfferNotExistType.metadata, fmt.Sprintf("Can't remove offer with item %s because you don't have an offer with this item. Need more help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failRemoveOnOfferNotExistType.metadata, fmt.Sprintf("Offer %s not found for removal. Need assistance?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failRemoveOnOfferNotExistType.metadata, fmt.Sprintf("The offer with item %s does not exist. How can we assist you?", name), name}
		}),
	},
	failRemoveOnNotHavingAnyOffers.metadata: {
		messageGenerator(FailRemoveOnNotHavingAnyOffers),
		messageGenerator(func(name string) Message {
			return Message{failRemoveOnNotHavingAnyOffers.metadata, fmt.Sprintf("Can't remove offer with item %s because you don't have any offers. Need more help? Just ask!", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failRemoveOnNotHavingAnyOffers.metadata, "No offers available to remove. How can we assist?", name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failRemoveOnNotHavingAnyOffers.metadata, fmt.Sprintf("Removing %s is not possible as you have no offers. Need more help?", name), name}
		}),
	},
	failUpdateOnNotHavingAnyOffersType.metadata: {
		messageGenerator(FailUpdateOnNotHavingAnyOffers),
		messageGenerator(func(name string) Message {
			return Message{failUpdateOnNotHavingAnyOffersType.metadata, fmt.Sprintf("Can't update offer with item %s because you don't have any offers. Please try again or ask for help.", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failUpdateOnNotHavingAnyOffersType.metadata, fmt.Sprintf("Updating %s isn't possible due to lack of offers. How can we help?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failUpdateOnNotHavingAnyOffersType.metadata, fmt.Sprintf("Update failed for %s as you have no offers. Need assistance?", name), name}
		}),
	},
	failUpdateOnAddType.metadata: {
		messageGenerator(FailUpdateOnAdd),
		messageGenerator(func(name string) Message {
			return Message{failUpdateOnAddType.metadata, fmt.Sprintf("We are sorry, we wanted to update offer with item %s instead of adding it because you already have one, but we failed. Please try again or ask for another help.", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failUpdateOnAddType.metadata, fmt.Sprintf("We didnt added nor updated offer with item %s, something wrong happened. Can we help in another way?", name), name}
		}),
		messageGenerator(func(name string) Message {
			return Message{failUpdateOnAddType.metadata, fmt.Sprintf("Adding offer with item %s failed, precisely speaking update did, because you already have same offer. We are sorry. Do you need another help?", name), name}
		}),
	},
}

type messageGenerator func(string) Message

type metadata struct {
	msgType string
}

type (
	add        struct{ metadata }
	remove     struct{ metadata }
	update     struct{ metadata }
	list       struct{ metadata }
	listVendor struct{ metadata }
)

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
