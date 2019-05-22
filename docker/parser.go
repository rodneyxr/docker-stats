// Copyright © 2019 Rodney Rodriguez
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docker

import (
	"fmt"
	"strings"

	"github.com/asottile/dockerfile"
	"github.com/mvdan/sh/syntax"
)

// ExtractRunCommandsFromDockerfile
func ExtractRunCommandsFromDockerfile(data string) ([]dockerfile.Command, error) {
	reader := strings.NewReader(data)
	commandList, err := dockerfile.ParseReader(reader)
	if err != nil {
		return nil, err
	}

	var commands []dockerfile.Command

	// Print all commands in the Dockerfile
	for _, cmd := range commandList {
		if cmd.Cmd == "run" {
			commands = append(commands, cmd)
			//fmt.Println(cmd.Cmd, cmd.Value)
		}
	}
	return commands, nil
}

func extractFlags(word *syntax.Word) []string {
	var flags []string

	flags = append(flags, "")

	return flags
}

func AnalyzeRunCommand(cmd dockerfile.Command) {
	commandString := strings.Join(cmd.Value, " ")
	in := strings.NewReader(commandString)
	f, err := syntax.NewParser().Parse(in, "")
	if err != nil {
		return
	}
	fmt.Println("\tRun command:", commandString)
	//p := syntax.NewPrinter()
	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.CallExpr:
			// only handle most common commands
			// go through all projects and rank most common commands
			cmd := x.Args[0].Lit()
			if cmd == "touch" {
				// Create a touch statement for each argument
				for _, s := range x.Args[1:] {
					fmt.Println("touch", s.Lit())
				}
			} else if cmd == "mkdir" {
				for _, s := range x.Args[1:] {
					fmt.Println("mkdir", s.Lit())
				}
			} else if cmd == "rm" || cmd == "rmdir" {
				// TODO: check for flags
				for _, s := range x.Args[1:] {
					fmt.Println("rm", s.Lit())
				}
			} else if cmd == "cp" {
				arg1, arg2 := x.Args[1].Lit(), x.Args[2].Lit()
				fmt.Println("cp", arg1, arg2)
			} else if cmd == "mv" {
				arg1, arg2 := x.Args[1].Lit(), x.Args[2].Lit()
				fmt.Println("cp", arg1, arg2)
				fmt.Println("rm", arg1)
			}
			break
		case *syntax.IfClause:
			break
		case *syntax.WhileClause:
			break
		case *syntax.ForClause:
			break
		case *syntax.CaseClause:
			break
		case *syntax.Block:
			break
		case *syntax.Subshell:
			break
		case *syntax.BinaryCmd:
			break
		case *syntax.FuncDecl:
			break
		case *syntax.ArithmCmd:
			break
		case *syntax.TestClause:
			break
		case *syntax.DeclClause:
			break
		case *syntax.LetClause:
			break
		case *syntax.TimeClause:
			break
		case *syntax.CoprocClause:
			break
		case *syntax.Assign:
			fmt.Print("assign: ")
			fmt.Println("$x? =", x.Name.Value)
			break
		default:
		}
		return true
	})
	//syntax.NewPrinter().Print(os.Stdout, f)
}
