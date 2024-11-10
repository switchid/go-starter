package hashing

type Hashing interface {
	Hash() (string, error)
	Verify() (bool, error)
}

func MakeHash(hash Hashing) (string, error) {
	return hash.Hash()
}

func VerifyHash(hash Hashing) (bool, error) {
	return hash.Verify()
}
