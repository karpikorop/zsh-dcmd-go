// Package template finds and fills {{name}} and {{name:default}} placeholders in a command.
package template

import (
	"regexp"
	"strings"
)

var placeholderPattern = regexp.MustCompile(`{{([^}]+)}}`)

type Variable struct {
	Name       string
	Default    string
	HasDefault bool
}

// Variables returns each variable of a command once, in the order of first use.
// If a variable appears more than once, its first placeholder sets the default.
func Variables(command string) []Variable {
	seen := make(map[string]bool)
	variables := []Variable{}

	for _, match := range placeholderPattern.FindAllStringSubmatch(command, -1) {
		v := parsePlaceholder(match[1])
		if !seen[v.Name] {
			seen[v.Name] = true
			variables = append(variables, v)
		}
	}
	return variables
}

// Render replaces each placeholder with its value in one pass. A placeholder with no value stays as it is.
func Render(command string, values map[string]string) string {
	return placeholderPattern.ReplaceAllStringFunc(command, func(placeholder string) string {
		v := parsePlaceholder(placeholder[2 : len(placeholder)-2])
		if value, ok := values[v.Name]; ok {
			return value
		}
		return placeholder
	})
}

// ValidName reports whether name can be used inside a placeholder.
func ValidName(name string) bool {
	return strings.TrimSpace(name) != "" && !strings.ContainsAny(name, ":{}")
}

// parsePlaceholder splits "name:default" on the first colon, so a default can contain colons.
func parsePlaceholder(inner string) Variable {
	name, defaultValue, hasDefault := strings.Cut(inner, ":")
	return Variable{
		Name:       strings.TrimSpace(name),
		Default:    strings.TrimSpace(defaultValue),
		HasDefault: hasDefault,
	}
}
