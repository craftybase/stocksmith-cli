package output

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsTTY() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return isatty.IsTerminal(os.Stdout.Fd())
}

func ColorEnabled(noColor bool) bool {
	if noColor {
		return false
	}
	return IsTTY()
}
