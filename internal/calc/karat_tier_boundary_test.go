package calc

import (
	"testing"

	"goldbar/internal/model"
)

// TestDiscountRateOnExactTierKarat asserts the discount rate for karats that
// are exactly equal to a configured tier: such a karat must use that tier's own
// rate, not the rate of a lower tier.
func TestDiscountRateOnExactTierKarat(t *testing.T) {
	cfg := sampleConfig()
	cases := []struct {
		karat int
		want  float64
	}{
		{karat: 999, want: 0.98},
		{karat: 990, want: 0.95},
		{karat: 900, want: 0.90},
		{karat: 750, want: 0.80},
	}
	for _, c := range cases {
		got, ok := DiscountRate(cfg.KaratTiers, c.karat)
		if !ok {
			t.Errorf("karat %d: no discount tier matched", c.karat)
			continue
		}
		if got != c.want {
			t.Errorf("karat %d: rate = %v, want %v", c.karat, got, c.want)
		}
	}
}

// TestDiscountRateBetweenTiersKarat is the contrast case: karats that fall
// strictly between two tiers must keep using the nearest lower tier.
func TestDiscountRateBetweenTiersKarat(t *testing.T) {
	cfg := sampleConfig()
	cases := []struct {
		karat int
		want  float64
	}{
		{karat: 1000, want: 0.98},
		{karat: 995, want: 0.95},
		{karat: 950, want: 0.90},
		{karat: 800, want: 0.80},
	}
	for _, c := range cases {
		got, ok := DiscountRate(cfg.KaratTiers, c.karat)
		if !ok {
			t.Errorf("karat %d: no discount tier matched", c.karat)
			continue
		}
		if got != c.want {
			t.Errorf("karat %d: rate = %v, want %v", c.karat, got, c.want)
		}
	}
}

// TestComputeOnExactTierKarat asserts the monetary outcome for the three karats
// that stores actually trade in most (999 / 990 / 900), each of which is exactly
// a configured tier.
func TestComputeOnExactTierKarat(t *testing.T) {
	cfg := sampleConfig()
	cases := []struct {
		name            string
		order           model.Order
		wantOldDiscount float64
		wantSalePrice   float64
		wantPayable     float64
	}{
		{
			name:            "karat 999",
			order:           model.Order{LineNumber: 2, OrderID: "E999", Customer: "C", OldKarat: 999, OldWeight: 2, NewProductCode: "G", NewWeight: 8, GoldPrice: 768.5},
			wantOldDiscount: 1504.75,
			wantSalePrice:   6244.00,
			wantPayable:     4739.25,
		},
		{
			name:            "karat 990",
			order:           model.Order{LineNumber: 3, OrderID: "E990", Customer: "C", OldKarat: 990, OldWeight: 3.5, NewProductCode: "G", NewWeight: 7, GoldPrice: 768.5},
			wantOldDiscount: 2529.71,
			wantSalePrice:   5463.50,
			wantPayable:     2933.79,
		},
		{
			name:            "karat 900",
			order:           model.Order{LineNumber: 4, OrderID: "E900", Customer: "C", OldKarat: 900, OldWeight: 4, NewProductCode: "G", NewWeight: 10, GoldPrice: 768.5},
			wantOldDiscount: 2489.94,
			wantSalePrice:   7805.00,
			wantPayable:     5315.06,
		},
	}
	for _, c := range cases {
		s, err := Compute(cfg, c.order)
		if err != nil {
			t.Errorf("%s: Compute: %v", c.name, err)
			continue
		}
		if s.OldDiscount != c.wantOldDiscount {
			t.Errorf("%s: old_discount = %v, want %v", c.name, s.OldDiscount, c.wantOldDiscount)
		}
		if s.NewSalePrice != c.wantSalePrice {
			t.Errorf("%s: new_sale_price = %v, want %v", c.name, s.NewSalePrice, c.wantSalePrice)
		}
		if s.Payable != c.wantPayable {
			t.Errorf("%s: payable = %v, want %v", c.name, s.Payable, c.wantPayable)
		}
	}
}
