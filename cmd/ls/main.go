package main

import (
	. "github.com/troesler95/go-cli-tools"
)

func main() {
	printer := NewColorizedPrinter(NewColor(FgBlue, BgGreen, DefaultText))
	printer.Println("Hello, world!")
}
