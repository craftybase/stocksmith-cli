package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/craftybase/craftybase-cli/internal/output"
)

func TestFormatMoney_USD(t *testing.T) {
	m := &output.Money{Amount: "8.75", CurrencyCode: "USD"}
	got := output.FormatMoney(m)
	if got != "$8.75" {
		t.Errorf("expected $8.75, got %q", got)
	}
}

func TestFormatMoney_GBP(t *testing.T) {
	m := &output.Money{Amount: "10.00", CurrencyCode: "GBP"}
	got := output.FormatMoney(m)
	if got != "£10.00" {
		t.Errorf("expected £10.00, got %q", got)
	}
}

func TestFormatMoney_Unknown(t *testing.T) {
	m := &output.Money{Amount: "8.75", CurrencyCode: "CHF"}
	got := output.FormatMoney(m)
	if got != "CHF 8.75" {
		t.Errorf("expected 'CHF 8.75', got %q", got)
	}
}

func TestFormatMoney_Absent(t *testing.T) {
	got := output.FormatMoney(nil)
	if got != "—" {
		t.Errorf("expected '—', got %q", got)
	}
}

func TestFormatMoney_StringAmount(t *testing.T) {
	m := &output.Money{Amount: "8.75000", CurrencyCode: "USD"}
	got := output.FormatMoney(m)
	if got != "$8.75000" {
		t.Errorf("expected $8.75000 (string preserved), got %q", got)
	}
}

func TestColorDisabled_NO_COLOR(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	enabled := output.ColorEnabled(false)
	if enabled {
		t.Error("expected color disabled when NO_COLOR is set")
	}
}

func TestColorDisabled_NoColorFlag(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	enabled := output.ColorEnabled(true)
	if enabled {
		t.Error("expected color disabled when --no-color flag is set")
	}
}

func TestFormatTable_NoColor(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"ID", "NAME", "SKU"}
	rows := [][]string{
		{"1", "Beeswax", "WAX-001"},
		{"2", "Soy Wax", "—"},
	}
	output.FormatTable(&buf, headers, rows, false)
	got := buf.String()

	if !strings.Contains(got, "ID") {
		t.Error("expected ID in output")
	}
	if !strings.Contains(got, "Beeswax") {
		t.Error("expected Beeswax in output")
	}
	if !strings.Contains(got, "—") {
		t.Error("expected placeholder in output")
	}
	if strings.Contains(got, "\033[") {
		t.Error("expected no ANSI codes in no-color mode")
	}
}

func TestPrintJSON(t *testing.T) {
	input := []byte(`{"material":{"id":1,"name":"Beeswax"}}`)
	var buf bytes.Buffer
	if err := output.PrintJSON(&buf, input); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, `"material"`) {
		t.Error("expected material key in JSON output")
	}
}

func TestMaterialsList_GoldenFixtureContracts(t *testing.T) {
	fixtureJSON := []byte(`{
		"materials": [
			{
				"id": 123,
				"name": "Organic Beeswax",
				"sku": "WAX-001",
				"category": "Waxes",
				"stock_on_hand": "12.5",
				"unit_measure": "kg",
				"unit_cost": {"amount": "8.75", "currency_code": "USD"}
			}
		],
		"meta": {
			"total_pages": 1,
			"total_count": 1,
			"per_page": 25,
			"page": 1
		}
	}`)

	var raw map[string]interface{}
	if err := json.Unmarshal(fixtureJSON, &raw); err != nil {
		t.Fatal(err)
	}

	materials, ok := raw["materials"].([]interface{})
	if !ok || len(materials) == 0 {
		t.Fatal("expected materials array")
	}

	mat := materials[0].(map[string]interface{})

	catVal := mat["category"]
	if _, isStr := catVal.(string); !isStr {
		t.Errorf("category must be a flat string, got %T: %v", catVal, catVal)
	}

	unitCost, ok := mat["unit_cost"].(map[string]interface{})
	if !ok {
		t.Fatal("expected unit_cost object")
	}
	amtVal := unitCost["amount"]
	if _, isStr := amtVal.(string); !isStr {
		t.Errorf("unit_cost.amount must be a string, got %T: %v", amtVal, amtVal)
	}
}
