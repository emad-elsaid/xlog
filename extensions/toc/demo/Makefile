OUT = static

.PHONY: all
all: $(OUT)/wasm_exec.js $(OUT)/main.wasm

$(OUT)/wasm_exec.js:
	@mkdir -p $(OUT)
	cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" $@

$(OUT)/main.wasm: $(wildcard *.go)
	@mkdir -p $(OUT)
	GOOS=js GOARCH=wasm go build -o $@
