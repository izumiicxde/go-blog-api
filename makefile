make:
	@go run ./cmd/
build:
	@go build -o /bin/main.exe .
run: build
	@/bin/main.exe