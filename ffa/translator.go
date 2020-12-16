package ffa

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mvdan/sh/syntax"
)

func TranslateDockerfile(data string) ([]string, error) {
	var ffaScript []string
	// Parse the Dockerfile
	commandList, err := ExtractAllCommandsFromDockerfile(data)
	if err != nil {
		log.Print(err)
	}

	for _, cmd := range commandList {
		switch cmd.Cmd {
		case "run":
			results, err := TranslateShellScript(strings.Join(cmd.Value, " "))
			if err != nil {
				return ffaScript, err
			}
			ffaScript = append(ffaScript, results...)
			break
		case "workdir":
			ffaScript = append(ffaScript, "cd '"+cmd.Value[0]+"';")
			break
		case "copy":
			if len(cmd.Value) == 2 {
				ffaScript = append(ffaScript, "cp '"+cmd.Value[0]+"' '"+cmd.Value[1]+"';")
			}
			break
		}
	}
	return ffaScript, nil
}

func TranslateShellScript(data string) ([]string, error) {
	var ffaList []string
	in := strings.NewReader(data)
	f, err := syntax.NewParser().Parse(in, "")
	if err != nil {
		return nil, err
	}

	ffaVarCounter := 0
	var varbank = make(map[string]string)

	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.Assign:
			// Check if varname is in bank
			if x.Name != nil {
				ffaVar, ok := varbank[x.Name.Value]
				if !ok {
					ffaVar = "$x" + strconv.Itoa(ffaVarCounter)
					// increment x? variable name
					ffaVarCounter++
				}

				// If RHS is unknown use 'INPUT'
				rhs := x.Value
				if rhs == nil || len(rhs.Parts) == 0 {
					return false
				}
				if _, ok = rhs.Parts[0].(*syntax.Lit); !ok {
					// If RHS is not of type Lit, then we use INPUT
					ffaList = append(ffaList, fmt.Sprintf("%s = INPUT;", ffaVar))
				} else {
					ffaList = append(ffaList, fmt.Sprintf("%s = '%s';", ffaVar, rhs.Lit()))
				}
			}
			break
		case *syntax.CallExpr:
			// Skip if empty command
			if len(x.Args) == 0 {
				return true
			}

			// We only handle most common commands
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
					// TODO: handle arguments with variables
					ffaList = append(ffaList, fmt.Sprintf("mkdir '%s';", s))
				}
				break
			case "rm":
			case "rmdir":
				// TODO: check for flags
				// TODO: check for -r and use rmr
				for _, s := range x.Args[1:] {
					ffaList = append(ffaList, fmt.Sprintf("rmr '%s';", s.Lit()))
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
				ffaList = append(ffaList, fmt.Sprintf("rmr '%s';", arg1))
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
				if len(x.Args) == 1 {
					// Typically 'cd' with no args with go to user's home directory...
					ffaList = append(ffaList, "cd '/';")
				} else {
					arg1 := x.Args[1].Lit()
					ffaList = append(ffaList, fmt.Sprintf("cd '%s';", arg1))
				}
				break
			case "wget":
				command := literize(x.Args)
				command, args := extractFlag(command, "-O", 1)
				if args != nil {
					// if -O is present, touch full path
					ffaList = append(ffaList, fmt.Sprintf("touch '%s';", args[1]))
				} else {
					command = removeFlagsLit(command)
					// if -O is not present, we don't always know what the filename will be
					//ffaList = append(ffaList, fmt.Sprintf("touch '%s';", filepath.Base(command[1])))
				}
				break
			case "curl":
				command := literize(x.Args)
				command, args := extractFlag(command, "-O", 1)
				if args != nil {
					ffaList = append(ffaList, fmt.Sprintf("touch '%s';", args[1]))
				}
				break
			case "chmod":
				command := literize(x.Args)
				command = removeFlagsLit(command)
				if len(command) >= 3 {
					for _, filename := range command[2:] {
						ffaList = append(ffaList, fmt.Sprintf("assert(exists '%s');", filename))
					}
				}
				break
			case "source":
				fallthrough
			case "python":
				fallthrough
			case "python2":
				fallthrough
			case "python3":
				command := literize(x.Args)
				command = removeFlagsLit(command)
				if len(command) >= 2 {
					ffaList = append(ffaList, fmt.Sprintf("assert(exists '%s');", command[1]))
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
			default:
				// if strings.HasPrefix("./")
				if m, err := regexp.MatchString(`^\.*?/`, cmd); err != nil {
					log.Fatal(err)
				} else if m {
					// Assert that unknown scripts/binaries exists if relative or absolute path is invoked
					ffaList = append(ffaList, fmt.Sprintf("assert(exists '%s');", cmd))
				}
			}
			break
		case *syntax.IfClause:
			// TODO: handle if clause
			// Axeldnahcram_zsh_install.sh
		case *syntax.WhileClause:
		case *syntax.ForClause:
		case *syntax.CaseClause:
		case *syntax.Block:
		case *syntax.Subshell:
		case *syntax.BinaryCmd:
		case *syntax.FuncDecl:
		case *syntax.ArithmCmd:
		case *syntax.TestClause:
		case *syntax.DeclClause:
		case *syntax.LetClause:
		case *syntax.TimeClause:
		case *syntax.CoprocClause:
		default:
		}
		return true
	})
	return ffaList, nil
}
