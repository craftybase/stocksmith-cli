package output

import "fmt"

type Money struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

var currencySymbols = map[string]string{
	"USD": "$",
	"GBP": "£",
	"EUR": "€",
	"AUD": "A$",
	"CAD": "C$",
	"NZD": "NZ$",
}

func FormatMoney(m *Money) string {
	if m == nil {
		return "—"
	}
	sym, ok := currencySymbols[m.CurrencyCode]
	if ok {
		return sym + m.Amount
	}
	return fmt.Sprintf("%s %s", m.CurrencyCode, m.Amount)
}
