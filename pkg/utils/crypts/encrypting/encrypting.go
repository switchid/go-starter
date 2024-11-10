package encrypting

import (
	"bytes"
	"crypto/aes"
)

type Encrypting interface {
	Encrypt() (string, error)
	Decrypt() (string, error)
}

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func Encryptor(enc Encrypting) (string, error) {
	return enc.Encrypt()
}

func Decryptor(dec Encrypting) (string, error) {
	return dec.Decrypt()
}
