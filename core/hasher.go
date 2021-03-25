package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"strings"
	"time"
)

func HashBytes(content []byte) string {
	hasher := sha1.New()
	hasher.Write(content)
	sha1Bytes := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return string(sha1Bytes)
}

func TimedHash(content string) string {
	return HashBytes([]byte(time.Now().String() + content))
}

func CompressBytes(content []byte) []byte {
	var buffer bytes.Buffer

	writer := zlib.NewWriter(&buffer)
	writer.Write(content)
	writer.Close()

	return buffer.Bytes()
}

func DecompressBytes(content []byte) (string, error) {
	buffer := bytes.NewBuffer(content)
	sb := new(strings.Builder)

	reader, err := zlib.NewReader(buffer)
	if err != nil {
		return "", err
	}
	io.Copy(sb, reader)

	return sb.String(), nil
}
