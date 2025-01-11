
run:
	# rm -rf ~/.go-hunt-the-wumpus/
	go run ./cmd/stuff.go

test:
	go test ./...

.PHONY: build
build:
	rm -rf ./build/
	mkdir -p ./build/
	go build \
		-o ./build/go-hunt-the-wumpus \
		-gcflags -m=2 \
		./cmd/ 
