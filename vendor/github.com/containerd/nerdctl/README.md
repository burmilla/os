# Customized nerdctl
This branch is based on nerdctl version 1.2.0 but it is modified on way that we use [trash](https://github.com/burmilla/trash) instead of Go modules and code is cleaned up on way that it can be build with modules which trash was able to vendor. 

Re-vendor command:
```bash
GO111MODULE=off trash -d
```

Build command:
```bash
make binaries
```