package main

import (
	"fmt"

	"github.com/bradenrayhorn/paper-backup/compress"
	"github.com/bradenrayhorn/paper-backup/crypt"
	"github.com/bradenrayhorn/paper-backup/encode"
)

func EncodeBackup(data []byte, key string) (string, string, error) {
	compressed, err := compress.Compress(data)
	if err != nil {
		return "", "", fmt.Errorf("compressing: %w", err)
	}

	encrypted, err := crypt.Encrypt(key, compressed)
	if err != nil {
		return "", "", fmt.Errorf("encrypting: %w", err)
	}

	textEncoding := encode.ToText(encrypted)
	qrEncoding := encode.ToQR(encrypted)

	return textEncoding, qrEncoding, nil
}
