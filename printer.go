package goclitools

import (
	"fmt"
	"io"
	"os"
)

// prints to the console using specified colors
type ColorizedPrinter struct {
	color Color
}

// returns a new printer that prints output with the given color
func NewColorizedPrinter(defaultColor Color) *ColorizedPrinter {
	return &ColorizedPrinter{
		color: defaultColor,
	}
}

// returns a new printer that prints using the terminal defaults (typically no color)
func NewDefaultPrinter() *ColorizedPrinter {
	return &ColorizedPrinter{
		color: NewDefaultColor(),
	}
}

// print a string to console according to the given format and objects in color
func (cp *ColorizedPrinter) Printf(format string, a ...any) (bytesWritten int, err error) {
	return cp.Fprintf(os.Stdout, format, a...)
}

// print a string of the given objects to console with color and ends the line with a newline
func (cp *ColorizedPrinter) Println(a ...any) (bytesWritten int, err error) {
	return cp.Fprintln(os.Stdout, a...)
}

// print a string of the given objects to console with color
func (cp *ColorizedPrinter) Print(a ...any) (bytesWritten int, err error) {
	return cp.Fprint(os.Stdout, a...)
}

// print a string to console according to the given format and objects in the given color
func (cp *ColorizedPrinter) PrintfColor(c Color, format string, a ...any) (bytesWritten int, err error) {
	return cp.FprintfColor(os.Stdout, c, format, a...)
}

// print a string with the given objects to console with the given color and ends the line with a newline
func (cp *ColorizedPrinter) PrintlnColor(c Color, a ...any) (bytesWritten int, err error) {
	return cp.FprintlnColor(os.Stdout, c, a...)
}

// print a string with the given objects to console with the given color
func (cp *ColorizedPrinter) PrintColor(c Color, a ...any) (bytesWritten int, err error) {
	return cp.FprintColor(os.Stdout, c, a...)
}

// print a string to the given writer according to the given format and objects in the given color
func (cp *ColorizedPrinter) FprintfColor(w io.Writer, c Color, format string, a ...any) (bytesWritten int, err error) {
	currentColor := cp.color
	defer func() { cp.color = currentColor }()
	cp.color = c

	return cp.Fprintf(w, format, a...)
}

// print a string with the given objects as a string to the given writer with the given color and ends the line with a newline
func (cp *ColorizedPrinter) FprintlnColor(w io.Writer, c Color, a ...any) (bytesWritten int, err error) {
	currentColor := cp.color
	defer func() { cp.color = currentColor }()
	cp.color = c

	return cp.Fprintln(w, a...)
}

// print a string with the given objects to the given writer source with the given color
func (cp *ColorizedPrinter) FprintColor(w io.Writer, c Color, a ...any) (bytesWritten int, err error) {
	currentColor := cp.color
	defer func() { cp.color = currentColor }()
	cp.color = c

	return cp.Fprint(w, a...)
}

// print a string to the given writer according to the given format and objects in color
func (cp *ColorizedPrinter) Fprintf(w io.Writer, format string, a ...any) (bytesWritten int, err error) {
	text := fmt.Sprintf(format, a...)
	return fmt.Fprint(w, Colorize(text, cp.color))
}

// print a string with the given objects to the given writer in color and ends the line with a newline
func (cp *ColorizedPrinter) Fprintln(w io.Writer, a ...any) (bytesWritten int, err error) {
	text := fmt.Sprint(a...)
	return fmt.Fprintln(w, Colorize(text, cp.color))
}

// print a string with the given objects to the given writer source with color
func (cp *ColorizedPrinter) Fprint(w io.Writer, a ...any) (bytesWritten int, err error) {
	text := fmt.Sprint(a...)
	return fmt.Fprint(w, Colorize(text, cp.color))
}

// set the default color for this printer
func (cp *ColorizedPrinter) SetColor(newColor Color) {
	cp.color = newColor
}

// get the currently configured color for this printer
func (cp *ColorizedPrinter) GetColor() Color {
	return cp.color
}

// overrides the current colorizor color setting to red
// if msg is a blank string "", it is ignored and only the error is written
func (cp *ColorizedPrinter) PrintError(e error, msg string) (bytesWritten int, err error) {
	text := fmt.Sprintln(e)
	if len(msg) > 0 {
		text = fmt.Sprintf("%s: %s", msg, text)
	}

	if cp.color == NewDefaultColor() {
		return fmt.Fprint(os.Stderr, text)
	}
	return fmt.Fprint(os.Stderr, Colorize(text, NewColor(FgRed, BgDefault, DefaultText)))
}
