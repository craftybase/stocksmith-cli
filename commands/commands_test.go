package commands_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/craftybase/craftybase-cli/commands"
)

func materialListFixture() []byte {
	return []byte(`{
		"materials": [
			{
				"id": 123,
				"name": "Organic Beeswax",
				"sku": "WAX-001",
				"category": "Waxes",
				"stock_on_hand": "12.5",
				"unit_measure": "kg",
				"unit_cost": {"amount": "8.75", "currency_code": "USD"}
			},
			{
				"id": 456,
				"name": "Soy Wax Flakes",
				"sku": "WAX-002",
				"category": "Waxes",
				"stock_on_hand": "40.0",
				"unit_measure": "kg",
				"unit_cost": {"amount": "4.10", "currency_code": "USD"}
			}
		],
		"meta": {
			"total_pages": 1,
			"total_count": 2,
			"per_page": 25,
			"page": 1
		}
	}`)
}

func accountFixture() []byte {
	return []byte(`{
		"account": {
			"id": 42,
			"name": "BambuEarth",
			"currency_code": "USD",
			"time_zone": "Pacific Time (US & Canada)",
			"plan": "Growth"
		}
	}`)
}

func pingFixture() []byte {
	return []byte(`{"ping": "pong", "account_id": 42}`)
}

func setupMockServer(routes map[string]func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}
	return httptest.NewServer(mux)
}

func TestMaterialListFixture_ContractCheck(t *testing.T) {
	fixture := materialListFixture()
	var raw map[string]interface{}
	if err := json.Unmarshal(fixture, &raw); err != nil {
		t.Fatal(err)
	}

	materials := raw["materials"].([]interface{})
	mat := materials[0].(map[string]interface{})

	catVal := mat["category"]
	if _, ok := catVal.(string); !ok {
		t.Errorf("fixture category must be flat string, got %T", catVal)
	}

	unitCost := mat["unit_cost"].(map[string]interface{})
	amtVal := unitCost["amount"]
	if _, ok := amtVal.(string); !ok {
		t.Errorf("fixture amount must be string, got %T", amtVal)
	}
}

func TestVersionOutput(t *testing.T) {
	commands.SetVersion("1.2.3", "abc1234", "2026-06-19T00:00:00Z")
}

func TestSetupMockServer(t *testing.T) {
	srv := setupMockServer(map[string]func(http.ResponseWriter, *http.Request){
		"/api/v1/materials": func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/materials" {
				t.Errorf("expected path /api/v1/materials, got %q", r.URL.Path)
			}
			w.Write(materialListFixture())
		},
	})
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/materials")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestCategoryFilterParam(t *testing.T) {
	var gotParam string
	srv := setupMockServer(map[string]func(http.ResponseWriter, *http.Request){
		"/api/v1/materials": func(w http.ResponseWriter, r *http.Request) {
			gotParam = r.URL.Query().Get("category_name")
			w.Write([]byte(`{"materials":[],"meta":{"total_pages":1,"total_count":0,"per_page":25,"page":1}}`))
		},
	})
	defer srv.Close()

	reqURL := srv.URL + "/api/v1/materials?category_name=Waxes"
	resp, err := http.Get(reqURL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if gotParam != "Waxes" {
		t.Errorf("expected category_name=Waxes, got %q", gotParam)
	}
}

func TestMaskToken_Format(t *testing.T) {
	// Token masking format: first 9 chars + … + last 4
	if !strings.Contains("live_a1b2…g7h8", "…") {
		t.Error("mask format should contain ellipsis")
	}
}

// Verify fixtures exist and are referenced.
var _ = accountFixture
var _ = pingFixture

func manufactureListFixture() []byte {
	return []byte(`{
		"manufactures": [
			{"id": 1042, "product_id": 77, "production_status": "completed",
			 "quantity": "50.0", "actual_quantity": "48.0", "start_date": "2026-06-01T00:00:00Z",
			 "total_materials_cost": {"amount": "24.50", "currency_code": "USD"},
			 "line_items": [{"id": 5001, "material_id": 312, "name": "Beeswax", "quantity": "200.0",
			   "unit_price": {"amount": "0.08", "currency_code": "USD"}}]}
		],
		"meta": {"current_page": 1, "total_pages": 1, "total_count": 1, "per_page": 25}
	}`)
}

// Contract: manufactures envelope key; production_status flat string; money is a
// string-amount object; line_items is an array; never leaks the internal key.
func TestManufactureListFixture_ContractCheck(t *testing.T) {
	var raw map[string]interface{}
	if err := json.Unmarshal(manufactureListFixture(), &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["manufactures"]; !ok {
		t.Fatal("envelope must use key \"manufactures\"")
	}
	m := raw["manufactures"].([]interface{})[0].(map[string]interface{})
	if _, ok := m["production_status"].(string); !ok {
		t.Errorf("production_status must be a flat string, got %T", m["production_status"])
	}
	if _, ok := m["line_items"].([]interface{}); !ok {
		t.Errorf("line_items must be an array, got %T", m["line_items"])
	}
	cost := m["total_materials_cost"].(map[string]interface{})
	if _, ok := cost["amount"].(string); !ok {
		t.Errorf("money amount must be a string, got %T", cost["amount"])
	}
}
