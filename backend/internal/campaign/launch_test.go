package campaign

import (
	"context"
	"errors"
	"testing"

	"github.com/phishguard/phishguard/internal/store"
)

func TestLaunchValidation(t *testing.T) {
	st := newFakeStore()
	st.templates["t1"] = store.Template{ID: "t1", OrgID: "o1"}
	svc := newService(st, &fakePublisher{}, newFakeCache(), "http://localhost")
	ctx := context.Background()

	cases := []LaunchInput{
		{OrgID: "o1", CreatedBy: "a1", Name: "", TemplateID: "t1", RecipientUserIDs: []string{"u1"}},
		{OrgID: "o1", CreatedBy: "a1", Name: "Camp", TemplateID: "", RecipientUserIDs: []string{"u1"}},
		{OrgID: "o1", CreatedBy: "a1", Name: "Camp", TemplateID: "t1", RecipientUserIDs: nil},
	}
	for _, in := range cases {
		_, err := svc.Launch(ctx, in)
		if !errors.Is(err, store.ErrInvalidInput) {
			t.Fatalf("input %#v: got %v want ErrInvalidInput", in, err)
		}
	}
}

func TestLaunchPartialPublishStaysLaunching(t *testing.T) {
	st := newFakeStore()
	st.templates["t1"] = store.Template{ID: "t1", OrgID: "o1"}
	pub := &fakePublisher{
		failOn: map[string]error{},
	}
	// Fail the second recipient after create assigns predictable IDs.
	// CreateCampaignWithRecipients uses camp-1-r1, camp-1-r2.
	pub.failOn["camp-1-r2"] = errors.New("sqs down")

	svc := newService(st, pub, newFakeCache(), "http://localhost")
	camp, err := svc.Launch(context.Background(), LaunchInput{
		OrgID:            "o1",
		CreatedBy:        "a1",
		Name:             "Phish",
		TemplateID:       "t1",
		RecipientUserIDs: []string{"u1", "u2"},
	})
	if err == nil {
		t.Fatal("expected publish incomplete error")
	}
	if camp.ID == "" {
		t.Fatal("campaign should still be returned")
	}
	got, err := st.GetCampaign(context.Background(), camp.ID, "o1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != store.StatusLaunching {
		t.Fatalf("status=%q want launching", got.Status)
	}
}

func TestRetryLaunchUndeliveredOnly(t *testing.T) {
	st := newFakeStore()
	st.campaigns["c1"] = store.Campaign{ID: "c1", OrgID: "o1", Status: store.StatusLaunching}
	st.recipients["r1"] = store.Recipient{ID: "r1", CampaignID: "c1", UserID: "u1"}
	st.recipients["r2"] = store.Recipient{ID: "r2", CampaignID: "c1", UserID: "u2"}
	st.delivered["r1"] = true // already delivered; RetryLaunch should skip

	pub := &fakePublisher{}
	svc := newService(st, pub, newFakeCache(), "http://localhost")
	camp, err := svc.RetryLaunch(context.Background(), "c1", "o1")
	if err != nil {
		t.Fatal(err)
	}
	if camp.Status != store.StatusActive {
		t.Fatalf("status=%q want active", camp.Status)
	}
	ids := pub.publishedIDs()
	if len(ids) != 1 || ids[0] != "r2" {
		t.Fatalf("published %#v want only r2", ids)
	}
}
