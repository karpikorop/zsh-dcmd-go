package template

import (
	"reflect"
	"testing"
)

func TestVariables(t *testing.T) {
	cases := map[string][]Variable{
		"echo no vars":                        {},
		"echo {{name}} {{ greeting }}":        {{Name: "name"}, {Name: "greeting"}},
		"echo {{name}} and {{name}} again":    {{Name: "name"}},
		"yarn prettier --write {{File Name}}": {{Name: "File Name"}},
		"deploy {{fn: cloud:invoca }} {{fn:other}} {{empty:}}": {
			{Name: "fn", Default: "cloud:invoca", HasDefault: true},
			{Name: "empty", Default: "", HasDefault: true},
		},
	}

	for command, want := range cases {
		if got := Variables(command); !reflect.DeepEqual(got, want) {
			t.Errorf("Variables(%q) = %v, want %v", command, got, want)
		}
	}
}

func TestRender(t *testing.T) {
	values := map[string]string{"a": "{{b}}", "b": "B"}

	got := Render("x {{a}} {{ b:default }} {{missing}} {{a}}", values)
	want := "x {{b}} B {{missing}} {{b}}"
	if got != want {
		t.Errorf("Render = %q, want %q", got, want)
	}
}
