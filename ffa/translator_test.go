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
	if script[0] != "assert(exists 'a');" || script[1] != "assert(exists 'b');" {
		t.Fail()
	}

}

func TestShellPython(t *testing.T) {
	script := getFFAScript(t, "python -d -f script.py arg1 arg2")
	fmt.Println(script[0])
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
	fmt.Println("------------------------------------ SCRIPT START")
	for _, x := range script {
		fmt.Println(x)
	}
	fmt.Println("------------------------------------ SCRIPT END")
	//fmt.Println(script[0])
	//if script[0] != "assert(exists 'script.py');" {
	//	t.Fail()
	//}
}
