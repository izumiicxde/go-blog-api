build:
	@go build -o /bin/main.exe .
run: build
	@/bin/main.exe