.PHONY: test testdata
.PRECIOUS: %.wasm

sdk.src.go = \
	$(wildcard sdk/go/*/*.go)

timecraft.src.go = \
	$(wildcard *.go) \
	$(wildcard cmd/*.go)

testdata.src.go = $(wildcard testdata/*.go)

%.wasm: %.go $(sdk.src.go)
	tinygo build -o $@ -target=wasi $<

timecraft: go.mod $(timecraft.src.go)
	go build -o timecraft

test:
	go test ./...

testdata: $(testdata.src.go:.go=.wasm)