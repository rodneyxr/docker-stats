# File Flow Analysis Toolkit
This tool can search through github repositories for dockerfiles and shell scripts and attempt to convert them to the
FFA scripts to be analyzed by the FFA framework.

## Usage
```bash
A data collection program for github repositories using Docker.

Usage:
  ffatoolkit [command]

Available Commands:
  help        Help about any command
  info        Show information about dockerfiles from GitHub repositories
  list        Lists the Dockerfiles found in each repo in the repo file
  rank        Ranks the number of occurrences for each run binary executed by the docker RUN command
  translate   Translate scripts to FFAL
  update      Updates/downloads the repo cache

Flags:
  -h, --help                 help for ffatoolkit
      --repos string         list of repos to update (default "repos.yaml")
      --resultsfile string   output file as json (default "results.json")
      --token string         file containing GitHub access token (default "token.txt")

Use "ffatoolkit [command] --help" for more information about a command.
```