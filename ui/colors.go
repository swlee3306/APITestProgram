package ui

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Color string

const (
	Reset  Color = "\033[0m"
	Red    Color = "\033[31m"
	Green  Color = "\033[32m"
	Yellow Color = "\033[33m"
	Blue   Color = "\033[34m"
	Purple Color = "\033[35m"
	Cyan   Color = "\033[36m"
	White  Color = "\033[37m"
	Gray   Color = "\033[90m"
	
	// Bright colors
	BrightRed    Color = "\033[91m"
	BrightGreen  Color = "\033[92m"
	BrightYellow Color = "\033[93m"
	BrightBlue   Color = "\033[94m"
	BrightPurple Color = "\033[95m"
	BrightCyan   Color = "\033[96m"
	BrightWhite  Color = "\033[97m"
	
	// Background colors
	BgRed    Color = "\033[41m"
	BgGreen  Color = "\033[42m"
	BgYellow Color = "\033[43m"
	BgBlue   Color = "\033[44m"
	BgPurple Color = "\033[45m"
	BgCyan   Color = "\033[46m"
	BgWhite  Color = "\033[47m"
	
	// Styles
	Bold      Color = "\033[1m"
	Dim       Color = "\033[2m"
	Italic    Color = "\033[3m"
	Underline Color = "\033[4m"
	Blink     Color = "\033[5m"
	Reverse   Color = "\033[7m"
	Strike    Color = "\033[9m"
)

type ColorPrinter struct {
	enabled bool
}

func NewColorPrinter() *ColorPrinter {
	return &ColorPrinter{
		enabled: isColorSupported(),
	}
}

func (cp *ColorPrinter) Enable() {
	cp.enabled = true
}

func (cp *ColorPrinter) Disable() {
	cp.enabled = false
}

func (cp *ColorPrinter) IsEnabled() bool {
	return cp.enabled
}

func (cp *ColorPrinter) Print(color Color, format string, args ...interface{}) {
	if cp.enabled {
		fmt.Printf("%s%s%s", color, fmt.Sprintf(format, args...), Reset)
	} else {
		fmt.Printf(format, args...)
	}
}

func (cp *ColorPrinter) Println(color Color, format string, args ...interface{}) {
	cp.Print(color, format+"\n", args...)
}

func (cp *ColorPrinter) Success(format string, args ...interface{}) {
	cp.Println(Green, "✅ "+format, args...)
}

func (cp *ColorPrinter) Error(format string, args ...interface{}) {
	cp.Println(Red, "❌ "+format, args...)
}

func (cp *ColorPrinter) Warning(format string, args ...interface{}) {
	cp.Println(Yellow, "⚠️  "+format, args...)
}

func (cp *ColorPrinter) Info(format string, args ...interface{}) {
	cp.Println(Blue, "ℹ️  "+format, args...)
}

func (cp *ColorPrinter) Debug(format string, args ...interface{}) {
	cp.Println(Gray, "🐛 "+format, args...)
}

func (cp *ColorPrinter) Header(format string, args ...interface{}) {
	cp.Println(Bold+Blue, "🚀 "+format, args...)
}

func (cp *ColorPrinter) SubHeader(format string, args ...interface{}) {
	cp.Println(Bold+Cyan, "  "+format, args...)
}

func (cp *ColorPrinter) Highlight(format string, args ...interface{}) {
	cp.Println(Bold+BrightYellow, "⭐ "+format, args...)
}

func (cp *ColorPrinter) Code(format string, args ...interface{}) {
	cp.Println(Gray+BgWhite, "  "+format+"  ", args...)
}

func (cp *ColorPrinter) URL(format string, args ...interface{}) {
	cp.Println(Underline+Blue, format, args...)
}

func (cp *ColorPrinter) Method(method string) {
	var color Color
	switch strings.ToUpper(method) {
	case "GET":
		color = Green
	case "POST":
		color = Blue
	case "PUT":
		color = Yellow
	case "DELETE":
		color = Red
	case "PATCH":
		color = Purple
	default:
		color = White
	}
	cp.Print(Bold+color, " %s ", strings.ToUpper(method))
}

func (cp *ColorPrinter) StatusCode(code int) {
	var color Color
	switch {
	case code >= 200 && code < 300:
		color = Green
	case code >= 300 && code < 400:
		color = Yellow
	case code >= 400 && code < 500:
		color = Red
	case code >= 500:
		color = BrightRed
	default:
		color = White
	}
	cp.Print(Bold+color, " %d ", code)
}

func (cp *ColorPrinter) Duration(duration time.Duration) {
	var color Color
	switch {
	case duration < 100*time.Millisecond:
		color = Green
	case duration < 500*time.Millisecond:
		color = Yellow
	default:
		color = Red
	}
	cp.Print(color, " %v ", duration)
}

func (cp *ColorPrinter) Size(size int64) {
	var color Color
	switch {
	case size < 1024:
		color = Green
	case size < 1024*1024:
		color = Yellow
	default:
		color = Red
	}
	cp.Print(color, " %d bytes ", size)
}

func (cp *ColorPrinter) Progress(current, total int, message string) {
	if total == 0 {
		return
	}
	
	percentage := float64(current) / float64(total) * 100
	barWidth := 30
	filledWidth := int(float64(barWidth) * percentage / 100)
	
	bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)
	
	cp.Print(Cyan, "\r%s [%s] %.1f%% (%d/%d)", message, bar, percentage, current, total)
	
	if current >= total {
		cp.Println(Reset, "")
	}
}

func (cp *ColorPrinter) Spinner(message string, done chan bool) {
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	
	go func() {
		i := 0
		for {
			select {
			case <-done:
				cp.Print(Reset, "\r%s ✅\n", message)
				return
			default:
				cp.Print(Cyan, "\r%s %s", message, spinnerChars[i%len(spinnerChars)])
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

func (cp *ColorPrinter) Table(header []string, rows [][]string) {
	if len(rows) == 0 {
		return
	}
	
	// Calculate column widths
	widths := make([]int, len(header))
	for i, h := range header {
		widths[i] = len(h)
	}
	
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	
	// Print header
	cp.Print(Bold+Blue, "┌")
	for i, h := range header {
		cp.Print(Bold+Blue, "─%s─", strings.Repeat("─", widths[i]))
		if i < len(header)-1 {
			cp.Print(Bold+Blue, "┬")
		}
	}
	cp.Println(Bold+Blue, "┐")
	
	cp.Print(Bold+Blue, "│")
	for i, h := range header {
		cp.Print(Bold+Blue, " %-*s │", widths[i], h)
	}
	cp.Println(Reset, "")
	
	cp.Print(Bold+Blue, "├")
	for i := range header {
		cp.Print(Bold+Blue, "─%s─", strings.Repeat("─", widths[i]))
		if i < len(header)-1 {
			cp.Print(Bold+Blue, "┼")
		}
	}
	cp.Println(Bold+Blue, "┤")
	
	// Print rows
	for _, row := range rows {
		cp.Print(Reset, "│")
		for i, cell := range row {
			cp.Print(Reset, " %-*s │", widths[i], cell)
		}
		cp.Println(Reset, "")
	}
	
	cp.Print(Bold+Blue, "└")
	for i := range header {
		cp.Print(Bold+Blue, "─%s─", strings.Repeat("─", widths[i]))
		if i < len(header)-1 {
			cp.Print(Bold+Blue, "┴")
		}
	}
	cp.Println(Bold+Blue, "┘")
}

func (cp *ColorPrinter) JSON(data string) {
	// Simple JSON syntax highlighting
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.Contains(line, "\"") {
			// Highlight strings
			parts := strings.Split(line, "\"")
			for i, part := range parts {
				if i%2 == 1 {
					cp.Print(Green, "\"%s\"", part)
				} else {
					cp.Print(Reset, "%s", part)
				}
			}
			cp.Println(Reset, "")
		} else {
			cp.Println(Reset, line)
		}
	}
}

func isColorSupported() bool {
	// Check if we're in a terminal
	if os.Getenv("TERM") == "" {
		return false
	}
	
	// Check for NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	
	// Check if output is redirected
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	
	return (stat.Mode() & os.ModeCharDevice) != 0
}
