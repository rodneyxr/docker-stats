// Copyright © 2020 Rodney Rodriguez
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
	"github.com/rodneyxr/ffatoolkit/ffa"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var fileTypeFlag string
var filepathFlag string
var resultsDir string

// translateCmd represents the list command
var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate scripts to FFAL",
	Run: func(cmd *cobra.Command, args []string) {
		var files []string

		// Stat the file
		info, err := os.Stat(filepathFlag)
		if err != nil {
			cmd.PrintErrln("could not read " + filepathFlag)
			os.Exit(1)
		}

		if info.IsDir() {
			// If the file is a directory, add all files to the files list
			if err := filepath.Walk(filepathFlag, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files = append(files, path)
				}
				return err
			}); err != nil {
				panic(err)
			}
		} else {
			// if it is not a directory, the file will be the only one in the list
			abs, _ := filepath.Abs(filepathFlag)
			files = append(files, abs)
		}

		// Create the results directory
		_ = os.Mkdir(resultsDir, os.ModeDir)

		for _, filename := range files {
			// Read the file data
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Println(err)
				continue
			}

			var ffaScript []string
			switch fileTypeFlag {
			case "docker":
				ffaScript, err = ffa.TranslateDockerfile(string(data))
				if err != nil {
					log.Println(err)
					continue
				}
				break
			case "shell":
				ffaScript, err = ffa.TranslateShellScript(string(data))
				if err != nil {
					// skip this file to avoid a partially translated file
					log.Printf("failed to parse %s: %s", filename, err)
					continue
				}
				//ffaScript = append(ffaScript, results...)
				break
			default:
				log.Fatal("unsupported file type")
			}

			// Save the ffa script to a file
			ffaFilename := filepath.Join(resultsDir, filepath.Base(filename)+".ffa")
			ffaScriptData := []byte(strings.Join(ffaScript, "\n"))
			if err = ioutil.WriteFile(ffaFilename, ffaScriptData, os.ModePerm); err != nil {
				log.Print(err)
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().StringVar(&fileTypeFlag, "type", "shell", "type of file to analyze (shell or docker)")
	translateCmd.Flags().StringVar(&filepathFlag, "filepath", "", "path to file or directory to analyze")
	translateCmd.Flags().StringVar(&resultsDir, "results", "results", "directory to save results")
	_ = translateCmd.MarkFlagRequired("filepath")
}
