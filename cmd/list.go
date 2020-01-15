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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/asottile/dockerfile"
	"github.com/rodneyxr/docker-stats/git"
	"github.com/spf13/cobra"
)

var saveFlag bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the existing results
		repoList := git.LoadRepos(resultsFile)
		repoMap := make(map[string]git.Repo)
		for _, repo := range repoList {
			repoMap[repo.URL] = repo
		}

		var goRepos []git.Repo

		for _, repo := range repoList {
			if len(repo.Languages) > 0 && repo.Languages[0].Name == "Go" {
				goRepos = append(goRepos, repo)
				/*fmt.Printf("%d: %s\n", numberOfGoRepos, repo.URL)*/
			}
		}

		// Create the directory to save the docker files to
		if saveFlag {
			if err := os.MkdirAll("dockerfiles", os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}

		// Iterate through all repos
		for i, repo := range goRepos {
			fmt.Printf("%d: %s\n", i, repo.URL)

			// For each first Dockerfile in each repo
			if len(repo.Dockerfiles) > 0 {
				// Save the dockerfile to a file
				if saveFlag {
					dockerFilename := strings.Join([]string{repo.Owner, repo.Repo, "Dockerfile"}, "_")
					if err := ioutil.WriteFile(filepath.Join("dockerfiles", dockerFilename), []byte(repo.Dockerfiles[0]), os.ModePerm); err != nil {
						log.Fatal(err)
					}
				}

				// Parse the Dockerfile
				reader := strings.NewReader(repo.Dockerfiles[0])
				commandList, err := dockerfile.ParseReader(reader)
				if err != nil {
					log.Print(err)
					continue
				}

				// Print all commands in the Dockerfile
				for _, cmd := range commandList {
					if cmd.Cmd == "run" {
						fmt.Println(cmd.Cmd, cmd.Value)
					}
				}

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Provide this flag if the Dockerfiles should be saved to a file
	listCmd.Flags().BoolVarP(&saveFlag, "save", "s", false, "Save the Dockerfiles to a file")
}
