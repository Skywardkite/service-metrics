package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// SignBody подписываем response с помощью секретного ключа.
func SignBody(body []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)

	return hex.EncodeToString(h.Sum(nil))
}
