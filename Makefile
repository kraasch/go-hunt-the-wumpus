
run:
	# rm -rf ~/.go-hunt-the-wumpus/
	go run ./cmd/wumpus.go

test:
	go test ./build/

build_and_install:
	rm -f ~/.local/bin/wumpus
	# builds a binary without dependencies.
	CGO_ENABLED=0 go build -a -tags netgo,osusergo -ldflags "-extldflags '-static' -s -w" -o ~/.local/bin/wumpus ./cmd

test_coverage:
	go test ./... -coverprofile=cover.temp
	go tool cover -html=cover.temp

.PHONY: build
build:
	rm -rf ./build/
	mkdir -p ./build/
	# builds a binary with dependencies.
	go build -o ./build/go-hunt-the-wumpus -gcflags -m=2 ./cmd/ 
