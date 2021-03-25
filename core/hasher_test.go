package core

import (
	"testing"
	"time"
)

func TestCompressAndDecompressing(t *testing.T) {
	message := "Compress me"
	compressed := CompressBytes([]byte(message))
	decompressed, err := DecompressBytes(compressed)
	if err != nil {
		t.Errorf("Error arose: %v", err)
	}

	if decompressed != message {
		t.Error("Data loss!")
	}
}

func TestTimedHash(t *testing.T) {
	message := "Hash me"
	hash1 := TimedHash(message)
	time.Sleep(2 * time.Second)
	hash2 := TimedHash(message)
	if hash1 == hash2 {
		t.Error("Expected TimedHash() to be time dependent")
	}
}
