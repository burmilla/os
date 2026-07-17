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

# BurmillaOS 3.x Plan — Align with Debian 13 (trixie)

## Goals

- Share the **kernel version** and **console binaries** with Debian 13 so security
  fixes flow from Debian with minimal BurmillaOS-side maintenance.
- Existing use cases (standalone Docker nodes, Docker Swarm) keep working;
  existing installations must be upgradeable with `ros os upgrade`.
- Minimize long-term maintenance of both code and installed systems. Bigger
  breaking-internals changes are acceptable in 3.x when they serve those goals
  (the libcompose -> compose-go migration already merged here is one of them).

## Current state found in review (2.0.x baseline)

| Component | Current (2.0.x) | Problem | 3.x target |
|---|---|---|---|
| Kernel | 5.10.248-burmilla, built from kernel.org sources in `os-kernel` with own config/patches | 5.10 LTS EOL end of 2026; all config/firmware maintenance on BurmillaOS | Debian 13 kernel (6.12 LTS), maintained by Debian security team |
| Console | `debian:bullseye-slim` (Debian 11) in `images/02-console` | Debian 11 LTS ends Aug 2026 | `debian:trixie-slim` (Debian 13) |
| os-base | Buildroot 2023.02.10 glibc userland (busybox, dhcpcd, e2fsprogs, xfsprogs, cryptsetup, lvm2, mdadm, rsyslog, logrotate, eudev, wpa_supplicant, ntpd...) | Buildroot 2023.02 LTS is EOL; every CVE requires a manual Buildroot bump + rebuild | Debian 13-based rootfs (binaries shared with Debian) |
| os-initrd-base | Buildroot 2023.02.10 **static uClibc busybox** initrd, kernel headers pinned to 5.10 | Same EOL problem; headers pin blocks 6.12 | Rebuild against 6.12 headers; keep static busybox (smallest) or use Debian `busybox-static` |
| Build environment | `Dockerfile.dapper` FROM `ubuntu:bionic` (18.04, EOL) | EOL base, old syslinux/xorriso toolchain | `debian:trixie` build image |
| Go toolchain | Go 1.19.5, `GO111MODULE=off`, `trash` + checked-in vendor | Go 1.19 EOL; `golang.org/x/crypto`, `x/net`, `x/sys` pinned to go1.15-era commits (known CVEs in x/crypto SSH, e.g. Terrapin) | Debian 13's Go (1.24.x); migrate to Go modules + `go mod vendor` |
| System Docker | 17.06.107 fork (pre-built from `burmilla/os-system-docker`) | cgroup v1 only; ancient runc; blocks cgroup v2-only futures | Decide: keep on cgroup v1 for 3.0 or upgrade (see open questions) |
| User Docker | `os-services` v2.0.x branch, engines up to 29.1.5 (`DOCKER_MIN_API_VERSION=1.24` workaround for old System Docker API) | Workarounds pile up because System Docker is old | Keep current engine cadence; new `v3.0.x` branch |

## Plan steps

### 1. os-kernel: consume the Debian 13 kernel

The `burmilla/os` build only needs a `kernel.tar.gz` with this layout
(see `scripts/layout-kernel`): `boot/vmlinuz-*`, `lib/modules/<ver>/`,
`lib/firmware/`. That contract makes it possible to stop compiling kernels:

1. Create a `v6.12.x-debian` branch in `os-kernel`.
2. Preferred approach (least maintenance): **repackage Debian binary packages**
   instead of building from source. Download `linux-image-<ver>-amd64` /
   `linux-image-<ver>-arm64` (and matching `firmware-linux-free`/
   `firmware-*` packages from `non-free-firmware`) for the current Debian 13
   point release, extract, and re-tar into the `kernel.tar.gz` contract above.
   A new kernel release then becomes "bump Debian package version + repackage",
   and CVE handling is entirely Debian's.
   - Fallback approach if the binary kernel is missing something: build from the
     Debian `linux` *source* package (which carries Debian's patches and config)
     with minimal config overrides, still tracking their version.
3. Verify Debian's kernel config against BurmillaOS needs before committing to
   the binary route. Known requirements to check: overlayfs, br_netfilter and
   friends for Docker/Swarm (vxlan, ipvs), BPF (enabled for 2.x in
   `76d5ad2`), squashfs, iscsi, zfs-compatible build options, and — critical —
   **cgroup v1 controllers** (`CONFIG_MEMCG_V1` etc. are no longer default-on in
   6.12) as long as System Docker 17.06 is kept (see step 6).
4. Modules that BurmillaOS loads from initrd must exist in Debian's (heavily
   modular) config; update `modules/x86/modules.list` + `modules-extra.list`
   equivalents or drop the check.
5. Update `os-services` `kernel-headers`, `kernel-headers-system-docker`,
   `kernel-extras` and `zfs` services to publish images matching the Debian
   kernel version tag (`burmilla/os-headers:<debian-kernel-version>`); for the
   Debian kernel these can simply install Debian's `linux-headers-*` package.

### 2. os-base: replace Buildroot userland with Debian 13

`os-base` provides the rootfs used by `images/01-base` and all `02-*` system
containers (acpid, bootstrap, logrotate, syslog) plus tools mounted into
system containers. To share binaries with Debian 13:

1. Rebuild `os-base` as a minimal Debian 13 rootfs (debootstrap/mmdebstrap
   `--variant=minbase`) containing the same tool set the Buildroot config
   provides today: busybox or coreutils+bash, `dhcpcd`, `e2fsprogs`,
   `xfsprogs`, `dosfstools`, `parted`, `cryptsetup`, `lvm2`, `mdadm`, `kmod`,
   `udev` (eudev today - switch to Debian's udev without systemd running, it
   works standalone), `wpa_supplicant` + `wireless-tools`, `rsyslog`,
   `logrotate`, `ntpd` (Debian: `ntpsec` or switch to `chrony`), `kexec-tools`,
   `open-iscsi`, `ipset`/`iptables`, CA certificates.
2. Watch image size: Buildroot rootfs is much smaller than even minbase Debian.
   Mitigations: `--variant=minbase`, `dpkg --path-exclude` for docs/locales,
   busybox for shell utilities. Some growth is acceptable — the payoff is that
   `apt` security updates can be consumed by simply rebuilding.
3. Keep the existing artifact contract: `os-base_<arch>.tar.xz` consumed by
   `Dockerfile.dapper` (OS_BASE_URL) and `images/00-rootfs`.
4. `images/01-base` currently uses busybox `adduser`/`addgroup` syntax and
   dhcpcd hook paths — adjust for Debian equivalents when the rootfs switches.
5. Fix `/etc/os-release` / `/etc/lsb-release` branding (BurmillaOS identity is
   currently overlaid on top; keep that behavior).

### 3. os-initrd-base: minimal update

The initrd base is a static uClibc busybox — it has no Debian equivalent
benefit (a static busybox is a static busybox) and is the smallest-risk piece:

1. Option A (minimal work): bump Buildroot to a current LTS only to refresh the
   static busybox/uclibc, and change `BR2_KERNEL_HEADERS_5_10` to 6.12-compatible
   headers.
2. Option B (fewer repos to maintain): drop the Buildroot build and use Debian
   13's `busybox-static` binary + `ca-certificates.crt` asset. Verify all
   busybox applets used by `scripts/layout-initrd` and early `ros init` exist in
   Debian's busybox build (Debian disables some applets!).
3. The real initrd content (ros binary, kernel modules, firmware) comes from
   the main repo build, so this repo stays tiny either way.

### 4. images/02-console: Debian 13 console

1. Switch `FROM debian:bullseye-slim` to `FROM debian:trixie-slim`.
2. Revisit `update-alternatives --set iptables ... iptables-legacy`: trixie
   defaults to nftables backend. User Docker >= 20.10 works with iptables-nft,
   but rules must be consistent between console tooling, `os-base` network
   service (which mounts `/usr/bin/iptables` from os-base into the network
   container) and Docker itself. Decide legacy vs nft **once, globally** —
   mixing backends breaks Swarm networking. Keeping legacy is the
   compatibility-safe 3.0 choice; nft migration can be its own later step.
3. Check trixie package renames/removals in the console package list
   (`net-tools`, `nvi`, `open-iscsi`, `apparmor` are still present in trixie;
   verify at build time).
4. Trixie images are merged-/usr; the console-init bind-mount logic in
   `cmd/control/console_init.go` should be re-tested against merged-usr paths.

### 5. Main repo build modernization

1. `Dockerfile.dapper`: move FROM `ubuntu:bionic` to `debian:trixie`. The
   apt package list maps almost 1:1 (`isolinux`, `syslinux-common`, `xorriso`,
   `genisoimage`→`xorriso`/`mkisofs` compat, `qemu-kvm`→`qemu-system-x86`).
   Drop the gccgo remnants. Keep `KERNEL_URL`/`OS_BASE_URL` override args.
2. Go toolchain: bump `GO_VERSION` to Debian 13's Go (1.24.x line).
   - Migrate from `trash`/GOPATH to **Go modules with a checked-in `vendor/`**
     (`go mod vendor`). `go get` no longer works in GOPATH mode since Go 1.22,
     and `trash` is unmaintained; modules are the only sustainable path.
   - When modules land, update `CLAUDE.md` build notes and remove
     `GO111MODULE=off` from `Dockerfile.dapper`/`Makefile`/scripts.
   - Update `golang.org/x/crypto` (SSH host key / cloud-init key handling) and
     `x/net`, `x/sys` off the go1.15 release branches — security relevant.
   - Keep the burmilla forks (netlink, candiedyaml, cli-1, docker, containerd,
     runc) initially; replace opportunistically only when a maintained upstream
     equivalent is verified to work.
3. `os-config.tpl.yml` + `Dockerfile.dapper` version pointers for 3.x:
   - `repositories.core.url` -> `${OS_SERVICES_REPO}/v3.0.x` (new branch, step 7)
   - `OS_RELEASES_YML` -> `https://raw.githubusercontent.com/burmilla/releases/v3.0.x`
   - `KERNEL_VERSION`/`KERNEL_URL` -> Debian-based kernel artifact from step 1
   - `OS_BASE_URL`/`OS_INITRD_BASE_URL` -> new Debian-based releases from steps 2-3
4. CI: `.github/workflows/*` still use `actions/checkout@v2` and bare
   `ubuntu-latest`; bump action versions while touching the files.

### 6. System Docker decision (biggest open question)

System Docker 17.06.107 is the highest-risk legacy piece. Facts found in review:

- `pkg/dfs/scratch.go` mounts **cgroup v1** hierarchies parsed from
  `/proc/cgroups`; there is no cgroup v2 mount path. `console_init.go` also
  mounts a v1 `name=systemd` hierarchy.
- Docker 17.06 cannot run on a cgroup v2-only kernel.
- Linux 6.12 made several v1 controllers optional (`CONFIG_MEMCG_V1=n` by
  default); Debian's config must be checked (step 1.3), and even where v1 still
  works it is on its way out kernel-side and Docker-side.
- `os-services` already carries `DOCKER_MIN_API_VERSION=1.24` downgrade hacks so
  modern user-Docker CLIs can talk to the old System Docker.

Recommended 3.x scope:

1. **3.0 ships with System Docker 17.06 kept on cgroup v1** *if and only if*
   the Debian 6.12 kernel (or a one-line config override in the fallback
   source build) still provides the v1 controllers System Docker needs
   (cpu, cpuacct, cpuset, memory, devices, freezer, blkio, pids). Boot with
   `cgroup_memory=1`-style cmdline entries as required. This keeps upgrade risk
   for existing installations near zero.
2. Start a parallel `os-system-docker` upgrade track (modern moby or plain
   containerd+nerdctl) targeting 3.1+: add cgroup2 mounting support to
   `pkg/dfs/scratch.go`, remove the API-version downgrade hacks, and drop the
   17.06-era vendored client pins in this repo. This is the single change that
   would retire the most forked-code maintenance (burmilla/docker,
   burmilla/containerd, burmilla/runc forks all exist because of 17.06).
3. If Debian's binary kernel turns out to be v1-incapable, the decision
   flips: System Docker upgrade becomes a 3.0 blocker (or the kernel falls back
   to source-build-with-config-override, step 1.2 fallback).

### 7. os-services + releases branches

1. Branch `v2.0.x` -> `v3.0.x` in `burmilla/os-services`; keep service set
   (index.yml) unchanged so existing `services_include` configs keep working.
2. Rebuild kernel-dependent service images (`kernel-headers`, `kernel-extras`,
   `zfs`) for the Debian kernel version (see step 1.5). For a Debian kernel
   these become thin wrappers around Debian's own `linux-headers`/`zfs-dkms`
   packages — less custom build machinery.
3. Add `v3.0.x` branch + `releases.yml` in `burmilla/releases` so
   `ros os upgrade` can see 3.x (the `upgrade.url` in step 5.3 points there).
4. Console list: 2.x already dropped non-Debian consoles, nothing to remove.

### 8. Upgrade path 2.x -> 3.x (must-not-break)

1. `ros os upgrade` runs the *new* version's os image on the *old* system:
   verify the 3.x upgrade container still works against the 2.x System Docker
   API (this is why the old API-version pin in `pkg/compose`/docker client
   matters; test explicitly).
2. Kernel jump 5.10 -> 6.12: existing `modprobe`/module-name assumptions,
   renamed modules, and removed drivers should be checked against the
   hardware/VM targets we officially support (VMware, Hyper-V, KVM/Proxmox,
   Azure, bare metal amd64/arm64, Raspberry Pi 64).
3. Boot stack: syslinux/isolinux from the trixie build env must still produce
   ISOs bootable on existing BIOS + UEFI installs, and `ros os upgrade` writes
   the new kernel/initrd into the existing boot partition — test on a disk
   installed with 2.0.x, including `system-docker.json`/`docker` data survival.
4. Config compatibility: all existing `/var/lib/rancher/conf/cloud-config.yml`
   keys must keep parsing (the compose-go migration already maintains v1
   service-format compatibility via `pkg/libcompose` + `config/compat.go` —
   keep its tests green).
5. Document the jump: minimum supported upgrade base (recommend: only from
   2.0.x, not 1.x), and known-removed kernel drivers if any.

## Suggested execution order

1. Step 5.1-5.2 (build env + Go modules) — unblocks everything else and makes
   CI trustworthy for the rest.
2. Step 1 (Debian kernel) with the cgroup v1 verification (6.1) done first —
   its outcome decides System Docker scope.
3. Steps 2-4 (os-base, os-initrd-base, console) — independent of each other,
   can proceed in parallel.
4. Step 7 (branches) once artifacts exist.
5. Step 8 (upgrade testing) continuously, formal pass before 3.0.0-rc1.
6. Step 6.2 (System Docker replacement) as the headline 3.1 item unless 6.3
   forces it into 3.0.
