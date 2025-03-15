package encode

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"strings"
)

const maxLineSize = 80

func ToText(data []byte) string {
	formatted := ""

	// convert to hexr format
	hexr := string(hexrEncode(data))

	// split into lines of length `maxLineSize`
	for i := 0; i < len(hexr); i += maxLineSize {
		secretLine := hexr[i:]
		if i+maxLineSize <= len(hexr) {
			secretLine = hexr[i : i+maxLineSize]
		}

		// create checksum for line
		lineChecksum := string(hexrEncode(createChecksum(secretLine)))
		formatted += addWhitespaceToLine(secretLine + lineChecksum)
	}

	// create final checksum
	shareChecksum := hexrEncode(createChecksum(hexr))
	formatted += addWhitespaceToLine(string(shareChecksum))

	return formatted
}

func FromText(data string) ([]byte, error) {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid format. requires at least two lines")
	}

	secret := ""
	for j, line := range lines[:len(lines)-1] {
		stripped := strings.ReplaceAll(line, " ", "")
		if len(stripped) < 10 {
			return nil, fmt.Errorf("invalid format. at line %d", j+1)
		}

		foundChecksum := stripped[len(stripped)-8:]
		foundSecret := stripped[:len(stripped)-8]
		secret += foundSecret

		// validate checksum for line
		decodedChecksum, err := hexrDecode([]byte(foundChecksum))
		if err != nil {
			return nil, err
		}
		if !verifyChecksum(foundSecret, decodedChecksum) {
			return nil, fmt.Errorf("checksum failed. at line %d", j+1)
		}
	}

	// validate share checksum
	shareChecksum := strings.TrimSpace(strings.ReplaceAll(lines[len(lines)-1], " ", ""))
	decodedChecksum, err := hexrDecode([]byte(shareChecksum))
	if err != nil {
		return nil, err
	}
	if !verifyChecksum(secret, decodedChecksum) {
		return nil, fmt.Errorf("file checksum failed")
	}

	// decode line secret and add to full secret
	decodedSecret, err := hexrDecode([]byte(secret))
	if err != nil {
		return nil, err
	}

	return decodedSecret, nil
}

func createChecksum(secret string) []byte {
	checksum := crc32.ChecksumIEEE([]byte(secret))

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, checksum)
	return b
}

func verifyChecksum(secret string, checksum []byte) bool {
	if len(checksum) != 4 {
		return false
	}
	secretChecksum := crc32.ChecksumIEEE([]byte(secret))

	return secretChecksum == binary.LittleEndian.Uint32([]byte(checksum))
}

func addWhitespaceToLine(data string) string {
	for i := 4; i < len(data); i += 5 {
		data = data[:i] + " " + data[i:]
	}

	return data + "\n"
}
