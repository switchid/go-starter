package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type SHA256 struct {
	Text   string
	Salt   string
	Hashes string
}

func NewSHA256(str string) *SHA256 {
	return &SHA256{Text: str}
}

func NewSHA256WithSalt(str string, salt string) *SHA256 {
	return &SHA256{Text: str, Salt: salt}
}

func VerifySHA256(str string, hash string) *SHA256 {
	return &SHA256{Text: str, Hashes: hash}
}

func VerifySHA256WithSalt(str string, salt string, hash string) *SHA256 {
	return &SHA256{Text: str, Salt: salt, Hashes: hash}
}

func (hash SHA256) Hash() (string, error) {
	var plaintext string
	text := hash.Text
	salt := hash.Salt
	if salt != "" {
		plaintext = salt + text
	} else {
		plaintext = text
	}
	sh := sha256.New()
	sh.Write([]byte(plaintext))
	return hex.EncodeToString(sh.Sum(nil)), nil
}

func (hash SHA256) Verify() (bool, error) {
	var plaintext string
	text := hash.Text
	salt := hash.Salt
	if salt != "" {
		plaintext = salt + text
	} else {
		plaintext = text
	}
	sh := sha256.New()
	sh.Write([]byte(plaintext))
	result := hex.EncodeToString(sh.Sum(nil))
	if result != hash.Hashes {
		return false, errors.New("Hashes does not match")
	}
	return true, nil
}

func (hash SHA256) MakeHash() (string, error) {
	return MakeHash(hash)
}

func (hash SHA256) VerifyHash() (bool, error) {
	return VerifyHash(hash)
}
