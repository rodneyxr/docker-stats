package docker

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

var sampleDockerfile = `FROM golang:1.10.0
RUN apt-get update && apt-get install -y --no-install-recommends \
                openssh-client \
                rsync \
                fuse \
                sshfs \
        && rm -rf /var/lib/apt/lists/*
RUN go get  golang.org/x/lint/golint \
            github.com/mattn/goveralls \
            golang.org/x/tools/cover
RUN git clone https://github.com/rodneyxr/repo
RUN wget https://github.com/rodneyxr/testfile.txt
RUN curl -XGET https://google.com -O google_output
ENV USER root
WORKDIR /go/src/github.com/docker/machine
COPY . ./
RUN mkdir bin
`

func TestParser(t *testing.T) {
	// Parse the Dockerfile
	runCommandList, err := ExtractRunCommandsFromDockerfile(sampleDockerfile)
	if err != nil {
		log.Print(err)
	}
	ffa := strings.Builder{}
	for _, cmd := range runCommandList {
		commands := AnalyzeRunCommand(cmd)
		for _, ffaCommand := range commands {
			ffa.WriteString(ffaCommand + "\n")
		}
	}
	fmt.Println(ffa.String())
}
