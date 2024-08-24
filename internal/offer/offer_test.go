package offer

import (
	"errors"
	"testing"
)

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
			if actual := test.offer1.IsEqual(test.offer2); actual != test.expected {
				t.Fatalf("Expected: %v Actual: %v", actual, test.expected)
			}
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
			merged := test.offer1.Merge(test.offer2)
			if actual := merged.count; actual != test.expected {
				t.Fatalf("Expected %d Actual %d", test.expected, actual)
			}
		})
	}
}

func TestNewOffer(t *testing.T) {
	t.Run("success creation", func(t *testing.T) {
		offer, err := NewOffer(product{"e", 1.0}, 2)
		if err != nil {
			t.Fatalf("creation failed %v", err)
		}
		if offer.count != 2 {
			t.Fatalf("expected count 2, actual %d", offer.count)
		}
		if offer.product.name != "e" {
			t.Fatalf("expected name \"e\", actual %s", offer.product.name)
		}
		if offer.product.price != 1.0 {
			t.Fatalf("expected price 1.0, actual %f", offer.product.price)
		}
	})
	t.Run("fail when offer count < 0", func(t *testing.T) {
		_, err := NewOffer(product{"e", 1.0}, -1)
		if err == nil {
			t.Fatalf("should fail, but didnt")
		}
		if !errors.Is(err, ErrOfferCountLessOrEqualZero) {
			t.Fatalf("expected wrapped err: %v, actual: %v", ErrOfferCountLessOrEqualZero, err)
		}
	})
	t.Run("fail when offer count = 0", func(t *testing.T) {
		_, err := NewOffer(product{"e", 1.0}, 0)
		if err == nil {
			t.Fatalf("should fail, but didnt")
		}
		if !errors.Is(err, ErrOfferCountLessOrEqualZero) {
			t.Fatalf("expected wrapped err: %v, actual: %v", ErrOfferCountLessOrEqualZero, err)
		}
	})
	t.Run("success when offer count > 0", func(t *testing.T) {
		_, err := NewOffer(product{"e", 1.0}, 2)
		if err != nil {
			t.Fatalf("expected to create offer but error occured: err %v", err)
		}
	})
}

func TestNewVendor(t *testing.T) {
	t.Run("should fail when name is empty", func(t *testing.T) {
		_, err := NewVendor("")
		if err == nil {
			t.Fatalf("expected to fail")
		}
		if !errors.Is(err, ErrVendorIsEmpty) {
			t.Fatalf("expected wrapped err: %v, actual: %v", ErrVendorIsEmpty, err)
		}
	})

	t.Run("success when name len = 1", func(t *testing.T) {
		vendor, err := NewVendor("a")
		if err != nil {
			t.Fatalf("expected to create vendor. actual err: %v", err)
		}
		if vendor.name != "a" {
			t.Fatalf("vendor name is wrong. expected: %s actual %s", "a", vendor.name)
		}
	})
}
