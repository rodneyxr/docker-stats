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
	"github.com/rodneyxr/docker-stats/ffa"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var repoURL string

// infoCmd represents the list command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information about a docker GitHub repository",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the existing results
		repoList := ffa.LoadRepos(resultsFile)
		repoMap := make(map[string]ffa.Repo)
		for _, repo := range repoList {
			repoMap[repo.URL] = repo
		}

		// Generate a list of Golang repos
		var goRepos []ffa.Repo
		if repoURL != "" {
			goRepos = append(goRepos, repoMap[repoURL])
		} else {
			for _, repo := range repoList {
				if len(repo.Languages) > 0 && repo.Languages[0].Name == "Go" {
					goRepos = append(goRepos, repo)
				}
			}
		}

		// TODO: Only handle one repo provided by argument
		for i, repo := range goRepos {
			fmt.Printf("%d: %s\n", i, repo.URL)

			// For each first Dockerfile in each repo
			if len(repo.Dockerfiles) > 0 {
				// Parse the Dockerfile
				runCommandList, err := ffa.ExtractRunCommandsFromDockerfile(repo.Dockerfiles[0])
				if err != nil {
					log.Print(err)
				}
				for _, cmd := range runCommandList {
					ffa.AnalyzeShellCommand(strings.Join(cmd.Value, " "))
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVar(&repoURL, "repo", "", "Git repo URL")
}
