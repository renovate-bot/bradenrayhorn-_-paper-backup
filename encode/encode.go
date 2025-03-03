package encode

import (
	"encoding/base64"
	"encoding/binary"
	"hash/crc32"
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

func ToQR(data []byte) string {
	return base64.RawStdEncoding.EncodeToString(data)
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
