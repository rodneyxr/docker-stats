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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/rodneyxr/docker-stats/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates/downloads the results using the GitHub REST API",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the existing results
		repoList := git.LoadRepos(resultsFile)
		repoMap := make(map[string]git.Repo)
		for _, repo := range repoList {
			repoMap[repo.URL] = repo
		}

		repoURLs := viper.GetStringSlice("repos")
		ctx := context.Background()
		token, err := ioutil.ReadFile(tokenFile)
		if err != nil {
			log.Fatal(err)
		}
		client := git.CreateClient(ctx, string(token))

		// Create the info list to hold all the results
		var results []git.Repo

		for i, repoURL := range repoURLs {
			// Display progress to the user
			fmt.Printf("(%d/%d) %s\n", i+1, len(repoURLs), repoURL)

			// Check if the repo exists already
			repo, ok := repoMap[repoURL]
			if !ok {
				// Create and add the repo object to the result set
				repo = git.NewRepo(ctx, client, repoURL)
				repoList = append(repoList, repo)
				repoMap[repoURL] = repo
			}
			if repo.Languages == nil {
				git.LoadLanguages(ctx, client, &repo)
			}
			if repo.Dockerfiles == nil {
				git.LoadDockerfiles(ctx, client, &repo)
			}
			results = append(results, repo)
		}

		// Write the results to a file
		repoInfoJson, _ := json.MarshalIndent(results, "", "    ")
		if err := ioutil.WriteFile(resultsFile, repoInfoJson, 0644); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
