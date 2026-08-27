package commands

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/craftybase/stocksmith-cli/internal/output"
)

func sampleMaterialJSON() json.RawMessage {
	return json.RawMessage(`{
		"id": 123, "name": "Organic Beeswax", "sku": "WAX-001",
		"category": "Waxes", "stock_on_hand": "12.5", "on_order": "8.0", "unit_measure": "kg",
		"unit_cost": {"amount": "8.75", "currency_code": "USD"}
	}`)
}

func TestRenderMaterialShow_SingleRowTable(t *testing.T) {
	var buf bytes.Buffer
	if err := renderMaterialShow(&buf, sampleMaterialJSON(), false); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"ID", "NAME", "SKU", "CATEGORY", "ON HAND", "ON ORDER", "UNIT COST", "123", "Organic Beeswax", "WAX-001", "Waxes", "12.5 kg", "8.0 kg", "$8.75"} {
		if !strings.Contains(out, want) {
			t.Errorf("show output missing %q\n---\n%s", want, out)
		}
	}
}

func TestMaterialToRow_NilCostEmptySKUNilCategoryRenderDash(t *testing.T) {
	m := Material{ID: 1, Name: "No Price", StockOnHand: "0", UnitCost: (*output.Money)(nil)}
	row := materialToRow(&m)
	// ID, NAME, SKU, CATEGORY, ON HAND, ON ORDER, UNIT COST
	if row[2] != "—" {
		t.Errorf("empty SKU should render —, got %q", row[2])
	}
	if row[3] != "—" {
		t.Errorf("nil category should render —, got %q", row[3])
	}
	if row[6] != "—" {
		t.Errorf("nil unit_cost should render —, got %q", row[6])
	}
}

// on_order is absent on an API that predates the field. It must render "—", not
// an invented zero, and must never be confused with the on-hand figure.
func TestMaterialToRow_AbsentOnOrderRendersDash(t *testing.T) {
	row := materialToRow(&Material{ID: 1, Name: "Legacy API", StockOnHand: "12.5", UnitMeasure: "kg"})
	if row[4] != "12.5 kg" {
		t.Errorf("on hand: want %q, got %q", "12.5 kg", row[4])
	}
	if row[5] != "—" {
		t.Errorf("absent on_order should render —, got %q", row[5])
	}
}

// on_order is reported in stock units, so it carries the same unit of measure as
// stock_on_hand — never the purchase unit.
func TestMaterialToRow_OnOrderUsesStockUnit(t *testing.T) {
	row := materialToRow(&Material{ID: 2, Name: "Beef Tallow", StockOnHand: "100.0", OnOrder: "8.0", UnitMeasure: "g"})
	if row[5] != "8.0 g" {
		t.Errorf("on order: want %q, got %q", "8.0 g", row[5])
	}
}

func TestWithUnit(t *testing.T) {
	for _, tc := range []struct{ qty, unit, want string }{
		{"12.5", "kg", "12.5 kg"},
		{"0.0", "kg", "0.0 kg"},
		{"12.5", "", "12.5"},
		{"", "kg", "—"},
		{"", "", "—"},
	} {
		if got := withUnit(tc.qty, tc.unit); got != tc.want {
			t.Errorf("withUnit(%q, %q): want %q, got %q", tc.qty, tc.unit, tc.want, got)
		}
	}
}
