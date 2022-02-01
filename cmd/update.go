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
	"encoding/json"
	"github.com/rodneyxr/ffatoolkit/ffa"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates/downloads the repo cache",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the existing cache
		repoList, err := ffa.LoadRepoCache(cacheFile)
		if err != nil {
			repoList = []ffa.Repo{}
		}

		repoMap := make(map[string]ffa.Repo)
		for _, repo := range repoList {
			repoMap[repo.URL] = repo
		}

		repoURLs := viper.GetStringSlice("repos")
		ctx := context.Background()
		if err != nil {
			log.Fatal(err)
		}
		client := ffa.CreateClient(ctx, gitToken)

		// Create the info list to hold all the results
		var results []ffa.Repo

		for i, repoURL := range repoURLs {
			// Display progress to the user
			log.Printf("(%d/%d) %s\n", i+1, len(repoURLs), repoURL)

			// Check if the repo exists already
			repo, ok := repoMap[repoURL]
			if !ok {
				// Create and add the repo object to the result set
				repo = ffa.NewRepo(ctx, client, repoURL)
				repoMap[repoURL] = repo
			}
			//if repo.Languages == nil {
			//	if err := ffa.LoadLanguages(ctx, client, &repo); err != nil {
			//		log.Println(err)
			//	}
			//}
			//if repo.Dockerfiles == nil {
			//	ffa.LoadDockerfiles(ctx, client, &repo)
			//}
			results = append(results, repo)
		}

		// Write the results to a cache file
		repoCacheJson, _ := json.MarshalIndent(results, "", "    ")
		if err := ioutil.WriteFile(cacheFile, repoCacheJson, 0644); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&gitToken, "token", "", "GitHub access token")
}
