// Copyright © 2018 Rodney Rodriguez
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
	"github.com/asottile/dockerfile"
	"github.com/rodneyxr/docker-stats/git"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

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

		for i, repo := range goRepos {
			fmt.Printf("%d: %s\n", i, repo.URL)

			// For each first Dockerfile in each repo
			if len(repo.Dockerfiles) > 0 {
				// Parse the Dockerfile
				reader := strings.NewReader(repo.Dockerfiles[0])
				commandList, err := dockerfile.ParseReader(reader)
				if err != nil {
					log.Print(err)
					continue
				}

				// Print all commands in the Dockerfile
				for _, cmd := range commandList {
					fmt.Println(cmd.Cmd)
				}

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
