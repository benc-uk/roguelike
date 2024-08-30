.PHONY: build watch serve editor
EDITOR_DEPLOY_BASE ?= wasm-dungeon/sprite-editor/

build-bin:
	env GOOS=linux GOARCH=amd64 go build -o bin/dungeon main.go
	env GOOS=windows GOARCH=amd64 go build -o bin/dungeon.exe main.go

build-wasm:
	env GOOS=js GOARCH=wasm go build -o web/main.wasm main.go
	rm -rf web/assets
	cp -r assets/ web/

watch:
	air -c .air.toml

serve:
	npx vite --port 3000 web

serve-editor:
	npx vite --port 8000 ./sprite-editor

build-editor:
	npx vite build ./sprite-editor --target esnext --base $(EDITOR_DEPLOY_BASE)

clean:
	rm -rf bin/ web/main.wasm site/ sprite-editor/dist/ web/assets site/

build-site: clean build-wasm build-editor
	mkdir -p site/sprite-editor
	cp -r ./sprite-editor/dist/* ./site/sprite-editor
	cp -r ./web/* ./site
