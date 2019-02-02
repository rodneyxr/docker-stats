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

	// Print all commands in the Dockerfile
	for _, cmd := range commandList {
		if cmd.Cmd == "run" {
			fmt.Println(cmd.Cmd, cmd.Value)
		}
	}
	return commandList, nil
}

func AnalyzeRunCommand(cmd dockerfile.Command) {
	commandString := strings.Join(cmd.Value, " ")
	in := strings.NewReader(commandString)
	f, err := syntax.NewParser().Parse(in, "")
	if err != nil {
		return
	}
	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.ParamExp:
			x.Param.Value = strings.ToUpper(x.Param.Value)
		case *syntax.arg
		}
		return true
	})
	syntax.NewPrinter().Print(os.Stdout, f)
}
