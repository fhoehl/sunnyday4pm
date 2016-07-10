build: makeicecream icecreamd

.PHONY: makeicecream
makeicecream:
	go build -o bin/makeicecream makeicecream/cli.go

.PHONY: icecreamd
icecreamd:
	go build -o bin/icecreamd icecreamd/server.go
