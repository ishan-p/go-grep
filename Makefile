.PHONY = all clean
all: build

build: main.go
	@echo "Building binary..."
	go build -o go-grep

clean:
	@echo "Cleaning up..."
	rm grep
	go clean
