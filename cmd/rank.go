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
	"sort"
	"strings"
)

var uniqueFlag bool

// rankCmd represents the rank command
var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Ranks the number of occurrences for each run binary executed by the docker RUN command",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the existing results
		repoList := ffa.LoadRepos(resultsFile)
		repoMap := make(map[string]ffa.Repo)
		keywordMap := make(map[string]int)
		for _, repo := range repoList {
			repoMap[repo.URL] = repo
		}

		var goRepos []ffa.Repo
		for _, repo := range repoList {
			if len(repo.Languages) > 0 && repo.Languages[0].Name == "Go" {
				goRepos = append(goRepos, repo)
			}
		}

		for i, repo := range goRepos {
			localKeywordMap := make(map[string]int)
			fmt.Printf("%d: %s\n", i, repo.URL)

			// For each Dockerfile in each repo
			// For each first Dockerfile in each repo
			for _, dockerfile := range repo.Dockerfiles {
				//if len(repo.Dockerfiles) > 0 {
				runCommandList, err := ffa.ExtractRunCommandsFromDockerfile(dockerfile)
				if err != nil {
					log.Print(err)
					continue
				}

				for _, cmd := range runCommandList {
					commandName := strings.Split(cmd.Value[0], " ")[0]
					if uniqueFlag {
						// Only count one occurrence of a command per project if unique flag is provided
						localKeywordMap[commandName] = 1
					} else {
						localKeywordMap[commandName] += 1
					}
					fmt.Printf("%s(%d)\n", commandName, keywordMap[commandName])
				}
			}

			// Add local keywords map to the total keywords map
			for k, v := range localKeywordMap {
				keywordMap[k] += v
			}
		}

		// Sort the occurrence list
		type kv struct {
			Key   string
			Value int
		}
		var ss []kv
		for k, v := range keywordMap {
			ss = append(ss, kv{k, v})
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		// Display the results
		fmt.Println()
		fmt.Println("================================================================================")
		fmt.Println()
		for _, kv := range ss {
			fmt.Printf("%d:\t%s\n", kv.Value, kv.Key)
		}
	},
}

func init() {
	rootCmd.AddCommand(rankCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rankCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	rankCmd.Flags().BoolVar(&uniqueFlag, "unique", false, "Only allow one command per project to be ranked")
}
