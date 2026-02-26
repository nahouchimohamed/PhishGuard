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

func TestParseBodyEmpty(t *testing.T) {
	_, err := ParseBody("")
	if err == nil {
		t.Fatal("expected error for empty body")
	}
}

func TestParseBodyBadJSON(t *testing.T) {
	_, err := ParseBody("{not-json")
	if err == nil {
		t.Fatal("expected error for bad json")
	}
}

func TestDeliverMessageJSONFieldNames(t *testing.T) {
	raw, err := json.Marshal(DeliverMessage{RecipientID: "r1", CampaignID: "c1"})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Fatalf("keys=%v", m)
	}
	if m["recipient_id"] != "r1" {
		t.Fatalf("recipient_id=%q", m["recipient_id"])
	}
	if m["campaign_id"] != "c1" {
		t.Fatalf("campaign_id=%q", m["campaign_id"])
	}
}
