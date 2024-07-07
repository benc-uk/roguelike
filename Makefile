.PHONY: build watch serve

build:
	env GOOS=js GOARCH=wasm go build -o web/main.wasm main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/dungeon main.go
	env GOOS=windows GOARCH=amd64 go build -o bin/dungeon.exe main.go
	cp -r assets/ web/

watch:
	air -c .air.toml

serve:
	vite --port 3000 web