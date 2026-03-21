package ui

import (
	"reflect"
	"testing"
)

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		template string
		want     []string
	}{
		{"echo {{name}}", []string{"{{name}}"}},
		{"echo {{name}} {{greeting}}", []string{"{{name}}", "{{greeting}}"}},
		{"echo {{name}} {{name}}", []string{"{{name}}"}}, // deduplication
		{"echo no vars", nil},
	}

	for _, tt := range tests {
		got := ExtractVariables(tt.template)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ExtractVariables(%q) = %v, want %v", tt.template, got, tt.want)
		}
	}
}
