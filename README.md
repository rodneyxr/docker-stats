# File Flow Analysis Toolkit
This tool can search through github repositories for dockerfiles and shell scripts and attempt to convert them to the
FFA scripts to be analyzed by the FFA framework.

## Usage
```bash
A data collection program for github repositories using Docker.

Usage:
  ffatoolkit [command]

Available Commands:
  analyze     Analyze a dockerfile or directory full of dockerfiles
  help        Help about any command
  info        Show information about a docker GitHub repository
  list        A brief description of your command
  rank        Ranks the number of occurrences for each run binary executed by the docker RUN command
  update      Updates/downloads the results using the GitHub REST API

Flags:
  -h, --help                 help for ffatoolkit
      --repos string         list of repos to update (default "repos.yaml")
      --resultsfile string   output file as json (default "results.json")
      --token string         file containing GitHub access token (default "token.txt")

Use "ffatoolkit [command] --help" for more information about a command.
```