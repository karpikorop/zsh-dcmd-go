// Package template finds and fills {{name}} placeholders in a command.
package template

import (
	"regexp"
	"strings"
)

var placeholderPattern = regexp.MustCompile(`{{([^}]+)}}`)

type Variable struct {
	Name string
}

// Variables returns each variable of a command once, in the order of first use.
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

func parsePlaceholder(inner string) Variable {
	return Variable{Name: strings.TrimSpace(inner)}
}
