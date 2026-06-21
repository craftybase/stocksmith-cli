package commands

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

func sampleManufactureJSON() json.RawMessage {
	return json.RawMessage(`{
		"id": 1042, "product_id": 77, "variation_id": null,
		"batch_code": "BATCH-A", "production_status": "completed",
		"quantity": "50.0", "actual_quantity": "48.0", "minutes_worked": 120,
		"notes": null,
		"start_date": "2026-06-01T00:00:00Z", "completed_at": "2026-06-02T09:00:00Z",
		"deadline_at": null, "expiry_at": null,
		"total_materials_cost": {"amount": "24.50", "currency_code": "USD"},
		"total_labour_cost": {"amount": "18.00", "currency_code": "USD"},
		"line_items": [
			{"id": 5001, "material_id": 312, "name": "Beeswax", "sku": "BW-001",
			 "unit_measure": "grams", "quantity": "200.0", "lot_number": null,
			 "unit_price": {"amount": "0.08", "currency_code": "USD"}}
		]
	}`)
}

func TestManufactureDate(t *testing.T) {
	if got := manufactureDate("2026-06-01T00:00:00Z"); got != "2026-06-01" {
		t.Errorf("want 2026-06-01, got %q", got)
	}
	if got := manufactureDate(""); got != "—" {
		t.Errorf("empty should render —, got %q", got)
	}
}

func TestManufactureFilters_Apply(t *testing.T) {
	f := manufactureFilters{productID: "77", status: "completed", from: "2026-01-01", to: "2026-06-30"}
	params := url.Values{}
	f.apply(params)
	want := map[string]string{"product_id": "77", "status": "completed", "from": "2026-01-01", "to": "2026-06-30"}
	for k, v := range want {
		if params.Get(k) != v {
			t.Errorf("param %q: want %q, got %q", k, v, params.Get(k))
		}
	}
	if len(params) != len(want) {
		t.Errorf("expected exactly %d params, got %d: %v", len(want), len(params), params)
	}
	empty := url.Values{}
	(&manufactureFilters{}).apply(empty)
	if len(empty) != 0 {
		t.Errorf("empty filters should set no params, got %v", empty)
	}
}

func TestManufacturesToTable_Columns(t *testing.T) {
	headers, rows := manufacturesToTable([]json.RawMessage{sampleManufactureJSON()})
	wantHeaders := []string{"ID", "PRODUCT", "STATUS", "QTY", "START", "ITEMS", "MATERIALS COST", "LABOUR COST"}
	for i, h := range wantHeaders {
		if headers[i] != h {
			t.Errorf("header %d: want %q, got %q", i, h, headers[i])
		}
	}
	// ID, PRODUCT, STATUS, QTY, START, ITEMS, MATERIALS, LABOUR
	want := []string{"1042", "77", "completed", "50.0", "2026-06-01", "1", "$24.50", "$18.00"}
	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %d", len(rows))
	}
	for i := range want {
		if rows[0][i] != want[i] {
			t.Errorf("col %d: want %q, got %q", i, want[i], rows[0][i])
		}
	}
}

func TestManufacturesToTable_GatedCostsRenderDash(t *testing.T) {
	// view-only key: financial fields absent entirely
	raw := json.RawMessage(`{"id":1,"product_id":2,"production_status":"not_started","quantity":"5.0","start_date":"2026-06-01T00:00:00Z","line_items":[]}`)
	_, rows := manufacturesToTable([]json.RawMessage{raw})
	if rows[0][5] != "0" {
		t.Errorf("ITEMS should be 0, got %q", rows[0][5])
	}
	if rows[0][6] != "—" || rows[0][7] != "—" {
		t.Errorf("absent costs should render —, got materials=%q labour=%q", rows[0][6], rows[0][7])
	}
}

func TestRenderManufactureShow_DetailAndLineItems(t *testing.T) {
	var buf bytes.Buffer
	if err := renderManufactureShow(&buf, sampleManufactureJSON(), false); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{
		"ID", "1042", "PRODUCT", "77", "VARIATION", "—", "BATCH", "BATCH-A",
		"STATUS", "completed", "PLANNED QTY", "50.0", "ACTUAL QTY", "48.0",
		"MINUTES WORKED", "120", "START", "2026-06-01", "MATERIALS COST", "$24.50",
		"LABOUR COST", "$18.00", "NOTES", "—",
		"COMPLETED", "2026-06-02", "DEADLINE", "EXPIRY",
		"LINE ITEMS (1)", "MATERIAL", "Beeswax", "BW-001", "200.0 grams", "$0.08",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("show output missing %q\n---\n%s", want, out)
		}
	}
}

func TestRenderManufactureShow_NoLineItems(t *testing.T) {
	raw := json.RawMessage(`{"id":7,"product_id":3,"production_status":"not_started","quantity":"0","start_date":"","line_items":[]}`)
	var buf bytes.Buffer
	if err := renderManufactureShow(&buf, raw, false); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "LINE ITEMS (0)") || !strings.Contains(out, "No line items.") {
		t.Errorf("expected empty line-items state, got:\n%s", out)
	}
}
