module github.com/containerd/go-cni/integration

go 1.17

require (
	github.com/containerd/continuity v0.2.2
	github.com/containerd/go-cni v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.0
)

require (
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/containernetworking/cni v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/containerd/go-cni => ../
