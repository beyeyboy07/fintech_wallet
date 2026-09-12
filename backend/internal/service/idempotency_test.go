package service

import "testing"

func TestGenerateRequestHash(t *testing.T) {
	hash1 := GenerateRequestHash(
		3,
		80000,
		"Transfer test",
	)

	hash2 := GenerateRequestHash(
		3,
		80000,
		"Transfer test",
	)

	hash3 := GenerateRequestHash(
		3,
		50000,
		"Transfer test",
	)

	if hash1 != hash2 {
		t.Error("same request should produce same hash")
	}

	if hash1 == hash3 {
		t.Error("different request should produce different hash")
	}

	if len(hash1) != 64 {
		t.Errorf("expected SHA-256 hash length 64, got %d", len(hash1))
	}
}
