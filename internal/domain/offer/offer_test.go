package offer

import (
	"errors"
	"testing"
)

// FIXME
// fix all tests -- added new type offerType / VendorIdentityIdentity
func TestOffersEquality(t *testing.T) {
	tests := []struct {
		name     string
		offer1   Offer
		offer2   Offer
		expected bool
	}{
		{
			"equal",
			Offer{product{"name1", 1.0}, 1},
			Offer{product{"name1", 1.0}, 1},
			true,
		},
		{
			"not equal",
			Offer{product{"name1", 1.1}, 1},
			Offer{product{"name1", 1.0}, 1},
			false,
		},
		{
			"equal",
			Offer{product{"name1", 1.0}, 2},
			Offer{product{"name1", 1.0}, 1},
			true,
		},
		{
			"not equal",
			Offer{product{"name1", 1.0}, 1},
			Offer{product{"name1", 1.1}, 1},
			false,
		},
		{
			"equal",
			Offer{product{"name1", 1.0}, 1},
			Offer{product{"name1", 1.0}, 2},
			true,
		},
		{
			"not equal",
			Offer{product{"name2", 1.0}, 1},
			Offer{product{"name1", 1.0}, 2},
			false,
		},
		{
			"not equal",
			Offer{product{"name1", 1.0}, 1},
			Offer{product{"name2", 1.0}, 2},
			false,
		},
		{
			"not equal",
			Offer{product{"name1", 1.0}, 1},
			Offer{product{"name2", 1.1}, 2},
			false,
		},
		{
			"not equal",
			Offer{product{"name2", 1.1}, 1},
			Offer{product{"name1", 1.0}, 2},
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// if actual := test.offer1.isSameOffer(test.offer2); actual != test.expected {
			// 	t.Fatalf("Expected: %v Actual: %v", actual, test.expected)
			// }
		})
	}
}

func TestOffersMerge(t *testing.T) {
	tests := []struct {
		name     string
		offer1   Offer
		offer2   Offer
		expected int
	}{
		{
			"merged",
			Offer{product{"name1", 1.0}, 1},
			Offer{product{"name1", 1.0}, 1},
			2,
		},
		{
			"merged",
			Offer{product{"name1", 1.1}, 1},
			Offer{product{"name1", 1.0}, 3},
			4,
		},
		{
			"merged",
			Offer{product{"name1", 1.0}, 2},
			Offer{product{"name1", 1.0}, 1},
			3,
		},
		{
			"merged",
			Offer{product{"name1", 1.0}, 1},
			Offer{product{"name1", 1.1}, 999},
			1000,
		},
		{
			"merged",
			Offer{product{"name1", 1.0}, 1},
			Offer{product{"name1", 1.0}, -1},
			0,
		},
		{
			"merged",
			Offer{product{"name1", 1.0}, 1},
			Offer{product{"name1", 1.0}, -2},
			-1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// merged := test.offer1.merge(test.offer2)
			// if actual := merged.count; actual != test.expected {
			// 	t.Fatalf("Expected %d Actual %d", test.expected, actual)
			// }
		})
	}
}

func TestNewOffer(t *testing.T) {
	t.Run("success creation", func(t *testing.T) {
		offer := NewOffer(product{"e", 1.0}, 2)
		err := offer.validate()
		if err != nil {
			t.Fatalf("creation failed %v", err)
		}
		if offer.count != 2 {
			t.Fatalf("expected count 2, actual %d", offer.count)
		}
		if offer.Product.name != "e" {
			t.Fatalf("expected name \"e\", actual %s", offer.Product.name)
		}
		if offer.Product.price != 1.0 {
			t.Fatalf("expected price 1.0, actual %f", offer.Product.price)
		}
	})
	t.Run("fail when offer count < 0", func(t *testing.T) {
		off := NewOffer(product{"e", 1.0}, -1)
		err := off.validate()
		if err == nil {
			t.Fatalf("should fail, but didnt")
		}
		if !errors.Is(err, ErrOfferCountLessOrEqualZero) {
			t.Fatalf("expected wrapped err: %v, actual: %v", ErrOfferCountLessOrEqualZero, err)
		}
	})
	t.Run("fail when offer count = 0", func(t *testing.T) {
		off := NewOffer(product{"e", 1.0}, 0)
		err := off.validate()
		if err == nil {
			t.Fatalf("should fail, but didnt")
		}
		if !errors.Is(err, ErrOfferCountLessOrEqualZero) {
			t.Fatalf("expected wrapped err: %v, actual: %v", ErrOfferCountLessOrEqualZero, err)
		}
	})
	t.Run("success when offer count > 0", func(t *testing.T) {
		off := NewOffer(product{"e", 1.0}, 2)
		err := off.validate()
		if err != nil {
			t.Fatalf("expected to create offer but error occured: err %v", err)
		}
	})
}

func TestNewVendorIdentity(t *testing.T) {
	t.Run("should fail when name is empty", func(t *testing.T) {
		VendorIdentity := NewVendorIdentity("aa")
		err := VendorIdentity.validate()
		if err == nil {
			t.Fatalf("expected to fail")
		}
		if !errors.Is(err, ErrVendorIdentityIsEmpty) {
			t.Fatalf("expected wrapped err: %v, actual: %v", ErrVendorIdentityIsEmpty, err)
		}
	})

	t.Run("success when name len = 1", func(t *testing.T) {
		VendorIdentity := NewVendorIdentity("a")
		err := VendorIdentity.validate()
		if err != nil {
			t.Fatalf("expected to create VendorIdentityIdentity. actual err: %v", err)
		}
		if VendorIdentity.id != "a" {
			t.Fatalf("VendorIdentityIdentity name is wrong. expected: %s actual %s", "a", VendorIdentity.id)
		}
	})
}
