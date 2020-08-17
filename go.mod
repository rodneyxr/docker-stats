module github.com/rodneyxr/ffatoolkit

go 1.14

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200227233006-38f52c9fec82

require (
	github.com/asottile/dockerfile v3.1.0+incompatible
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/moby/buildkit v0.7.1 // indirect
	github.com/mvdan/sh v2.6.4+incompatible
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	mvdan.cc/sh v2.6.4+incompatible // indirect
)
