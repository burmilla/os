module github.com/containerd/continuity/cmd/continuity

go 1.17

require (
	bazil.org/fuse v0.0.0-20200524192727-fb710f7dfd05
	github.com/containerd/continuity v0.0.0-00010101000000-000000000000 // see replace
	github.com/dustin/go-humanize v1.0.0
	github.com/golang/protobuf v1.5.2
	github.com/opencontainers/go-digest v1.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.4.0
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220405210540-1e041c57c461 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

// use local source for the main module
replace github.com/containerd/continuity => ../../
