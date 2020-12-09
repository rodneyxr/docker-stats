module github.com/rodneyxr/ffatoolkit

go 1.14

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200227233006-38f52c9fec82

require (
	github.com/asottile/dockerfile v3.1.0+incompatible
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/moby/buildkit v0.7.1 // indirect
	github.com/mvdan/sh v2.6.4+incompatible
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200519105757-fe76b779f299 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	mvdan.cc/sh v2.6.4+incompatible // indirect
)
