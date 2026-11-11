package ui

import (
	"errors"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
)

var (
	ErrCancelled  = errors.New("cancelled")
	ErrNoTerminal = errors.New("this needs an interactive terminal. Pass the values as flags instead (see 'dcmd help')")
)

// IsInteractive reports whether stdin is a terminal, so prompts can read an answer.
func IsInteractive() bool {
	return readline.IsTerminal(int(os.Stdin.Fd()))
}

func Input(label string) (string, error) {
	if !IsInteractive() {
		return "", ErrNoTerminal
	}

	prompt := promptui.Prompt{Label: label, Stdin: os.Stdin, Stdout: os.Stderr}
	result, err := prompt.Run()
	return result, normalizeError(err)
}

// InputWithDefault shows an editable prompt that starts with defaultValue.
func InputWithDefault(label, defaultValue string) (string, error) {
	if !IsInteractive() {
		return "", ErrNoTerminal
	}

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
	if !IsInteractive() {
		return false, ErrNoTerminal
	}

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
