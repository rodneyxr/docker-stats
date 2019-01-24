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

package git

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Repos holds information about a GitHub repository
type Repo struct {
	URL         string     `json:"url"`
	Owner       string     `json:"owner"`
	Repo        string     `json:"repo"`
	Languages   []Language `json:"languages"`
	Dockerfiles []string   `json:"dockerfiles"`
	Images      []string   `json:"images"`
}

// Language holds information about a language used by a GitHub repository
type Language struct {
	Name       string  `json:"name"`
	Percentage float32 `json:"percentage"`
}

// LoadRepos will load a list of Repo objects from a json file.
func LoadRepos(filePath string) []Repo {
	var repos []Repo
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("repo file does not exist: %s", filePath)
		return repos
	}
	err = json.Unmarshal(data, &repos)
	if err != nil {
		log.Fatal(err)
	}
	return repos
}

// NewRepo creates the repo object given a URL
func NewRepo(ctx context.Context, client *github.Client, url string) Repo {
	tokens := strings.Split(url, "/")
	owner := tokens[len(tokens)-2]
	repo := tokens[len(tokens)-1]
	repoInfo := Repo{
		URL:   fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		Owner: owner,
		Repo:  repo,
	}
	LoadLanguages(ctx, client, &repoInfo)
	LoadDockerfiles(ctx, client, &repoInfo)
	return repoInfo
}

// LoadDockerfiles loads all files matching dockerfile in their path
func LoadDockerfiles(ctx context.Context, client *github.Client, repoInfo *Repo) {
	// Get the list of docker files in the repo
	var codeResults *github.CodeSearchResult
	var success bool
	var err error
	for !success {
		codeResults, _, err = client.Search.Code(ctx, fmt.Sprintf("dockerfile+in:path+repo:%s/%s", repoInfo.Owner, repoInfo.Repo), &github.SearchOptions{})
		if rateLimitError, ok := err.(*github.RateLimitError); ok {
			waitTime := int(math.Abs(float64(rateLimitError.Rate.Reset.Second())-float64(time.Now().Second()))) + 1
			log.Printf("error searching code: api rate limit reached - sleeping for %d second(s)...", waitTime)
			time.Sleep(time.Duration(waitTime) * time.Second)
		} else {
			success = true
		}
	}

	// Search for FROM statements in each docker file
	var images []string
	var dockerfiles []string
	for _, result := range codeResults.CodeResults {
		fmt.Println("\t", *result.Path)
		success = false
		var contents io.ReadCloser
		for !success {
			contents, err = client.Repositories.DownloadContents(ctx, repoInfo.Owner, repoInfo.Repo, *result.Path, &github.RepositoryContentGetOptions{})
			if rateLimitError, ok := err.(*github.RateLimitError); ok {
				waitTime := int(math.Abs(float64(rateLimitError.Rate.Reset.Second())-float64(time.Now().Second()))) + 1
				log.Printf("error downloading contents: api rate limit reached - sleeping for %d second(s)...", waitTime)
				time.Sleep(time.Duration(waitTime) * time.Second)
			} else {
				success = true
			}
		}

		// Save the dockerfile
		var data []byte
		data, err = ioutil.ReadAll(contents)
		dockerfiles = append(dockerfiles, string(data))

		// Extract the images
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			text := scanner.Text()
			if strings.HasPrefix(text, "FROM ") {
				image := strings.Split(text, "FROM ")[1]
				image = strings.Trim(image, " ")
				images = append(images, image)
				fmt.Println("\t\tFROM", image)
			}
		}
	}
	repoInfo.Dockerfiles = dockerfiles
	repoInfo.Images = images
}

// LoadLanguages loads all the
func LoadLanguages(ctx context.Context, client *github.Client, repoInfo *Repo) {
	// List the languages for the repo
	languages, _, err := client.Repositories.ListLanguages(ctx, repoInfo.Owner, repoInfo.Repo)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the total number bytes to be used later when computing the language percentage
	totalBytes := 0
	for _, byteData := range languages {
		totalBytes += byteData
	}

	// Create the list of languages along with their percentages
	var languageInfos []Language
	for language, byteData := range languages {
		languageInfos = append(languageInfos, Language{
			Name:       language,
			Percentage: float32(byteData) / float32(totalBytes) * 100,
		})
	}

	// Sort the language info list from highest to lowest percentage
	sort.Slice(languageInfos, func(i, j int) bool {
		return languageInfos[i].Percentage > languageInfos[j].Percentage
	})
	repoInfo.Languages = languageInfos
}

// CreateClient authenticates and creates a client to use
func CreateClient(ctx context.Context, accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

// GetRepoList creates a list of repos from a json array file
func GetRepoList(jsonFilePath string) []string {
	// Open the input file
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// Read the data from the file
	data, _ := ioutil.ReadAll(jsonFile)
	var repos []string
	if err = json.Unmarshal(data, &repos); err != nil {
		log.Fatal("failed to open")
	}
	return repos
}
