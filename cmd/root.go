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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rodneyxr/docker-stats/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reposFile string
var resultsFile string
var tokenFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "docker-stats",
	Short: "A data collection program for github repositories using Docker.",

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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&reposFile, "repos", "repos.yaml", "config file (default repos.yaml)")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token", "token.txt", "config file (default token.txt)")
	rootCmd.PersistentFlags().StringVar(&resultsFile, "results", "results.json", "config file (default results.json)")
}

// initConfig reads in the list of repos defined in a yaml file.
func initConfig() {
	if reposFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(reposFile)
	} else {
		viper.SetConfigType("yaml")
		viper.SetConfigName("repos")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using repos file:", viper.ConfigFileUsed())
	}
}
