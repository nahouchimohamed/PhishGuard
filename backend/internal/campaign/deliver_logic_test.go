package campaign

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/phishguard/phishguard/internal/store"
)

func TestDeliverIdempotent(t *testing.T) {
	st := newFakeStore()
	st.users["u1"] = store.User{ID: "u1", OrgID: "o1"}
	st.recipients["r1"] = store.Recipient{ID: "r1", CampaignID: "c1", UserID: "u1"}
	cache := newFakeCache()
	svc := newService(st, &fakePublisher{}, cache, "http://localhost")

	ctx := context.Background()
	if err := svc.Deliver(ctx, "r1", "ignored-campaign-id"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Deliver(ctx, "r1", "ignored-campaign-id"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Deliver(ctx, "r1", "ignored-campaign-id"); err != nil {
		t.Fatal(err)
	}

	if len(st.events) != 1 {
		t.Fatalf("expected 1 delivered event, got %d (%v)", len(st.events), st.events)
	}
	// First delivery + two duplicates each invalidate.
	if cache.inv != 3 {
		t.Fatalf("expected 3 cache invalidations, got %d", cache.inv)
	}
}

func TestOpenReportIdempotent(t *testing.T) {
	st := newFakeStore()
	st.inbox["r1"] = store.InboxItem{RecipientID: "r1", CampaignID: "c1"}
	cache := newFakeCache()
	svc := newService(st, &fakePublisher{}, cache, "http://localhost")
	ctx := context.Background()

	if err := svc.Open(ctx, "r1", "u1", "o1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Open(ctx, "r1", "u1", "o1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Report(ctx, "r1", "u1", "o1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Report(ctx, "r1", "u1", "o1"); err != nil {
		t.Fatal(err)
	}
	if len(st.events) != 2 {
		t.Fatalf("events=%v want opened+reported once", st.events)
	}
}

func TestClickCountsFirstOnly(t *testing.T) {
	st := newFakeStore()
	st.users["u1"] = store.User{ID: "u1", OrgID: "o1"}
	st.templates["t1"] = store.Template{ID: "t1", OrgID: "o1", Name: "T", LessonText: "L"}
	st.campaigns["c1"] = store.Campaign{ID: "c1", OrgID: "o1", TemplateID: "t1"}
	st.recipients["r1"] = store.Recipient{ID: "r1", CampaignID: "c1", UserID: "u1", ClickToken: "tok"}
	st.byToken["tok"] = "r1"
	cache := newFakeCache()
	svc := newService(st, &fakePublisher{}, cache, "http://localhost")
	ctx := context.Background()

	if _, _, err := svc.Click(ctx, "tok"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.Click(ctx, "tok"); err != nil {
		t.Fatal(err)
	}
	if len(st.events) != 1 {
		t.Fatalf("events=%v want one clicked", st.events)
	}
	if cache.inv != 1 {
		t.Fatalf("inv=%d want 1", cache.inv)
	}
}

func TestRetryLaunchConflictWhenActive(t *testing.T) {
	st := newFakeStore()
	st.campaigns["c1"] = store.Campaign{ID: "c1", OrgID: "o1", Status: store.StatusActive}
	svc := newService(st, &fakePublisher{}, newFakeCache(), "http://localhost")
	_, err := svc.RetryLaunch(context.Background(), "c1", "o1")
	if !errors.Is(err, store.ErrConflict) {
		t.Fatalf("err=%v want ErrConflict", err)
	}
}

func TestPublishAllHonoursCancel(t *testing.T) {
	st := newFakeStore()
	pub := &fakePublisher{}
	svc := newService(st, pub, newFakeCache(), "http://localhost")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := svc.publishAll(ctx, "c1", []string{"r1", "r2"})
	if err == nil {
		t.Fatal("expected context error")
	}
	if len(pub.publishedIDs()) != 0 {
		t.Fatalf("published=%v", pub.publishedIDs())
	}
}

func TestRecipientLabelPending(t *testing.T) {
	now := time.Now()
	if got := RecipientLabel(nil, nil, nil, nil); got != "Pending" {
		t.Fatalf("got %q", got)
	}
	if got := RecipientLabel(&now, nil, nil, nil); got != "Delivered" {
		t.Fatalf("got %q", got)
	}
}
