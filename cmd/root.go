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
	"os"

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
	//Run: func(cmd *cobra.Command, args []string) {},
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
	rootCmd.PersistentFlags().StringVar(&reposFile, "repos", "repos.yaml", "list of repos to update")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token", "token.txt", "file containing GitHub access token")
	rootCmd.PersistentFlags().StringVar(&resultsFile, "results", "results.json", "output file as json")
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
