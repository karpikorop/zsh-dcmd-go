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
	}

	for command, want := range cases {
		if got := Variables(command); !reflect.DeepEqual(got, want) {
			t.Errorf("Variables(%q) = %v, want %v", command, got, want)
		}
	}
}

func TestRender(t *testing.T) {
	values := map[string]string{"a": "{{b}}", "b": "B"}

	got := Render("x {{a}} {{ b }} {{missing}} {{a}}", values)
	want := "x {{b}} B {{missing}} {{b}}"
	if got != want {
		t.Errorf("Render = %q, want %q", got, want)
	}
}
