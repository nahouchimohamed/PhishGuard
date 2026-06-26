package campaign

import (
	"testing"

	"github.com/phishguard/phishguard/internal/store"
)

func TestNormalizeRecipientIDs(t *testing.T) {
	ids, err := normalizeRecipientIDs([]string{
		"11111111-1111-4111-8111-111111111111",
		"11111111-1111-4111-8111-111111111111",
		"22222222-2222-4222-8222-222222222222",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 {
		t.Fatalf("got %d", len(ids))
	}

	_, err = normalizeRecipientIDs(nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(store.ValidationError); !ok {
		t.Fatalf("want ValidationError, got %T", err)
	}
}
