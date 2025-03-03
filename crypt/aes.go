package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func Encrypt(passphrase string, data []byte) ([]byte, error) {
	salt := make([]byte, 48)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	encrypted := gcm.Seal(nil, nonce, data, nil)
	return append(append(salt, nonce...), encrypted...), nil
}

func Decrypt(passphrase string, encrypted []byte) ([]byte, error) {
	if len(encrypted) < 49 {
		return nil, fmt.Errorf("crypto data invalid")
	}

	salt := encrypted[:48]
	encrypted = encrypted[48:]

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(encrypted) < gcm.NonceSize() {
		return nil, fmt.Errorf("crypto data invalid")
	}
	nonce := encrypted[:gcm.NonceSize()]
	encrypted = encrypted[gcm.NonceSize():]

	return gcm.Open(nil, nonce, encrypted, nil)
}

func deriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(passphrase),
		salt,
		8,
		1024*1024,
		8,
		32,
	)
}
