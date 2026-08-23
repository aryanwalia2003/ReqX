.PHONY: dev dev-frontend build build-cli lint lint-go lint-frontend fmt fmt-go fmt-frontend test

## Run the desktop app in dev mode (requires the Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
dev:
	cd app && wails dev

## Run just the Vite dev server, no Go backend attached
dev-frontend:
	cd app/frontend && npm run dev

## Build the CLI binary. For the desktop app: cd app && wails build
build: build-cli

build-cli:
	go build -o reqx .

lint: lint-go lint-frontend

lint-go:
	golangci-lint run ./...

lint-frontend:
	cd app/frontend && npm run lint

fmt: fmt-go fmt-frontend

fmt-go:
	gofmt -l -w $(shell find . -name '*.go' -not -path './app/frontend/*')

fmt-frontend:
	cd app/frontend && npm run format

test:
	go test ./...
	cd app/frontend && npm run test
