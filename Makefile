.PHONY: build run clean

run:
	@go run cmd/bot/main.go

build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/bot-ctl cmd/bot/main.go
	@echo "[OK] Bot was built!"

clean:
	@rm -rf ./build
	@echo "[OK] Build directory cleaned!"
