//go:build js && wasm

package main

import (
	"errors"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("paperBackup", js.FuncOf(paperBackup))

	<-c
}

func paperBackup(this js.Value, args []js.Value) any {
	dataUInt8Array := args[0]
	passphrase := args[1].String()

	if dataUInt8Array.Type() != js.TypeObject || dataUInt8Array.Get("constructor").Get("name").String() != "Uint8Array" {
		return makeJsError(errors.New("data must be uint8array"))
	}

	length := dataUInt8Array.Length()
	data := make([]byte, length)
	js.CopyBytesToGo(data, dataUInt8Array)

	text, qr, err := EncodeBackup([]byte(data), passphrase)
	if err != nil {
		return makeJsError(err)
	}

	// assemble response
	obj := map[string]any{}
	obj["text"] = text
	obj["qr"] = qr

	return js.ValueOf(obj)
}

func makeJsError(err error) js.Value {
	return js.Global().Get("Error").New(err.Error())
}
