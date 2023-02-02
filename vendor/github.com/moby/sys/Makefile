PACKAGES ?= mountinfo mount sequential signal symlink
BINDIR ?= _build/bin
CROSS ?= linux/arm linux/arm64 linux/ppc64le linux/s390x \
	freebsd/amd64 openbsd/amd64 darwin/amd64 darwin/arm64 windows/amd64
SUDO ?= sudo -n
test test-local: RUN_VIA_SUDO = $(shell $(SUDO) true && echo -exec \"$(SUDO)\")

.PHONY: all
all: clean lint test cross

.PHONY: clean
clean:
	$(RM) mount/go-local.*

.PHONY: test
test: test-local
	set -eu; \
	for p in $(PACKAGES); do \
		(cd $$p; go test $(RUN_VIA_SUDO) -v .); \
	done

.PHONY: tidy
tidy:
	set -eu; \
		for p in $(PACKAGES); do \
		(cd $$p; go mod tidy); \
	done

# Test the mount module against the local mountinfo source code instead of the
# release specified in its go.mod. This allows catching regressions / breaking
# changes in mountinfo.
.PHONY: test-local
test-local: MOD = -modfile=go-local.mod
test-local:
	echo 'replace github.com/moby/sys/mountinfo => ../mountinfo' | cat mount/go.mod - > mount/go-local.mod
	# Run go mod tidy to make sure mountinfo dependency versions are met.
	cd mount && go mod tidy $(MOD) && go test $(MOD) $(RUN_VIA_SUDO) -v .
	$(RM) mount/go-local.*

.PHONY: lint
lint: $(BINDIR)/golangci-lint
	$(BINDIR)/golangci-lint version
	set -eu; \
	for p in $(PACKAGES); do \
		(cd $$p; \
		go mod download; \
		../$(BINDIR)/golangci-lint run); \
	done

$(BINDIR)/golangci-lint: $(BINDIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BINDIR) v1.45.2

$(BINDIR):
	mkdir -p $(BINDIR)

.PHONY: cross
cross:
	set -eu; \
	for osarch in $(CROSS); do \
		export GOOS=$${osarch%/*} GOARCH=$${osarch#*/}; \
		echo "# building for $$GOOS/$$GOARCH"; \
		for p in $(PACKAGES); do \
			(cd $$p; go build .); \
		done; \
	done
