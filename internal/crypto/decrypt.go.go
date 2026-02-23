package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func Decrypt(priv *rsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.New()
	return rsa.DecryptOAEP(hash, rand.Reader, priv, data, nil)
}
