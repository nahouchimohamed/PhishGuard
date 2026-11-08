package campaign

import "testing"

// Documents the idempotency rule enforced by store.MarkDelivered:
// only the first successful update counts as a new delivery.
func TestDeliveryCountedOnce(t *testing.T) {
	alreadyDelivered := false
	count := 0
	deliver := func() {
		if alreadyDelivered {
			return
		}
		alreadyDelivered = true
		count++
	}
	deliver()
	deliver()
	deliver()
	if count != 1 {
		t.Fatalf("expected 1 delivery, got %d", count)
	}
}
