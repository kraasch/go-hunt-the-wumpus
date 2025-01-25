
VERSION=$(shell git describe --tags --long --dirty 2>/dev/null)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

ifeq ($(VERSION),)
	VERSION = UNKNOWN
endif

run:
	# rm -rf ~/.go-hunt-the-wumpus/
	go run ./cmd/wumpus.go

info:
	@echo version: $(VERSION)
	@echo branch:  $(BRANCH)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: committed
committed:
	@git diff --exit-code >/dev/null || (echo " ** COMMIT YOUR CHANGES FIRST **"; exit 1)

test:
	go test ./build/

build_and_install:
	rm -f ~/.local/bin/wumpus
	@# BUILDS A BINARY WITHOUT DEPENDENCIES.
	@# CGO‥    ⇒ do not link c libraries (eg libc).
	@# -a      ⇒ explicitly redo all.
	@# -tags‥  ⇒ do not link network or user libraries.
	@# -X‥     ⇒ set the version as a string.
	@# -static ⇒ make a static and make smaller.
	@#  -s -w  ⇒ tell linker strips symbol table and DWARF debug info.
	@# -o      ⇒ name of binary.
	@# .       ⇒ target.
	CGO_ENABLED=0 go build                 \
	-a                                     \
	-tags netgo,osusergo                   \
	-ldflags "-X main.version=$(VERSION)"  \
	-ldflags "-extldflags '-static' -s -w" \
	-o ~/.local/bin/wumpus                 \
	./cmd

test_coverage:
	go test ./... -coverprofile=cover.temp
	go tool cover -html=cover.temp

.PHONY: build
build: committed lint
	rm -rf ./build/
	mkdir -p ./build/
	# builds a binary with dependencies.
	go build -o ./build/go-hunt-the-wumpus -gcflags -m=2 ./cmd/ 
