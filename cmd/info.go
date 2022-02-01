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
	"context"
	"fmt"
	"github.com/rodneyxr/ffatoolkit/ffa"
	"github.com/spf13/cobra"
	"log"
)

var repoURL string
var repoLanguage string
var gitToken string

// infoCmd represents the list command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information about dockerfiles from GitHub repositories",
	Run: func(cmd *cobra.Command, args []string) {
		var repoList []ffa.Repo
		var err error

		// if repoURL was provided load the repo
		if repoURL != "" {
			ctx := context.Background()
			client := ffa.CreateClient(ctx, gitToken)
			log.Println("Fetching repo info:", repoURL)
			repo := ffa.NewRepo(ctx, client, repoURL)
			repoList = append(repoList, repo)
		} else {
			// Load repos from the cache
			repoList, err = ffa.LoadRepoCache(cacheFile)
			if err != nil {
				log.Fatal(err)
			}
			if len(repoList) == 0 {
				fmt.Println("No repos in cache.")
				fmt.Println("Try running the following command first:")
				fmt.Println("    $ ffatoolkit update")
				return
			}
		}

		// Filter if filter-lang was provided
		var filteredRepos []ffa.Repo
		if repoLanguage != "" {
			for _, repo := range repoList {
				if len(repo.Languages) > 0 && repo.Languages[0].Name == repoLanguage {
					filteredRepos = append(filteredRepos, repo)
				}
			}
		} else {
			filteredRepos = repoList
		}

		if len(filteredRepos) == 0 {
			fmt.Println("No repos matched the filter.")
			return
		}

		// For each filtered repo, extract RUN commands from
		for i, repo := range filteredRepos {
			fmt.Printf("%d: %s\n", i, repo.URL)

			// For each first Dockerfile in each repo
			if len(repo.Dockerfiles) > 0 {
				// Parse the Dockerfile
				ffa, err := ffa.TranslateDockerfile(repo.Dockerfiles[0])
				if err != nil {
					log.Print(err)
				} else {
					log.Println(ffa)
				}
			} else {
				log.Println("No Dockerfile found.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVar(&repoURL, "repo", "", "Git repo URL")
	infoCmd.Flags().StringVar(&repoLanguage, "filter-lang", "", "Repo language to filter")
	infoCmd.Flags().StringVar(&gitToken, "token", "", "GitHub access token")
}
