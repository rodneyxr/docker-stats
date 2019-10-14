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
	"path/filepath"
	"strings"

	"github.com/asottile/dockerfile"
	"github.com/mvdan/sh/syntax"
)

// ExtractAllCommandsFromDockerfile
func ExtractAllCommandsFromDockerfile(data string) ([]dockerfile.Command, error) {
	reader := strings.NewReader(data)
	commandList, err := dockerfile.ParseReader(reader)
	if err != nil {
		return nil, err
	}
	return commandList, nil
}

// ExtractRunCommandsFromDockerfile
func ExtractRunCommandsFromDockerfile(data string) ([]dockerfile.Command, error) {
	commandList, err := ExtractAllCommandsFromDockerfile(data)
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

func AnalyzeRunCommand(cmd dockerfile.Command) []string {
	var ffaList []string
	commandString := strings.Join(cmd.Value, " ")
	in := strings.NewReader(commandString)
	f, err := syntax.NewParser().Parse(in, "")
	if err != nil {
		return ffaList
	}
	//fmt.Println("\tRun command:", commandString)
	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.CallExpr:
			// only handle most common commands
			// go through all projects and rank most common commands
			cmd := x.Args[0].Lit()
			switch cmd {
			case "touch":
				// Create a touch statement for each argument
				for _, s := range x.Args[1:] {
					ffaList = append(ffaList, fmt.Sprintf("touch '%s';", s.Lit()))
				}
				break
			case "mkdir":
				args := removeFlags(x.Args)
				for _, s := range args[1:] {
					ffaList = append(ffaList, fmt.Sprintf("mkdir '%s';", s))
				}
				break
			case "rm":
			case "rmdir":
				// TODO: check for flags
				for _, s := range x.Args[1:] {
					ffaList = append(ffaList, fmt.Sprintf("rm '%s';", s.Lit()))
				}
				break
			case "cp":
				args := removeFlags(x.Args)
				arg1, arg2 := args[1], args[2]
				ffaList = append(ffaList, fmt.Sprintf("cp '%s' '%s';", arg1, arg2))
				break
			case "mv":
				args := removeFlags(x.Args)
				arg1, arg2 := args[1], args[2]
				ffaList = append(ffaList, fmt.Sprintf("cp '%s' '%s';", arg1, arg2))
				ffaList = append(ffaList, fmt.Sprintf("rm '%s';", arg1))
				break
			case "git":
				// TODO: handle git rm
				arg1 := x.Args[1].Lit()
				if arg1 == "clone" {
					dirname := filepath.Base(x.Args[2].Lit())
					ffaList = append(ffaList, fmt.Sprintf("mkdir '%s';", dirname))
				}
				break
			case "cd":
				arg1 := x.Args[1].Lit()
				ffaList = append(ffaList, fmt.Sprintf("cd '%s';", arg1))
				break
			case "wget":
				// TODO: handle wget
				arg1 := x.Args[1].Lit()
				if !strings.HasPrefix("-", arg1) {
					ffaList = append(ffaList, fmt.Sprintf("touch '%s';", filepath.Base(arg1)))
				}
				// TODO: handle -O parameter
				break
			case "curl":
				index := -1
				for i, word := range x.Args {
					arg := word.Lit()
					if strings.HasPrefix("-O", arg) {
						index = i + 1
					}
				}
				if index < len(x.Args) && index >= 0 {
					arg := x.Args[index].Lit()
					ffaList = append(ffaList, fmt.Sprintf("touch '%s';", arg))
				}
				break
			case "tar":
				// TODO: handle tar
				break
			case "set":
				// TODO: handle variables
				break
			case "ln":
				// TODO: handle symlinks
				break
			case "export":
				// TODO: handle variables
				break
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
			ffaList = append(ffaList, fmt.Sprintf("$x? = '%s';", x.Name.Value))
			break
		default:
		}
		return true
	})
	return ffaList
}

func removeFlags(arguments []*syntax.Word) []string {
	var args []string
	for _, arg := range arguments {
		if !strings.HasPrefix(arg.Lit(), "-") {
			args = append(args, arg.Lit())
		}
	}
	return args
}
