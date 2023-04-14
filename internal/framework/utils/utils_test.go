package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

func TestHasUnicodeChars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "empty string",
			s:    "",
			want: false,
		},
		{
			name: "ascii string",
			s:    "hello-world",
			want: false,
		},
		{
			name: "unicode string",
			s:    "hello-世界",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, utils.HasUnicodeChars(tt.s))
		})
	}
}
