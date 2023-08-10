package goclitools

import "fmt"

// a terminal text output color
type FgColor = int
type BgColor = int

const resetFormatting string = "\033[0m"

// enumerates possible text output colors
const (
	FgBlack FgColor = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgPurple
	FgCyan
	FgWhite
	_ // increment iota
	FgDefault
)

const (
	BgBlack BgColor = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgPurple
	BgCyan
	BgWhite
	_ // increment iota
	BgDefault
)

type TextModifier uint32

const (
	DefaultText   TextModifier = 0
	ItalicizeText TextModifier = 1 << iota
	BoldText
	UnderlineText
)

type Color struct {
	Foreground FgColor
	Background BgColor
	italic     bool
	bold       bool
	underline  bool
}

func NewColor(fgColor FgColor, bgColor BgColor, modifiers TextModifier) Color {
	return Color{
		Foreground: fgColor,
		Background: bgColor,
		italic:     modifiers.hasModifier(ItalicizeText),
		bold:       modifiers.hasModifier(BoldText),
		underline:  modifiers.hasModifier(UnderlineText),
	}
}

func NewDefaultColor() Color {
	return NewColor(FgDefault, BgDefault, DefaultText)
}

func (m TextModifier) hasModifier(modifierToTest TextModifier) bool {
	return m&modifierToTest != 0
}

// adds one or more text modifiers
func (c *Color) AddTextModifiers(modifiers TextModifier) error {
	if modifiers == DefaultText {
		return fmt.Errorf("unable to add text modifier. modifier DefaultText is not valid")
	}

	if modifiers.hasModifier(ItalicizeText) {
		c.italic = true
	}
	if modifiers.hasModifier(BoldText) {
		c.bold = true
	}
	if modifiers.hasModifier(UnderlineText) {
		c.underline = true
	}

	return nil
}

func (c *Color) ClearTextModifiers(modifiers TextModifier) error {
	if modifiers == DefaultText {
		return fmt.Errorf("unable to add text modifier. modifier DefaultText is not valid")
	}

	if modifiers.hasModifier(ItalicizeText) {
		c.italic = false
	}
	if modifiers.hasModifier(BoldText) {
		c.bold = false
	}
	if modifiers.hasModifier(UnderlineText) {
		c.underline = false
	}

	return nil
}

func (c *Color) ClearAllTextModifiers() {
	c.bold = false
	c.italic = false
	c.underline = false
}

// returns a string formatted with ANSI escape codes for colorized output
func Colorize(msg string, color Color) string {
	boldModifier := 22
	if color.bold {
		boldModifier = 1
	}

	italicModifier := 23
	if color.italic {
		italicModifier = 3
	}

	underlineModifier := 24
	if color.underline {
		underlineModifier = 4
	}

	return fmt.Sprintf("\033[%d;%d;%d;%d;%dm%s%s",
		color.Foreground,
		color.Background,
		boldModifier,
		italicModifier,
		underlineModifier,
		msg,
		resetFormatting,
	)
}
