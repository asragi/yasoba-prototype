GO ?= go

.PHONY: all gen

all: gen
	$(GO) run main.go

gen: app/wire_gen.go

app/wire_gen.go: app/wire.go
	$(GO) generate ./app
