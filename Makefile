GO ?= go

.PHONY: all gen

all: gen
	$(GO) mod tidy
	$(GO) run main.go

gen: app/wire.go
	$(GO) generate ./...

app/wire_gen.go: app/wire.go
	$(GO) generate ./app
