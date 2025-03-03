#!/bin/bash

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ui/src/wasm/

# Get base64 encoding of main.wasm
GOOS=js GOARCH=wasm go build -o main.wasm .
base64_wasm=$(cat main.wasm | base64 | tr -d '\n')

# Create the wasm_load.js file
cat > ui/src/wasm/load.js <<EOF
const wasmBinary = Uint8Array.from(atob('${base64_wasm}'), c => c.charCodeAt(0));
const go = new Go();
WebAssembly.instantiate(wasmBinary.buffer, go.importObject).then((result) => {
  go.run(result.instance);
});
EOF

