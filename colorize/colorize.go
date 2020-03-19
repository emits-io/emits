package colorize

import (
	"fmt"
)

// Color string
type Color string

const (
	// Black Color
	Black Color = "\u001b[30"
	// Red Color
	Red Color = "\u001b[31"
	// Green Color
	Green Color = "\u001b[32"
	// Yellow Color
	Yellow Color = "\u001b[33"
	// Blue Color
	Blue Color = "\u001b[34"
	// Magenta Color
	Magenta Color = "\u001b[35"
	// Cyan Color
	Cyan Color = "\u001b[36"
	// White Color
	White Color = "\u001b[37"
)

func underline(text string) string {
	return fmt.Sprintf("\u001b[4m%s\u001b[0m", text)
}

func bold(text string) string {
	return fmt.Sprintf("\u001b[1m%s\u001b[0m", text)
}

// Printc func
func Printc(text string, color Color, bright bool) string {
	brightColor := ""
	if bright {
		brightColor = ";1"
	}
	return fmt.Sprintf("%s%sm%s\u001b[0m", color, brightColor, text)
}
