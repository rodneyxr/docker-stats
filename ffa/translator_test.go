package ffa

import (
	"strings"
	"testing"
)

func getFFAScript(t *testing.T, sh string) []string {
	script, err := TranslateShellScript(sh)
	if err != nil {
		t.Fatal(err)
	}
	return script
}

func TestShellChmod(t *testing.T) {
	// Parse the Dockerfile
	script := getFFAScript(t, "chmod u+x a")
	tokens := []string{"assert(exists 'a');"}
	tokenCount := verifyTokens(tokens, script)
	if tokenCount != len(tokens) {
		t.Errorf("token '%s' not found", tokens[tokenCount])
	}

	script = getFFAScript(t, "chmod u+x -R a b")
	tokens = []string{"assert(exists 'a');", "assert(exists 'b');"}
	tokenCount = verifyTokens(tokens, script)
	if tokenCount != len(tokens) {
		t.Errorf("token '%s' not found", tokens[tokenCount])
	}
}

func TestShellPython(t *testing.T) {
	script := getFFAScript(t, "python -d -f script.py arg1 arg2")
	if script[0] != "assert(exists 'script.py');" {
		t.Fail()
	}
}

func TestShellIfStatements(t *testing.T) {
	script := getFFAScript(t, `
#!/bin/bash
touch before
if [ $1 -gt 100 ]
then
	touch inside
	pwd
fi
touch after
`)
	tokens := []string{"before", "if", "inside", "after"}
	tokenCount := verifyTokens(tokens, script)
	if tokenCount != len(tokens) {
		t.Errorf("token '%s' not found", tokens[tokenCount])
	}
}

func TestShellWhileStatements(t *testing.T) {
	script := getFFAScript(t, `
#!/bin/bash
touch before
i=0

while [ $i -le 2 ]
do
  touch inside
  ((i++))
done
touch after
`)
	//fmt.Println(strings.Join(script, "\n"))

	tokens := []string{"before", "while", "inside", "after"}
	tokenCount := verifyTokens(tokens, script)
	if tokenCount != len(tokens) {
		t.Errorf("token '%s' not found", tokens[tokenCount])
	}
}

func TestShellNestedWhileStatements(t *testing.T) {
	script := getFFAScript(t, `
#!/bin/bash
touch before
i=0
j=0

while [ $i -le 2 ]
do
	touch inside1
	((i++))
	while [ $j -le 2 ]
	do
	  touch inside2
	  ((j++))
	done
done
touch after
`)
	//fmt.Println(strings.Join(script, "\n"))

	tokens := []string{"before", "while", "inside1", "while", "inside2", "}", "after"}
	tokenCount := verifyTokens(tokens, script)
	if tokenCount != len(tokens) {
		t.Errorf("token '%s' not found", tokens[tokenCount])
	}
}

func TestShellIfElseStatements(t *testing.T) {
	script := getFFAScript(t, `
#!/bin/bash
touch before

if [ $ANSI_INSTALLED -eq $SUCC ]; then
	touch inside1
else
	touch inside2
fi

touch after
`)
	//fmt.Println(strings.Join(script, "\n"))

	tokens := []string{"before", "if", "inside1", "else", "inside2", "}", "after"}
	tokenCount := verifyTokens(tokens, script)
	if tokenCount != len(tokens) {
		t.Errorf("token '%s' not found", tokens[tokenCount])
	}
}

func TestShellElifStatements(t *testing.T) {
	script := getFFAScript(t, `
#!/bin/bash
touch before

if [ $ANSI_INSTALLED -eq $SUCC ]; then
	touch inside1
elif [ $ANSI_INSTALLED ]; then
	touch inside2
elif [ $ANSI_INSTALLED ]; then
	touch inside3
else
    touch inside4
fi

touch after
`)
	//fmt.Println(strings.Join(script, "\n"))

	tokens := []string{"before", "if", "inside1", "else if", "inside2", "else if", "inside3", "after"}
	tokenCount := verifyTokens(tokens, script)
	if tokenCount != len(tokens) {
		t.Errorf("token '%s' not found", tokens[tokenCount])
	}
}

// verifyTokens ensures that all tokens are found in the script provided.
// Returns the number of tokens found. If the return value equals the length of
// the tokens array, then all tokens were found.
func verifyTokens(tokens, script []string) int {
	// Verify the token list in order
	tokenCount := 0
outer:
	for _, x := range script {
		for strings.Contains(x, tokens[tokenCount]) {
			tokenCount++
			if tokenCount == len(tokens) {
				break outer
			}
		}
	}
	return tokenCount
}
