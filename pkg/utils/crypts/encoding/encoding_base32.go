package encoding

import "encoding/base32"

const base32enc = "0123456789aBcDeFgHiJkLmNoPqRsTuV"

type BASE32 struct {
	Text string
}

func NewEncodingBASE32(text string) *BASE32 {
	return &BASE32{Text: text}
}

func (enc *BASE32) Encode() (string, error) {
	text := []byte(enc.Text)
	encoding := base32.NewEncoding(base32enc).WithPadding(base32.StdPadding)
	return encoding.EncodeToString(text), nil
}

func (enc *BASE32) Decode() (string, error) {
	encoding := base32.NewEncoding(base32enc).WithPadding(base32.StdPadding)
	decoder, errDecoder := encoding.DecodeString(enc.Text)
	if errDecoder != nil {
		return "", errDecoder
	}
	return string(decoder), nil
}

func (enc *BASE32) EncodingText() (string, error) {
	return EncodingText(enc)
}

func (enc *BASE32) DecodingText() (string, error) {
	return DecodingText(enc)
}
