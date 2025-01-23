
run:
	# rm -rf ~/.go-hunt-the-wumpus/
	go run ./cmd/wumpus.go

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=cover.temp
	go tool cover -html=cover.temp

.PHONY: build
build:
	rm -rf ./build/
	mkdir -p ./build/
	go build \
		-o ./build/go-hunt-the-wumpus \
		-gcflags -m=2 \
		./cmd/ 
