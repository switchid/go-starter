package hashing

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
)

type SHA512 struct {
	Text   string
	Salt   string
	Hashes string
}

func NewSHA512(str string) *SHA512 {
	return &SHA512{Text: str}
}

func NewSHA512WithSalt(str string, salt string) *SHA512 {
	return &SHA512{Text: str, Salt: salt}
}

func VerifySHA512(str string, hash string) *SHA512 {
	return &SHA512{Text: str, Hashes: hash}
}

func VerifySHA512WithSalt(str string, salt string, hash string) *SHA512 {
	return &SHA512{Text: str, Salt: salt, Hashes: hash}
}

func (hash SHA512) Hash() (string, error) {
	var plaintext string
	text := hash.Text
	salt := hash.Salt
	if salt != "" {
		plaintext = salt + text
	} else {
		plaintext = text
	}
	sh := sha512.New()
	sh.Write([]byte(plaintext))
	return hex.EncodeToString(sh.Sum(nil)), nil
}

func (hash SHA512) Verify() (bool, error) {
	var plaintext string
	text := hash.Text
	salt := hash.Salt
	if salt != "" {
		plaintext = salt + text
	} else {
		plaintext = text
	}
	sh := sha512.New()
	sh.Write([]byte(plaintext))
	result := hex.EncodeToString(sh.Sum(nil))
	if result != hash.Hashes {
		return false, errors.New("Hashes does not match")
	}
	return true, nil
}

func (hash SHA512) MakeHash() (string, error) {
	return MakeHash(hash)
}

func (hash SHA512) VerifyHash() (bool, error) {
	return VerifyHash(hash)
}
