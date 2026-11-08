package auth

import "testing"

func TestHashTokenStable(t *testing.T) {
	a := HashToken("abc")
	b := HashToken("abc")
	c := HashToken("xyz")
	if a != b {
		t.Fatal("hash should be stable")
	}
	if a == c {
		t.Fatal("different inputs should differ")
	}
	if len(a) != 64 {
		t.Fatalf("expected sha256 hex length 64, got %d", len(a))
	}
}
