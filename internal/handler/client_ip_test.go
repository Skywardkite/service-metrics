package handler

import (
	"net/http"
	"testing"
)

func TestClientIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		want       string
	}{
		{
			name: "x-forwarded-for single ip",
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.1",
			},
			remoteAddr: "10.0.0.1:1234",
			want:       "192.168.1.1",
		},
		{
			name: "x-forwarded-for multiple ips",
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.1, 10.0.0.1",
			},
			remoteAddr: "10.0.0.2:1234",
			want:       "192.168.1.1",
		},
		{
			name:       "remote addr with port",
			remoteAddr: "127.0.0.1:8080",
			want:       "127.0.0.1",
		},
		{
			name:       "remote addr without port",
			remoteAddr: "127.0.0.1",
			want:       "127.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{
				Header:     make(http.Header),
				RemoteAddr: tt.remoteAddr,
			}

			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}

			got := clientIP(r)
			if got != tt.want {
				t.Errorf("clientIP() = %q, want %q", got, tt.want)
			}
		})
	}
}
