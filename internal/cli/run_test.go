package cli

import (
	"reflect"
	"strings"
	"testing"
)

func TestResolveValuesWithoutTerminal(t *testing.T) {
	command := "deploy {{name}} --project {{project}}"

	values, err := resolveValues(command, map[string]string{"name": "fn", "project": "dev"})
	if err != nil || !reflect.DeepEqual(values, map[string]string{"name": "fn", "project": "dev"}) {
		t.Errorf("all values given: got %v, %v", values, err)
	}

	_, err = resolveValues(command, map[string]string{"name": "fn"})
	if err == nil || !strings.Contains(err.Error(), "'project'") {
		t.Errorf("missing value: got %v, want an error that names 'project'", err)
	}

	_, err = resolveValues(command, map[string]string{"name": "fn", "project": "dev", "typo": "x"})
	if err == nil || !strings.Contains(err.Error(), "'typo'") {
		t.Errorf("unknown value: got %v, want an error that names 'typo'", err)
	}
}

func TestParseVarFlags(t *testing.T) {
	values, err := parseVarFlags([]string{"File Name=a=b.ts", "empty="})
	want := map[string]string{"File Name": "a=b.ts", "empty": ""}
	if err != nil || !reflect.DeepEqual(values, want) {
		t.Errorf("got %v, %v; want %v", values, err, want)
	}

	for _, bad := range []string{"novalue", "=x"} {
		if _, err := parseVarFlags([]string{bad}); err == nil {
			t.Errorf("parseVarFlags(%q): expected an error", bad)
		}
	}
}
