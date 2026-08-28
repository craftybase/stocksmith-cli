package commands

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/craftybase/stocksmith-cli/internal/brand"
	"github.com/craftybase/stocksmith-cli/internal/output"
)

// Material is one stock item. on_order is the inbound quantity still outstanding
// on purchases, in the same stock unit as stock_on_hand (not purchase units), and
// is ungated — unlike unit_cost it is visible without the material-costs
// permission. It is absent on an API that predates the field, which withUnit
// renders as "—" rather than inventing a zero.
type Material struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	SKU         string        `json:"sku"`
	Category    *string       `json:"category"`
	StockOnHand string        `json:"stock_on_hand"`
	OnOrder     string        `json:"on_order"`
	UnitMeasure string        `json:"unit_measure"`
	UnitCost    *output.Money `json:"unit_cost"`
}

var materialsCmd = &cobra.Command{
	Use:   "materials",
	Short: "Manage materials",
}

var (
	materialsFilters    projectFilters
	materialsPagination paginationFlags
)

func materialsToTable(rawItems []json.RawMessage) ([]string, [][]string) {
	headers := []string{"ID", "NAME", "SKU", "CATEGORY", "ON HAND", "ON ORDER", "UNIT COST"}
	rows := make([][]string, 0, len(rawItems))
	for i, raw := range rawItems {
		var m Material
		if err := json.Unmarshal(raw, &m); err != nil {
			warnSkip(i, err)
			continue
		}
		rows = append(rows, materialToRow(&m))
	}
	return headers, rows
}

// withUnit renders a quantity in the material's unit of measure. An absent
// quantity renders "—" — a material that reports no figure is not the same as
// one reporting zero.
func withUnit(qty, unit string) string {
	if qty == "" {
		return "—"
	}
	if unit == "" {
		return qty
	}
	return qty + " " + unit
}

func materialToRow(m *Material) []string {
	sku := m.SKU
	if sku == "" {
		sku = "—"
	}
	category := "—"
	if m.Category != nil && *m.Category != "" {
		category = *m.Category
	}
	unitCost := output.FormatMoney(m.UnitCost)

	return []string{
		strconv.Itoa(m.ID),
		m.Name,
		sku,
		category,
		withUnit(m.StockOnHand, m.UnitMeasure),
		withUnit(m.OnOrder, m.UnitMeasure),
		unitCost,
	}
}

// renderMaterialShow adapts the one-row materials table to the
// resourceConfig.renderShow signature.
func renderMaterialShow(w io.Writer, raw json.RawMessage, useColor bool) error {
	headers, rows := materialsToTable([]json.RawMessage{raw})
	output.FormatTable(w, headers, rows, useColor)
	return nil
}

func init() {
	res := resourceConfig{
		pathSegment: "materials",
		collection:  "materials",
		singular:    "material",
		listLong: "List materials from your " + brand.ProductName + " account.\n\n" +
			"ON ORDER is the quantity still inbound on outstanding purchases, in the same\n" +
			"stock unit as ON HAND. Filter by SKU, name, category, or state. Use --all to\n" +
			"fetch all pages, or --ndjson for streaming NDJSON output suitable for data\n" +
			"pipelines.",
		toTable:    materialsToTable,
		renderShow: renderMaterialShow,
	}
	materialsCmd.AddCommand(newResourceListCmd(res, &materialsFilters, &materialsPagination))
	materialsCmd.AddCommand(newResourceShowCmd(res))
	rootCmd.AddCommand(materialsCmd)
}
