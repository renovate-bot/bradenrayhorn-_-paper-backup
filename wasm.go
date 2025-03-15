//go:build js && wasm

package main

import (
	"errors"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("paperBackup", js.FuncOf(paperBackup))
	js.Global().Set("paperBackupDecodeQR", js.FuncOf(paperBackupDecodeQR))
	js.Global().Set("paperBackupDecodeText", js.FuncOf(paperBackupDecodeText))

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

	array := js.Global().Get("Uint8Array").New(len(qr))
	js.CopyBytesToJS(array, qr)
	obj["qr"] = array

	return js.ValueOf(obj)
}

func paperBackupDecodeQR(this js.Value, args []js.Value) any {
	key := args[0].String()
	input := args[1].String()

	data, err := DecodeBackupFromQR(input, key)
	if err != nil {
		return makeJsError(err)
	}

	array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(array, data)
	return array
}

func paperBackupDecodeText(this js.Value, args []js.Value) any {
	key := args[0].String()
	input := args[1].String()

	data, err := DecodeBackupFromText(input, key)
	if err != nil {
		return makeJsError(err)
	}

	array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(array, data)
	return array
}

func makeJsError(err error) js.Value {
	return js.Global().Get("Error").New(err.Error())
}
