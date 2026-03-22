# CLAUDE.md - BurmillaOS Development Guide

## Project Overview

BurmillaOS is a minimal Linux distribution where **everything runs as a Docker container**, including system services. It is the community-maintained successor to RancherOS. The OS uses a dual-Docker architecture: **System Docker** (PID 1) manages OS-level containers, and **User Docker** runs user workloads in isolation.

## Build System

### Prerequisites
- Docker (for `dapper` containerized builds)
- `make`

### Key Commands
```bash
make build           # Build the ros binary and images
make test            # Run Go tests with race detection
make validate        # Run go vet + go fmt checks
make pr-validation   # CI validation (skips kernel build, runs test + validate)
make release         # Build all release artifacts (ISO, initrd, etc.)
```

### Build Architecture
- Uses **dapper** (containerized build tool) with `Dockerfile.dapper` (Ubuntu 18.04 base)
- Go 1.19.5 with `GO111MODULE=off` (GOPATH-based vendoring, NOT Go modules)
- Cross-compilation targets: `amd64`, `arm64`
- Build scripts live in `scripts/` directory

### Dependency Management
- Uses **trash** tool (not Go modules) - dependencies declared in `trash.conf`
- Vendor directory is checked in (`vendor/`)
- To update dependencies: edit `trash.conf`, then run `make deps`

## Project Structure

```
cmd/                    # CLI entry points
  control/              # Main `ros` CLI (subcommands for OS management)
  cloudinitexecute/     # Cloud-init execution
  cloudinitsave/        # Cloud-init config saving
  init/                 # System initialization
  network/              # Network configuration (netconf)
  power/                # halt, poweroff, reboot, shutdown
  sysinit/              # System init
  respawn/              # Process respawning

pkg/                    # Core libraries
  init/                 # Init subsystem (bootstrap, cloudinit, docker, fsmount, etc.)
  compose/              # Docker Compose integration
  dfs/                  # Docker filesystem utilities
  docker/               # Docker client operations
  hostname/             # Hostname management
  libcompose/           # Internal libcompose replacement (12 subpackages)
  log/                  # Logging
  netconf/              # Network configuration (bridges, bonds, VLANs, DHCP, IPv4LL)
  sysinit/              # System initialization
  util/                 # Utilities (network, versioning)

config/                 # Configuration types, parsing, validation
  cloudinit/            # Cloud-init documentation and schemas

images/                 # Docker container image definitions
  00-rootfs/            # Base root filesystem
  01-base/              # Base image
  02-*/                 # System services (acpid, bootstrap, console, logrotate, syslog)

scripts/                # Build, package, release, and CI scripts
```

## Key Entry Point

`main.go` registers reexec handlers that map binary names to functions:
- `init` -> system init
- `cloud-init-execute/save` -> cloud-init
- `netconf` -> network configuration
- `ros-sysinit` / `ros-bootstrap` -> system/bootstrap init
- Default -> `control.Main()` (the `ros` CLI)

## Testing

```bash
# Run all tests
make test

# Run tests directly (from inside dapper container or with correct GOPATH)
go test -v -cover -tags=test ./...

# Validate code (vet + fmt)
make validate
```

Test files are alongside source: `config/*_test.go`, etc.

## Custom/Forked Dependencies

BurmillaOS maintains forked versions of several critical dependencies under the `burmilla` GitHub organization. These are declared in `trash.conf` with custom repository URLs.

### Why Custom runc?
- **Package**: `github.com/opencontainers/runc` -> fork at `github.com/burmilla/runc.git`
- **Usage**: Only the `libcontainer/user` package is imported (user/group lookup from `/etc/passwd` and `/etc/group`)
- **Used by**: Vendored Docker packages (`docker/go-connections`, `docker/docker/pkg/homedir`)
- **Reason**: Inherited from RancherOS. The fork maintains compatibility with the specific Docker version used by BurmillaOS's System Docker (v17.06.107). The actual runc binary is bundled with System Docker, downloaded as a pre-built binary from `github.com/burmilla/os-system-docker`

### Why Custom netlink?
- **Package**: `github.com/vishvananda/netlink` -> fork at `github.com/burmilla/netlink`
- **Usage**: Extensively used in `pkg/netconf/` for network configuration
- **Features used**: Link management (bridge, bond, VLAN creation), IP address management, routing, IPv4 link-local addressing
- **Reason**: The fork contains patches for BurmillaOS-specific network configuration needs including bridge creation for container networking, VLAN tagging, interface bonding, and zero-config networking (IPv4LL). Migrated from `niusmallnan/netlink` to `burmilla/netlink` in May 2021

### Other Custom Forks
- `burmilla/docker.git` - Docker engine (System Docker patches)
- `burmilla/containerd.git` - containerd runtime
- `burmilla/candiedyaml` - YAML parser
- `burmilla/cli-1` - CLI framework
- `burmilla/libcompose.git` - Docker Compose library (replaced by compose-spec/compose-go + pkg/libcompose)

## Configuration

- **OS config template**: `os-config.tpl.yml` - defines defaults (DHCP, DNS, services)
- **Config types**: `config/types.go` - constants, paths, labels
- **Network config types**: `pkg/netconf/types.go`

## CI/CD

- **PR validation**: `.github/workflows/pull-request-validation.yml` (runs `make pr-validation`)
- **Releases**: `.github/workflows/create-release.yml` (manual trigger, builds + publishes)

## Common Pitfalls

- This project uses `GO111MODULE=off` - do NOT add `go.mod`/`go.sum` files
- Dependencies are managed via `trash.conf`, not `go mod`
- The `ros` binary is a multi-call binary (behavior changes based on argv[0])
- System Docker is a separate pre-built binary, not built from this repo's vendor tree
- Network configuration (`pkg/netconf/`) operates at the Linux netlink level - test on real/virtual hardware
