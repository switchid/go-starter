package stringers

import (
	"GoStarter/pkg/utils/consts"
)

type ColorText int

const (
	RED ColorText = iota
	GREEN
	BLUE
	CYAN
	MAGENTA
	YELLOW
	WHITE
)

type String struct {
	Text string
}

func NewString(s string) *String {
	return &String{Text: s}
}

func (s *String) String() string {
	return s.Text
}

func (s *String) SetTextColor(color ColorText) string {
	var c consts.Colors
	switch color {
	case RED:
		return c.GetRed() + s.Text + c.GetWhite()
	case GREEN:
		return c.GetGreen() + s.Text + c.GetWhite()
	case BLUE:
		return c.GetBlue() + s.Text + c.GetWhite()
	case CYAN:
		return c.GetCyan() + s.Text + c.GetWhite()
	case MAGENTA:
		return c.GetMagenta() + s.Text + c.GetWhite()
	case YELLOW:
		return c.GetYellow() + s.Text + c.GetWhite()
	case WHITE:
		return c.GetWhite() + s.Text
	default:
		return s.Text
	}
}
