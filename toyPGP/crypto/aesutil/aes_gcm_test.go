package aesutil

import (
	"bytes"
	"testing"
)

func TestAESGCM_RoundTrip(t *testing.T) {
	key, err := NewKey(32) // AES-256
	if err != nil {
		t.Fatal(err)
	}
	plaintext := []byte("Chicken Avocados")
	aad := []byte("demo-metadata")

	nonce, ct, err := EncryptGCM(key, plaintext, aad)
	if err != nil {
		t.Fatal(err)
	}
	pt, err := DecryptGCM(key, nonce, ct, aad)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt, plaintext) {
		t.Fatalf("round-trip mismatch: got %q want %q", pt, plaintext)
	}
}