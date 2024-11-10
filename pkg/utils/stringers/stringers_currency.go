package stringers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type CurrencyCode string

const (
	IDR CurrencyCode = "IDR"
	USD CurrencyCode = "USD"
	SGD CurrencyCode = "SGD"
)

type CurrencyConfig struct {
	Symbol        string
	DecimalPlaces int
	ThousandSep   string
	DecimalSep    string
	Scale         int64
}

var currencyConfigs = map[CurrencyCode]CurrencyConfig{
	IDR: {
		Symbol:        "Rp ",
		DecimalPlaces: 0,
		ThousandSep:   ".",
		DecimalSep:    ",",
		Scale:         1,
	},
	USD: {
		Symbol:        "$",
		DecimalPlaces: 2,
		ThousandSep:   ",",
		DecimalSep:    ".",
		Scale:         100,
	},
	SGD: {
		Symbol:        "S$",
		DecimalPlaces: 2,
		ThousandSep:   ",",
		DecimalSep:    ".",
		Scale:         100,
	},
}

type Currency struct {
	amount int64
	code   CurrencyCode
}

type ErrInvalidCurrency struct {
	Code CurrencyCode
}

func (e ErrInvalidCurrency) Error() string {
	return fmt.Sprintf("invalid currency code: %s", e.Code)
}

func NewCurrency(amount interface{}, code CurrencyCode) (*Currency, error) {
	config, ok := currencyConfigs[code]
	if !ok {
		return nil, &ErrInvalidCurrency{Code: code}
	}

	var internalAmount int64
	var err error

	switch v := amount.(type) {
	case string:
		internalAmount, err = parseAmount(v, config.Scale)
	case int:
		internalAmount = int64(v) * config.Scale
	case int64:
		internalAmount = v * config.Scale
	case float64:
		internalAmount = int64(math.Round(v * float64(config.Scale)))
	default:
		return nil, fmt.Errorf("unsupported amount type: %T", amount)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	return &Currency{amount: internalAmount, code: code}, nil
}

func (c *Currency) Format(code CurrencyCode) (string, error) {
	config, ok := currencyConfigs[code]
	if !ok {
		return "", &ErrInvalidCurrency{Code: code}
	}

	formatted := formatNumber(c.amount, config)
	return config.Symbol + formatted, nil
}

func (c *Currency) String() string {
	config, ok := currencyConfigs[c.code]
	if !ok {
		return fmt.Sprintf("invalid currency code: %s", c.code)
	}

	formatted := formatNumber(c.amount, config)
	return config.Symbol + formatted
}

func (c *Currency) GetAmount(code CurrencyCode) (float64, error) {
	config, ok := currencyConfigs[code]
	if !ok {
		return 0, &ErrInvalidCurrency{Code: code}
	}
	return float64(c.amount) / float64(config.Scale), nil
}

func parseAmount(s string, scale int64) (int64, error) {
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid amount format")
	}

	var intPart, decPart int64
	var err error

	// Parse integer part
	intPart, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer part: %w", err)
	}

	// Handle decimal part if exists
	if len(parts) == 2 {
		// Pad or truncate decimal part to match scale
		decStr := parts[1]
		decLen := int(math.Log10(float64(scale)))
		if len(decStr) > decLen {
			decStr = decStr[:decLen]
		} else {
			decStr = decStr + strings.Repeat("0", decLen-len(decStr))
		}

		decPart, err = strconv.ParseInt(decStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid decimal part: %w", err)
		}
	}

	// Combine parts with proper scaling
	if intPart >= 0 {
		return intPart*scale + decPart, nil
	}
	return intPart*scale - decPart, nil
}

func formatNumber(amount int64, config CurrencyConfig) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}

	intPart := amount / config.Scale
	decPart := amount % config.Scale

	intStr := strconv.FormatInt(intPart, 10)
	var result strings.Builder

	if negative {
		result.WriteString("-")
	}

	for i := 0; i < len(intStr); i++ {
		if i > 0 && (len(intStr)-i)%3 == 0 {
			result.WriteString(config.ThousandSep)
		}
		result.WriteByte(intStr[i])
	}

	if config.DecimalPlaces > 0 {
		result.WriteString(config.DecimalSep)
		decStr := strconv.FormatInt(decPart, 10)
		result.WriteString(strings.Repeat("0", config.DecimalPlaces-len(decStr)))
		result.WriteString(decStr)
	}

	return result.String()
}
