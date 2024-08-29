package offer

import (
	"errors"
	"testing"
)

// FIXME fix all tests
func TestProductValidation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := NewProduct("a", 1.0)
		err := p.validate()
		if err != nil {
			t.Fatalf("error occured but sohuldnt, %v", err)
		}
		if p.name != "a" {
			t.Fatalf("product name not match, expected: %s actual %s", "a", p.name)
		}
		if p.price != 1.0 {
			t.Fatalf("product price not match, expected: %f, actual: %f", 1.0, p.price)
		}
	})

	t.Run("should fail when name is empty", func(t *testing.T) {
		prod := NewProduct("", 1.0)
		err := prod.validate()
		if err == nil {
			t.Fatalf("error should occur")
		}
		if !errors.Is(err, ErrProductNameEmpty) {
			t.Fatalf("expected err: %v, actual: %v", ErrProductNameEmpty, err)
		}
	})
	t.Run("should fail when price < 0", func(t *testing.T) {
		prod := NewProduct("a", -1.0)
		err := prod.validate()
		if err == nil {
			t.Fatalf("error should occur")
		}
		if !errors.Is(err, ErrProductPriceLessThanZero) {
			t.Fatalf("expected err: %v, actual: %v", ErrProductPriceLessThanZero, err)
		}
	})

	t.Run("success when price = 0", func(t *testing.T) {
		prod := NewProduct("a", 0.0)
		err := prod.validate()
		if err != nil {
			t.Fatalf("error occured but shouldnot. err: %v", err)
		}
	})
	t.Run("success when name len = 1", func(t *testing.T) {
		prod := NewProduct("a", 1)
		err := prod.validate()
		if err != nil {
			t.Fatalf("error occured, but shouldnt, %v", err)
		}

	})
}
