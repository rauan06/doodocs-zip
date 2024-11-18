format:
	gofumpt -l -w .
run:
	gofumpt -l -w .
	go run cmd/main.go
build:
	gofumpt -l -w .
	go build -o doodocs-zip .