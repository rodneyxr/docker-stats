package ffa

import (
	"fmt"
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
	if script[0] != "assert 'a';" {
		t.Fail()
	}

	script = getFFAScript(t, "chmod u+x -R a b")
	fmt.Println(script)
	if script[0] != "assert 'a';" || script[1] != "assert 'b';" {
		t.Fail()
	}

}

func TestShellPython(t *testing.T) {
	script := getFFAScript(t, "python -d -f script.py arg1 arg2")
	if script[0] != "assert 'script.py';" {
		t.Fail()
	}
}
