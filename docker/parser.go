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
	"os"
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

func AnalyzeRunCommand(cmd dockerfile.Command) {
	commandString := strings.Join(cmd.Value, " ")
	in := strings.NewReader(commandString)
	f, err := syntax.NewParser().Parse(in, "")
	if err != nil {
		return
	}
	fmt.Println("\tRun command:", commandString)
	p := syntax.NewPrinter()
	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.CallExpr:
			p.Print(os.Stdout, x.Args[0])
			fmt.Println()
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
			//syntax.NewPrinter().Print(os.Stdout, x.Op)
			fmt.Println(x.Op.String())
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
			p.Print(os.Stdout, x.Name)
			break
		default:
		}
		return true
	})
	//syntax.NewPrinter().Print(os.Stdout, f)
}
