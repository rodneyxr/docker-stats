package ffa

import (
	"fmt"
	"github.com/mvdan/sh/syntax"
	"path/filepath"
	"strconv"
	"strings"
)

func AnalyzeShellCommand(cmd string) ([]string, error) {
	var ffaList []string
	in := strings.NewReader(cmd)
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
				ffaList = append(ffaList, fmt.Sprintf("touch '%s'", filepath.Base(args[0])))
				//arg1 := x.Args[1].Lit()
				//if !strings.HasPrefix("-", arg1) {
				//	ffaList = append(ffaList, fmt.Sprintf("touch '%s';", filepath.Base(arg1)))
				//}
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
	marker := -1
	for i, s := range command {
		if flag == s {
			marker = i
		}
	}
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

func removeFlags(arguments []*syntax.Word) []string {
	var args []string
	for _, arg := range arguments {
		if !strings.HasPrefix(arg.Lit(), "-") {
			args = append(args, arg.Lit())
		}
	}
	return args
}
