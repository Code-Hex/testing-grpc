package level

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestLog(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want zerolog.Level
	}{
		{
			name: "DEBUG",
			s:    "",
			want: zerolog.DebugLevel,
		},
		{
			name: "INFO",
			s:    "info",
			want: zerolog.InfoLevel,
		},
		{
			name: "WARN",
			s:    "WARN",
			want: zerolog.WarnLevel,
		},
		{
			name: "ERROR",
			s:    "error",
			want: zerolog.ErrorLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Log(tt.s)
			if tt.want != got {
				t.Errorf("Log() = %v, want %v", got, tt.want)
			}
		})
	}
}
