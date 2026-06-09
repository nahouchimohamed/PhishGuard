package campaign

import (
	"context"
	"errors"
	"testing"

	"github.com/phishguard/phishguard/internal/store"
)

const (
	testOrgID      = "11111111-1111-1111-1111-111111111111"
	testAdminID    = "22222222-2222-2222-2222-222222222222"
	testTemplateID = "33333333-3333-3333-3333-333333333333"
	testUser1      = "44444444-4444-4444-4444-444444444441"
	testUser2      = "44444444-4444-4444-4444-444444444442"
)

func TestLaunchValidation(t *testing.T) {
	st := newFakeStore()
	st.templates[testTemplateID] = store.Template{ID: testTemplateID, OrgID: testOrgID}
	svc := newService(st, &fakePublisher{}, newFakeCache(), "http://localhost")
	ctx := context.Background()

	cases := []LaunchInput{
		{OrgID: testOrgID, CreatedBy: testAdminID, Name: "", TemplateID: testTemplateID, RecipientUserIDs: []string{testUser1}},
		{OrgID: testOrgID, CreatedBy: testAdminID, Name: "Camp", TemplateID: "", RecipientUserIDs: []string{testUser1}},
		{OrgID: testOrgID, CreatedBy: testAdminID, Name: "Camp", TemplateID: testTemplateID, RecipientUserIDs: nil},
		{OrgID: testOrgID, CreatedBy: testAdminID, Name: "Camp", TemplateID: "not-uuid", RecipientUserIDs: []string{testUser1}},
		{OrgID: testOrgID, CreatedBy: testAdminID, Name: "Camp", TemplateID: testTemplateID, RecipientUserIDs: []string{"bad"}},
	}
	for _, in := range cases {
		_, err := svc.Launch(ctx, in)
		if !errors.Is(err, store.ErrInvalidInput) {
			t.Fatalf("input %#v: got %v want ErrInvalidInput", in, err)
		}
	}
}

func TestLaunchDedupesRecipients(t *testing.T) {
	st := newFakeStore()
	st.templates[testTemplateID] = store.Template{ID: testTemplateID, OrgID: testOrgID}
	pub := &fakePublisher{}
	svc := newService(st, pub, newFakeCache(), "http://localhost")
	camp, err := svc.Launch(context.Background(), LaunchInput{
		OrgID:            testOrgID,
		CreatedBy:        testAdminID,
		Name:             "Phish",
		TemplateID:       testTemplateID,
		RecipientUserIDs: []string{testUser1, testUser1, testUser2},
	})
	if err != nil {
		t.Fatal(err)
	}
	if camp.Status != store.StatusActive {
		t.Fatalf("status=%q", camp.Status)
	}
	if len(pub.publishedIDs()) != 2 {
		t.Fatalf("published=%v want 2 unique", pub.publishedIDs())
	}
}

func TestLaunchPartialPublishStaysLaunching(t *testing.T) {
	st := newFakeStore()
	st.templates[testTemplateID] = store.Template{ID: testTemplateID, OrgID: testOrgID}
	pub := &fakePublisher{
		failOn: map[string]error{},
	}
	// Fail the second recipient after create assigns predictable IDs.
	pub.failOn["camp-1-r2"] = errors.New("sqs down")

	svc := newService(st, pub, newFakeCache(), "http://localhost")
	camp, err := svc.Launch(context.Background(), LaunchInput{
		OrgID:            testOrgID,
		CreatedBy:        testAdminID,
		Name:             "Phish",
		TemplateID:       testTemplateID,
		RecipientUserIDs: []string{testUser1, testUser2},
	})
	if err == nil {
		t.Fatal("expected publish incomplete error")
	}
	if camp.ID == "" {
		t.Fatal("campaign should still be returned")
	}
	got, err := st.GetCampaign(context.Background(), camp.ID, testOrgID)
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
