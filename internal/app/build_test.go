package app

import "testing"

func TestValueOrNA(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "N/A",
		},
		{
			name:  "non empty string",
			input: "1.0.0",
			want:  "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := valueOrNA(tt.input)
			if got != tt.want {
				t.Errorf("valueOrNA(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
