package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func signBody(body []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)

	return hex.EncodeToString(h.Sum(nil))
}

func checkKey(req *http.Request, key string) (bool, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))

	expected := signBody(body, key)
	received := req.Header.Get("HashSHA256")

	if received == "" || !hmac.Equal([]byte(received), []byte(expected)) {
		return false, nil
	}

	return true, nil
}