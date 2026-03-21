package ui

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
)

func ExtractVariables(template string) []string {
	re := regexp.MustCompile(`{{([^}]+)}}`)
	matches := re.FindAllStringSubmatch(template, -1)
	
	seen := make(map[string]bool)
	var variables []string
	for _, match := range matches {
		if !seen[match[0]] {
			seen[match[0]] = true
			variables = append(variables, match[0])
		}
	}
	return variables
}

func PromptForVariables(template string) (string, error) {
	variables := ExtractVariables(template)
	if len(variables) == 0 {
		return template, nil
	}

	finalCmd := template
	for _, varTag := range variables {
		cleanName := varTag[2 : len(varTag)-2]
		
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("Enter value for '%s'", cleanName),
			// Use Stderr for the prompt itself so it doesn't pollute stdout
			Stdin:  os.Stdin,
			Stdout: os.Stderr,
		}

		result, err := prompt.Run()
		if err != nil {
			return "", err
		}

		finalCmd = strings.ReplaceAll(finalCmd, varTag, result)
	}

	return finalCmd, nil
}

func PromptConfirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label + " (Y/n)",
		IsConfirm: true,
		Stdin:     os.Stdin,
		Stdout:    os.Stderr,
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	return result == "y" || result == "Y" || result == "", nil
}

func PromptInput(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:  label,
		Stdin:  os.Stdin,
		Stdout: os.Stderr,
	}
	return prompt.Run()
}

func PromptInputWithDefault(label, defaultValue string) (string, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: fmt.Sprintf("%s: ", label),
		Stdout: os.Stderr,
		Stdin:  os.Stdin,
	})
	if err != nil {
		return "", err
	}
	defer rl.Close()

	rl.WriteStdin([]byte(defaultValue))

	result, err := rl.Readline()
	if err != nil {
		return "", err
	}
	return result, nil
}

func PromptInputNoPipe(label string) (string, error) {
	for {
		result, err := PromptInput(label)
		if err != nil {
			return "", err
		}
		if strings.Contains(result, "|") {
			fmt.Fprintf(os.Stderr, "Error: input cannot contain the '|' character. Please try again.\n")
			continue
		}
		return result, nil
	}
}

func PromptInputWithDefaultNoPipe(label, defaultValue string) (string, error) {
	for {
		result, err := PromptInputWithDefault(label, defaultValue)
		if err != nil {
			return "", err
		}
		if strings.Contains(result, "|") {
			fmt.Fprintf(os.Stderr, "Error: input cannot contain the '|' character. Please try again.\n")
			continue
		}
		return result, nil
	}
}
