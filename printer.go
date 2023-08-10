package goclitools

import "fmt"

type Printer interface {
	Printf(format string, a ...any) (bytesWritten int, err error)
	Println(a ...any) (bytesWritten int, err error)
	Print(a ...any) (bytesWritten int, err error)
}

type colorizedPrinter struct {
	color Color
}

func NewColorizedPrinter(defaultColor Color) Printer {
	return &colorizedPrinter{
		color: defaultColor,
	}
}

func NewDefaultPrinter() Printer {
	return &colorizedPrinter{
		color: NewDefaultColor(),
	}
}

func (cp *colorizedPrinter) Printf(format string, a ...any) (bytesWritten int, err error) {
	text := fmt.Sprintf(format, a...)
	return fmt.Printf(Colorize(text, cp.color))
}

func (cp *colorizedPrinter) Println(a ...any) (bytesWritten int, err error) {
	text := fmt.Sprint(a...)
	return fmt.Println(Colorize(text, cp.color))
}

func (cp *colorizedPrinter) Print(a ...any) (bytesWritten int, err error) {
	text := fmt.Sprint(a...)
	return fmt.Print(Colorize(text, cp.color))
}

func (cp *colorizedPrinter) Color(newColor Color) {
	cp.color = newColor
}
