package helpers

import (
	"GoStarter/pkg/utils/consts"
	"strconv"
)

var Strings struct {
	text
	currency
}

type text struct{}
type currency struct{}

func (s *text) SetTextRed(str string) string {
	var c consts.Colors
	strRed := c.GetRed() + str + c.GetWhite()
	return strRed
}

func (s *text) SetTextGreen(str string) string {
	var c consts.Colors
	strGreen := c.GetGreen() + str + c.GetWhite()
	return strGreen
}

func (s *text) SetTextYellow(str string) string {
	var c consts.Colors
	strYellow := c.GetYellow() + str + c.GetWhite()
	return strYellow
}

func (s *text) SetTextBlue(str string) string {
	var c consts.Colors
	strBlue := c.GetBlue() + str + c.GetWhite()
	return strBlue
}

func (s *text) SetTextMagenta(str string) string {
	var c consts.Colors
	strMagenta := c.GetMagenta() + str + c.GetWhite()
	return strMagenta
}

func (s *text) SetTextCyan(str string) string {
	var c consts.Colors
	strCyan := c.GetCyan() + str + c.GetWhite()
	return strCyan
}

func (s *text) SetTextWhite(str string) string {
	var c consts.Colors
	strWhite := c.GetWhite() + str
	return strWhite
}

func (s *currency) SetTextRupiah(num int) string {
	return "Rp. " + strconv.Itoa(num)
}
