package campaign

import (
	"context"
	"testing"

	"github.com/phishguard/phishguard/internal/store"
)

func TestDeliverIdempotent(t *testing.T) {
	st := newFakeStore()
	st.users["u1"] = store.User{ID: "u1", OrgID: "o1"}
	st.recipients["r1"] = store.Recipient{ID: "r1", CampaignID: "c1", UserID: "u1"}
	cache := newFakeCache()
	svc := newService(st, &fakePublisher{}, cache, "http://localhost")

	ctx := context.Background()
	if err := svc.Deliver(ctx, "r1", "c1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Deliver(ctx, "r1", "c1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Deliver(ctx, "r1", "c1"); err != nil {
		t.Fatal(err)
	}

	if len(st.events) != 1 {
		t.Fatalf("expected 1 delivered event, got %d (%v)", len(st.events), st.events)
	}
	if cache.inv != 1 {
		t.Fatalf("expected 1 cache invalidation, got %d", cache.inv)
	}
}
