// Package shell holds the shell code that connects dcmd to the user's prompt.
package shell

import (
	_ "embed"
	"strings"
)

//go:embed zsh.zsh
var zshScript string

const binPlaceholder = "__DCMD_BIN__"

// Zsh returns the zsh function that calls the dcmd binary at binPath.
func Zsh(binPath string) string {
	return strings.ReplaceAll(zshScript, binPlaceholder, quote(binPath))
}

// quote wraps s in single quotes for zsh, so a path with spaces or quotes stays one word.
func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
