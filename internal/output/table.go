package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const placeholder = "—"

func FormatTable(w io.Writer, headers []string, rows [][]string, useColor bool) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	if useColor {
		boldHeaders := make([]string, len(headers))
		for i, h := range headers {
			boldHeaders[i] = bold(h)
		}
		fmt.Fprintln(tw, strings.Join(boldHeaders, "\t"))
	} else {
		fmt.Fprintln(tw, strings.Join(headers, "\t"))
	}

	for _, row := range rows {
		if useColor {
			coloredRow := make([]string, len(row))
			for i, cell := range row {
				if cell == placeholder {
					coloredRow[i] = dim(cell)
				} else {
					coloredRow[i] = cell
				}
			}
			fmt.Fprintln(tw, strings.Join(coloredRow, "\t"))
		} else {
			fmt.Fprintln(tw, strings.Join(row, "\t"))
		}
	}

	tw.Flush()
}

func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func dim(s string) string {
	return "\033[2m" + s + "\033[0m"
}
