package ui

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
)

var ErrCancelled = errors.New("cancelled")

var variablePattern = regexp.MustCompile(`{{([^}]+)}}`)

func ExtractVariables(template string) []string {
	seen := make(map[string]bool)
	var variables []string

	for _, match := range variablePattern.FindAllString(template, -1) {
		if !seen[match] {
			seen[match] = true
			variables = append(variables, match)
		}
	}
	return variables
}

func PromptForVariables(template string) (string, error) {
	finalCmd := template

	for _, tag := range ExtractVariables(template) {
		name := tag[2 : len(tag)-2]
		value, err := Input(fmt.Sprintf("Enter value for '%s'", name))
		if err != nil {
			return "", err
		}
		finalCmd = strings.ReplaceAll(finalCmd, tag, value)
	}
	return finalCmd, nil
}

func Input(label string) (string, error) {
	prompt := promptui.Prompt{Label: label, Stdin: os.Stdin, Stdout: os.Stderr}
	result, err := prompt.Run()
	return result, normalizeError(err)
}

// InputWithDefault shows an editable prompt that starts with defaultValue.
func InputWithDefault(label, defaultValue string) (string, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: label + ": ",
		Stdin:  os.Stdin,
		Stdout: os.Stderr,
	})
	if err != nil {
		return "", err
	}
	defer rl.Close()

	rl.WriteStdin([]byte(defaultValue))
	result, err := rl.Readline()
	return result, normalizeError(err)
}

// Confirm asks a yes/no question. An empty answer means yes.
func Confirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Default:   "y",
		Stdin:     os.Stdin,
		Stdout:    os.Stderr,
	}

	_, err := prompt.Run()
	if errors.Is(err, promptui.ErrAbort) {
		return false, nil
	}
	return err == nil, normalizeError(err)
}

func normalizeError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, promptui.ErrInterrupt),
		errors.Is(err, promptui.ErrEOF),
		errors.Is(err, readline.ErrInterrupt),
		errors.Is(err, io.EOF):
		return ErrCancelled
	default:
		return err
	}
}
