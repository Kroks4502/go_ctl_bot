.PHONY: build

run:
	@go run ./...

build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/bot-ctlv ./...
	@echo "[OK] Bot was build!"
