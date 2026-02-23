package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func Encrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	hash := sha256.New()
	return rsa.EncryptOAEP(hash, rand.Reader, pub, data, nil)
}
