.DEFAULT_GOAL := help

help: ## This help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

serve: ## Serve the sprite editor
	npx vite --port 8001 ./

build: ## Bundle the sprite editor web app
	npx vite build ./ --target esnext --base $(EDITOR_DEPLOY_BASE)
