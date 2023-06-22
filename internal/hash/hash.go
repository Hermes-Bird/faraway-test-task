package hash

import (
	"crypto/sha256"
)

func HashBytes(bs []byte) []byte {
	hashBytes := sha256.Sum256(bs)
	return hashBytes[:]
}
