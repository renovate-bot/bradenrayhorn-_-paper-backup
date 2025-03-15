#!/bin/bash

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
export EMCC_CFLAGS="-O2 -Oz -s FILESYSTEM=0 -s SINGLE_FILE=1 -s EXPORT_ES6=1" # bundle wasm into js file
cmake --build .

# 5. copy to ui project
cp $WORKING/build/zxing.js $BASE/ui/src/wasm/

# cleanup.. reset patch
cd $WORKING/zxing-cpp
git restore .

