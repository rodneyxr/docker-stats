package ffa

import (
	"mvdan.cc/sh/v3/syntax"
	"strings"
)

// literize converts a []*syntax.Word to a []string
func literize(arguments []*syntax.Word) []string {
	var args []string
	for _, arg := range arguments {
		args = append(args, arg.Lit())
	}
	return args
}

// extractFlag strips the flag and nFlags number of flags after the flag from
// the command string provided.
// Returns the stripped command string and extracted flags (flag inclusive)
func extractFlag(command []string, flag string, nFlags int) ([]string, []string) {
	// Search for the flag in the command
	marker := -1
	for i, s := range command {
		if flag == s {
			marker = i
		}
	}

	// Return if the flag was not found
	if marker == -1 {
		return command, nil
	}

	// Remove the specified flag and arguments from the command
	var newCommand []string
	var extractedFlags []string
	for i, s := range command {
		if i >= marker && i <= marker+nFlags {
			extractedFlags = append(extractedFlags, s)
		} else {
			newCommand = append(newCommand, s)
		}
	}
	return newCommand, extractedFlags
}

// removeFlags removes flags from a []*syntax.Word
func removeFlags(arguments []*syntax.Word) []string {
	return removeFlagsLit(literize(arguments))
}

// removeFlagsLit removes flags from a literized array
func removeFlagsLit(arguments []string) []string {
	var args []string
	for _, arg := range arguments {
		if !strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		}
	}
	return args
}
