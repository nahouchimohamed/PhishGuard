package queue

import (
	"encoding/json"
	"testing"
)

func TestParseBody(t *testing.T) {
	raw, err := json.Marshal(DeliverMessage{RecipientID: "r1", CampaignID: "c1"})
	if err != nil {
		t.Fatal(err)
	}
	msg, err := ParseBody(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	if msg.RecipientID != "r1" || msg.CampaignID != "c1" {
		t.Fatalf("unexpected %#v", msg)
	}
}
