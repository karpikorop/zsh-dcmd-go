package cli

import (
	"reflect"
	"strings"
	"testing"
)

func TestResolveValuesWithoutTerminal(t *testing.T) {
	command := "deploy {{name}} --project {{project}} --region {{region:us-central1}}"
	projectVars := map[string]string{"project": "app-dev", "unused": "x"}

	values, err := resolveValues(command, projectVars, map[string]string{"name": "fn"})
	want := map[string]string{"name": "fn", "project": "app-dev", "region": "us-central1"}
	if err != nil || !reflect.DeepEqual(values, want) {
		t.Errorf("project value and default: got %v, %v; want %v", values, err, want)
	}

	values, err = resolveValues(command, projectVars, map[string]string{"name": "fn", "project": "prod", "region": "eu"})
	want = map[string]string{"name": "fn", "project": "prod", "region": "eu"}
	if err != nil || !reflect.DeepEqual(values, want) {
		t.Errorf("given values win: got %v, %v; want %v", values, err, want)
	}

	_, err = resolveValues(command, projectVars, nil)
	if err == nil || !strings.Contains(err.Error(), "'name'") || strings.Contains(err.Error(), "'region'") {
		t.Errorf("missing value: got %v, want an error that names only 'name'", err)
	}

	_, err = resolveValues(command, projectVars, map[string]string{"name": "fn", "typo": "x"})
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
