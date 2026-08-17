package settle

import (
	"context"
	"testing"

	"goldbar/internal/model"
)

// TestRunBatchOfExactTierKarats settles a batch whose karats are all exactly
// equal to a configured discount tier and asserts the per-order amounts and the
// store daily totals.
func TestRunBatchOfExactTierKarats(t *testing.T) {
	cfg := testConfig()
	orders := []model.Order{
		mkOrder(2, "A1", 999, 2, 8, 768.5),
		mkOrder(3, "A2", 990, 3.5, 7, 768.5),
		mkOrder(4, "A3", 900, 4, 10, 768.5),
	}

	res, err := Run(context.Background(), cfg, orders, 2)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected row errors: %v", res.Errors)
	}
	if len(res.Settlements) != 3 {
		t.Fatalf("settlements = %d, want 3", len(res.Settlements))
	}

	wantDiscount := []float64{1504.75, 2529.71, 2489.94}
	wantPayable := []float64{4739.25, 2933.79, 5315.06}
	for i, s := range res.Settlements {
		if s.OldDiscount != wantDiscount[i] {
			t.Errorf("%s: old_discount = %v, want %v", s.OrderID, s.OldDiscount, wantDiscount[i])
		}
		if s.Payable != wantPayable[i] {
			t.Errorf("%s: payable = %v, want %v", s.OrderID, s.Payable, wantPayable[i])
		}
	}

	if res.Summary.TotalOldDiscount != 6524.40 {
		t.Errorf("total_old_discount = %v, want 6524.40", res.Summary.TotalOldDiscount)
	}
	if res.Summary.TotalNewSalePrice != 19512.50 {
		t.Errorf("total_new_sale_price = %v, want 19512.50", res.Summary.TotalNewSalePrice)
	}
	if res.Summary.TotalCraftFee != 300.00 {
		t.Errorf("total_craft_fee = %v, want 300.00", res.Summary.TotalCraftFee)
	}
	if res.Summary.TotalPayable != 12988.10 {
		t.Errorf("total_payable = %v, want 12988.10", res.Summary.TotalPayable)
	}
	if res.Summary.NetStoreReceivable != 12988.10 {
		t.Errorf("net_store_receivable = %v, want 12988.10", res.Summary.NetStoreReceivable)
	}
}

// TestRunBatchOfBetweenTierKarats is the contrast case: karats that fall
// strictly between two configured tiers.
func TestRunBatchOfBetweenTierKarats(t *testing.T) {
	cfg := testConfig()
	orders := []model.Order{
		mkOrder(2, "B1", 995, 2, 8, 768.5),
		mkOrder(3, "B2", 950, 2.5, 8, 768.5),
	}

	res, err := Run(context.Background(), cfg, orders, 2)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(res.Settlements) != 2 {
		t.Fatalf("settlements = %d, want 2", len(res.Settlements))
	}
	// 995 uses the 990 tier (rate 0.95): 2 * 768.5 * 0.995 * 0.95
	if res.Settlements[0].OldDiscount != 1452.85 {
		t.Errorf("B1 old_discount = %v, want 1452.85", res.Settlements[0].OldDiscount)
	}
	// 950 uses the 900 tier (rate 0.90): 2.5 * 768.5 * 0.950 * 0.90
	if res.Settlements[1].OldDiscount != 1642.67 {
		t.Errorf("B2 old_discount = %v, want 1642.67", res.Settlements[1].OldDiscount)
	}
}
