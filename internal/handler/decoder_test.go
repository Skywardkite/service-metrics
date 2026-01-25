package handler

import (
	"testing"
)

func TestSignBody(t *testing.T) {
	tests := []struct {
		name string
		body []byte
		key  string
		want string
	}{
		{
			name: "success",
			body: []byte("hello"),
			key:  "secret",
			want: "88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b",
		},
		{
			name: "empty body",
			body: []byte(""),
			key:  "secret",
			want: "f9e66e179b6747ae54108f82f8ade8b3c25d76fd30afde6c395822c530196169",
		},
		{
			name: "different key",
			body: []byte("hello"),
			key:  "other",
			want: "ccba5a99c3e0c7f2bc6c6f02bc752217c3c71fb91b41211dffc0c1d0981d1663",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SignBody(tt.body, tt.key)
			if got != tt.want {
				t.Errorf("SignBody() = %q, want %q", got, tt.want)
			}
		})
	}
}
