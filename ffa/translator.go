package ffa

// FIXME: if [... will be seen as a command and assert that '[' does not exist
// FIXME: quoted values do not appear in translation. ex: touch 'a' will translate to touch ''

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"mvdan.cc/sh/v3/syntax"
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

type stack []syntax.Node

func (s stack) Push(node syntax.Node) stack {
	return append(s, node)
}

func (s stack) Pop() (stack, syntax.Node) {
	// If stack is empty return the empty stack with a nil node
	if len(s) == 0 {
		return s, nil
	}
	return s[:len(s)-1], s[len(s)-1]
}

// nodes is a stack of syntax.nodes for scopes
var nodes stack
var scopeCounter int

func appendFFAList(ffaList []string, commandStr string) []string {
	commandStr = strings.Repeat("    ", scopeCounter) + commandStr
	ffaList = append(ffaList, commandStr)
	return ffaList
}

func isScope(node syntax.Node, data string) bool {
	var empty syntax.Pos
	switch x := node.(type) {
	case *syntax.IfClause:
		// elif and else if's should not count as a new scope
		if x.ThenPos == empty && x.FiPos != empty {
			return false
		} else if x.ThenPos != empty && string(data[x.Pos().Offset()]) == "e" {
			return false
		}
		return true
	case *syntax.WhileClause, *syntax.ForClause:
		return true
	}
	return false
}

func isElseScope(node syntax.Node) bool {
	var empty syntax.Pos
	switch x := node.(type) {
	case *syntax.IfClause:
		if x.ThenPos == empty {
			return true
		}
	}
	return false
}

func TranslateShellScript(data string) ([]string, error) {
	var ffaList []string
	in := strings.NewReader(data)
	parser := syntax.NewParser()
	f, err := parser.Parse(in, "")
	if err != nil {
		return nil, err
	}
	ffaVarCounter := 0
	scopeCounter = 0
	nodes = stack{}
	var varbank = make(map[string]string)

	syntax.Walk(f, func(node syntax.Node) bool {
		if node == nil {
			var x syntax.Node
			nodes, x = nodes.Pop()
			if isElseScope(x) {
				return false
			} else if isScope(x, data) {
				scopeCounter--
				ffaList = appendFFAList(ffaList, "}")
				return false
			}
		} else {
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
						ffaList = appendFFAList(ffaList, fmt.Sprintf("%s = INPUT;", ffaVar))
					} else {
						ffaList = appendFFAList(ffaList, fmt.Sprintf("%s = '%s';", ffaVar, rhs.Lit()))
					}
				}
				break
			case *syntax.CallExpr:
				// Skip if empty command
				if len(x.Args) == 0 {
					break
				}

				// We only handle most common commands
				cmd := x.Args[0].Lit()
				switch cmd {
				case "read":

				case "touch":
					// Create a touch statement for each argument
					for _, s := range x.Args[1:] {
						ffaList = appendFFAList(ffaList, fmt.Sprintf("touch '%s';", s.Lit()))
					}
					break
				case "mkdir":
					args := removeFlags(x.Args)
					for _, s := range args[1:] {
						// TODO: handle arguments with variables
						ffaList = appendFFAList(ffaList, fmt.Sprintf("mkdir '%s';", s))
					}
					break
				case "rm":
					fallthrough
				case "rmdir":
					// TODO: check for flags
					// TODO: check for -r and use rmr
					args := removeFlags(x.Args)
					for _, s := range args[1:] {
						ffaList = appendFFAList(ffaList, fmt.Sprintf("rmr '%s';", s))
					}
					break
				case "cp":
					args := removeFlags(x.Args)
					arg1, arg2 := args[1], args[2]
					ffaList = appendFFAList(ffaList, fmt.Sprintf("cp '%s' '%s';", arg1, arg2))
					break
				case "mv":
					args := removeFlags(x.Args)
					arg1, arg2 := args[1], args[2]
					ffaList = appendFFAList(ffaList, fmt.Sprintf("cp '%s' '%s';", arg1, arg2))
					ffaList = appendFFAList(ffaList, fmt.Sprintf("rmr '%s';", arg1))
					break
				case "git":
					// TODO: handle git rm
					arg1 := x.Args[1].Lit()
					if arg1 == "clone" {
						dirname := filepath.Base(x.Args[2].Lit())
						ffaList = appendFFAList(ffaList, fmt.Sprintf("mkdir '%s';", dirname))
					}
					break
				case "cd":
					if len(x.Args) == 1 {
						// Typically 'cd' with no args with go to user's home directory...
						ffaList = appendFFAList(ffaList, "cd '/';")
					} else {
						arg1 := x.Args[1].Lit()
						ffaList = appendFFAList(ffaList, fmt.Sprintf("cd '%s';", arg1))
					}
					break
				case "wget":
					command := literize(x.Args)
					command, args := extractFlag(command, "-O", 1)
					if args != nil {
						// if -O is present, touch full path
						ffaList = appendFFAList(ffaList, fmt.Sprintf("touch '%s';", args[1]))
					} else {
						command = removeFlagsLit(command)
						// if -O is not present, we don't always know what the filename will be
						//ffaList = appendFFAList(ffaList, fmt.Sprintf("touch '%s';", filepath.Base(command[1])))
					}
					break
				case "curl":
					command := literize(x.Args)
					command, args := extractFlag(command, "-O", 1)
					if args != nil {
						ffaList = appendFFAList(ffaList, fmt.Sprintf("touch '%s';", args[1]))
					}
					break
				case "chmod":
					command := literize(x.Args)
					command = removeFlagsLit(command)
					if len(command) >= 3 {
						for _, filename := range command[2:] {
							ffaList = appendFFAList(ffaList, fmt.Sprintf("assert(exists '%s');", filename))
						}
					}
					break
				case "file":
					fallthrough
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
						ffaList = appendFFAList(ffaList, fmt.Sprintf("assert(exists '%s');", command[1]))
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
						ffaList = appendFFAList(ffaList, fmt.Sprintf("assert(exists '%s');", cmd))
					} else {
						// Ignore if conditions
						if len(cmd) > 0 && cmd[0] != '[' {
							// Assert that the binary does not exist locally
							ffaList = appendFFAList(ffaList, fmt.Sprintf("assert(! exists '%s');", cmd))
						}
					}
				}
			case *syntax.IfClause:
				var empty syntax.Pos

				// Condition to check if IfClause node is an Else statement
				if x.ThenPos == empty && x.FiPos != empty {
					scopeCounter--
					ffaList = appendFFAList(ffaList, "} else {")
					scopeCounter++
					// Condition to check if IfClause node is an elif statement
				} else if x.ThenPos != empty && string(data[x.Pos().Offset()]) == "e" {
					scopeCounter--
					ffaList = appendFFAList(ffaList, "} else if (other) {")
					scopeCounter++
				} else {
					ffaList = appendFFAList(ffaList, "if (other) {")
				}
			case *syntax.WhileClause:
				ffaList = appendFFAList(ffaList, "while (other) {")
			case *syntax.ForClause:
				ffaList = appendFFAList(ffaList, "while (other) {")
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
			}

			if isScope(node, data) {
				scopeCounter++
			}
			nodes = nodes.Push(node)
		}
		return true
	})

	return ffaList, nil
}
