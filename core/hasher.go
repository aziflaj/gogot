package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
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

func CompressBytes(content []byte) string {
	var buffer bytes.Buffer

	writer := zlib.NewWriter(&buffer)
	writer.Write(content)
	writer.Close()

	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func DecompressBytes(content string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}

	buffer := bytes.NewReader(data)
	reader, err := zlib.NewReader(buffer)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(result), nil

	// TODO: This should be done in bytes, not base64 strings!
	// buffer := bytes.NewBuffer(content)
	// sb := new(strings.Builder)

	// reader, err := zlib.NewReader(buffer)
	// if err != nil {
	// 	return "", err
	// }
	// io.Copy(sb, reader)

	// return sb.String(), nil
}
