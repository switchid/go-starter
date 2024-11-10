package encoding

import "encoding/base64"

const base64enc = "ABCDEFGHIJKLMNOPQRSTUVabcdefghijklmnopqrstuv1234567890!#"

type BASE64 struct {
	Text string
}

func NewEncodingBASE64(text string) *BASE64 {
	return &BASE64{Text: text}
}

func (enc *BASE64) Encode() (string, error) {
	text := []byte(enc.Text)
	encoding := base64.NewEncoding(base64enc)
	return encoding.EncodeToString(text), nil
}

func (enc *BASE64) Decode() (string, error) {
	encoding := base64.NewEncoding(base64enc).WithPadding(base64.StdPadding)
	decoder, errDecoder := encoding.DecodeString(enc.Text)
	if errDecoder != nil {
		return "", errDecoder
	}
	return string(decoder), nil
}

func (enc *BASE64) EncodingText() (string, error) {
	return EncodingText(enc)
}

func (enc *BASE64) DecodingText() (string, error) {
	return DecodingText(enc)
}
