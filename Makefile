GO ?= go
GENERATED := app/wire_gen.go

.PHONY: run go-generate

run: $(GENERATED)
	$(GO) run main.go

go-generate: $(GENERATED)

$(GENERATED): app/wire.go
	$(GO) generate ./...
