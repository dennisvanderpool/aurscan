BIN     := aurscan
PREFIX  ?= /usr/local
PKGPATH := github.com/manticore-projects/aurscan/internal/version

# Version info from git (works in a checkout; release tarballs override via
# VERSION=... on the make command line, as the PKGBUILD does).
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short=12 HEAD 2>/dev/null)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X $(PKGPATH).Version=$(VERSION) \
	-X $(PKGPATH).Commit=$(COMMIT) \
	-X $(PKGPATH).Date=$(DATE)
GOFLAGS := -trimpath -ldflags="$(LDFLAGS)"

build:
	CGO_ENABLED=0 go build $(GOFLAGS) -o $(BIN) ./cmd/aurscan

version: build
	./$(BIN) --version

test:
	go vet ./...
	go test ./...

compress: build
	@command -v upx >/dev/null || { echo "upx not found (pacman -S upx)"; exit 1; }
	upx --best --lzma $(BIN)
	upx -t $(BIN)

release:
	@command -v upx >/dev/null || { echo "upx not found (pacman -S upx)"; exit 1; }
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o aurscan-linux-amd64 ./cmd/aurscan
	upx --best --lzma aurscan-linux-amd64 && upx -t aurscan-linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o aurscan-linux-arm64 ./cmd/aurscan

install: build
	install -Dm755 $(BIN) $(DESTDIR)$(PREFIX)/bin/$(BIN)
	ln -sf $(BIN) $(DESTDIR)$(PREFIX)/bin/syay
	ln -sf $(BIN) $(DESTDIR)$(PREFIX)/bin/aurscan-edit

clean:
	rm -f $(BIN) aurscan-linux-amd64 aurscan-linux-arm64

.PHONY: build version test compress release install clean
