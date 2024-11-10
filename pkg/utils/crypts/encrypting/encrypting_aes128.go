package encrypting

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

type AES128 struct {
	Text string
	Key  string
}

func NewAES128(text string, key string) *AES128 {
	return &AES128{Text: text, Key: key}
}

func (enc AES128) Encrypt() (string, error) {
	key := []byte(enc.Key)
	if len(key) < 16 {
		return "", errors.New("key too short")
	} else if len(key) > 16 {
		return "", errors.New("key too long")
	}
	block, errBlock := aes.NewCipher(key)
	if errBlock != nil {
		return "", errBlock
	}

	plaintext := pad([]byte(enc.Text))
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	// Generate a random IV.
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

func (enc AES128) Decrypt() (string, error) {
	key := []byte(enc.Key)
	if len(key) < 16 {
		return "", errors.New("key too short")
	} else if len(key) > 16 {
		return "", errors.New("key too long")
	}
	block, errBlock := aes.NewCipher(key)
	if errBlock != nil {
		return "", errBlock
	}
	ciphertext, errCiphertext := hex.DecodeString(enc.Text)
	if errCiphertext != nil {
		return "", errCiphertext
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return string(unpad(ciphertext)), nil
}

func (enc AES128) Encryptor() (string, error) {
	return Encryptor(enc)
}

func (enc AES128) Dercryptor() (string, error) {
	return Decryptor(enc)
}
