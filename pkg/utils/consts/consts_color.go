package consts

const (
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[0m"
)

type Colors struct{}

func (c *Colors) GetRed() string {
	return red
}

func (c *Colors) GetGreen() string {
	return green
}

func (c *Colors) GetYellow() string {
	return yellow
}

func (c *Colors) GetBlue() string {
	return blue
}

func (c *Colors) GetMagenta() string {
	return magenta
}

func (c *Colors) GetCyan() string {
	return cyan
}

func (c *Colors) GetWhite() string {
	return white
}
