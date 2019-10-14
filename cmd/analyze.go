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

package cmd

import (
	"fmt"
	"github.com/rodneyxr/docker-stats/docker"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var dockerfilePath string
var resultsDir string

// analzyeCmd represents the list command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze a dockerfile or directory full of dockerfiles",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("%d: %s\n", i, repo.URL)

		var files []string
		if info, err := os.Stat(dockerfilePath); err == nil && info.IsDir() {
			if err := filepath.Walk(dockerfilePath, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files = append(files, path)
				}
				return err
			}); err != nil {
				panic(err)
			}
		}

		// Create the results directory
		_ = os.Mkdir(resultsDir, os.ModeDir)

		for i, dockerFilename := range files {
			fmt.Printf("%d: %s\n", i, dockerFilename)
			data, err := ioutil.ReadFile(dockerFilename)
			if err != nil {
				log.Println(err)
				continue
			}
			// Parse the Dockerfile
			commandList, err := docker.ExtractAllCommandsFromDockerfile(string(data))
			if err != nil {
				log.Print(err)
			}
			var ffaScript []string
			for _, cmd := range commandList {
				switch cmd.Cmd {
				case "run":
					ffaScript = append(ffaScript, docker.AnalyzeRunCommand(cmd)...)
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

			// Save the ffa script to a file
			ffaFilename := filepath.Join(resultsDir, filepath.Base(dockerFilename)+".ffa")
			ffaScriptData := []byte(strings.Join(ffaScript, "\n"))
			if err = ioutil.WriteFile(ffaFilename, ffaScriptData, os.ModePerm); err != nil {
				log.Print(err)
				continue
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVar(&dockerfilePath, "dockerfile", "", "dockerfile or dockerfile directory")
	analyzeCmd.Flags().StringVar(&resultsDir, "results", "results", "directory to save results")
}
