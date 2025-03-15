#!/bin/bash

ts_nocheck_file() {
    local file_path="$1"
    
    cat /dev/stdin "$file_path" <<EOI > "${file_path}.tmp"
// @ts-nocheck
EOI

    mv "${file_path}.tmp" "$file_path"
}

set -e

EMSCRIPTEN_VERSION=4.0.5
EMSDK_DIR="emsdk"

BASE=$(pwd)
WORKING=$(pwd)/zxing

cd $WORKING

# these steps are mostly based on zxing-cpp WASM README.


# 1. install emsdk
if [ ! -d "$EMSDK_DIR" ]; then
	git clone https://github.com/emscripten-core/emsdk.git
fi

cd $WORKING/$EMSDK_DIR

./emsdk install $EMSCRIPTEN_VERSION
./emsdk activate $EMSCRIPTEN_VERSION
source ./emsdk_env.sh

cd $WORKING

# 2. init cmake project
mkdir -p $WORKING/build
cd $WORKING/build
emcmake cmake ../zxing-cpp/wrappers/wasm

# 3. apply patch
cd $WORKING/zxing-cpp
git restore .
git apply $WORKING/patches/BarcodeWriter.patch


# 4. build for wasm
cd $WORKING/build
export EMCC_CFLAGS="-O2 -Oz -s FILESYSTEM=0 -s SINGLE_FILE=1 -s EXPORT_ES6=1"
cmake --build .

# 5. copy to ui project
ts_nocheck_file "$WORKING/build/zxing.js"
cp $WORKING/build/zxing.js $BASE/ui/src/wasm/

# add types
cat > $BASE/ui/src/wasm/zxing.d.ts <<EOF
interface ZXingInstance {
  generateBarcodeFromBinary(data: Uint8Array, format: string, encoding: string, margin: number, width: number, height: number, eccLevel: number): WriteResult;
}

interface WriteResult {
  image: Uint8Array
  error: string
  delete: () => void
}

declare function ZXing(): Promise<ZXingInstance>;

export default ZXing;
EOF

# cleanup.. reset patch
cd $WORKING/zxing-cpp
git restore .

