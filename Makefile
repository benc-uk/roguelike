EDITOR_DEPLOY_BASE ?= wasm-dungeon/sprite-editor/
.DEFAULT_GOAL := help

help: ## This help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build-bin: ## Build binaries for linux and windows
	env GOOS=linux GOARCH=amd64 go build -o bin/dungeon main.go
	env GOOS=windows GOARCH=amd64 go build -o bin/dungeon.exe main.go

build-wasm: ## Build as WASM for web
	env GOOS=js GOARCH=wasm go build -o web/main.wasm main.go
	rm -rf web/assets
	cp -r assets/ web/

watch: ## Watch for changes and rebuild
	air -c .air.toml

serve: ## Serve the web app
	npx vite --port 3000 ./web

serve-editor: ## Serve the sprite editor
	npx vite --port 8000 ./sprite-editor

build-editor: ## Build the sprite editor
	npx vite build ./sprite-editor --target esnext --base $(EDITOR_DEPLOY_BASE)

clean: ## Clean up
	rm -rf bin/ web/main.wasm site/ sprite-editor/dist/ web/assets site/

build-site: clean build-wasm build-editor ## Build the site for deployment
	mkdir -p site/sprite-editor
	cp -r ./sprite-editor/dist/* ./site/sprite-editor
	cp -r ./web/* ./site
