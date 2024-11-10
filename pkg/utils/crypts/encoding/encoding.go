package encoding

type Encoding interface {
	Encode() (string, error)
	Decode() (string, error)
}

func EncodingText(enc Encoding) (string, error) {
	return enc.Encode()
}

func DecodingText(enc Encoding) (string, error) {
	return enc.Decode()
}
