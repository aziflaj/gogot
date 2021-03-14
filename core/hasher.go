package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

func HashBytes(content []byte) string {
	hasher := sha1.New()
	hasher.Write(content)
	sha1Bytes := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return string(sha1Bytes)
}

func CompressBytes(content []byte) []byte {
	var buffer bytes.Buffer

	writer := zlib.NewWriter(&buffer)
	writer.Write(content)
	writer.Close()

	return buffer.Bytes()
}

func DecompressBytes(content []byte) string {
	buffer := bytes.NewBuffer(content)
	sb := new(strings.Builder)

	reader, err := zlib.NewReader(buffer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	io.Copy(sb, reader)

	return sb.String()
}
