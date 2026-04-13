package campaign

import (
	"context"
	"errors"
	"testing"

	"github.com/phishguard/phishguard/internal/store"
)

func TestGetStatsHit(t *testing.T) {
	st := newFakeStore()
	cache := newFakeCache()
	want := store.CampaignStats{Total: 3, Delivered: 2}
	cache.stats["c1"] = want
	svc := newService(st, &fakePublisher{}, cache, "http://localhost")

	got, err := svc.GetStats(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
	if cache.hits != 1 {
		t.Fatalf("hits=%d", cache.hits)
	}
	if cache.sets != 0 {
		t.Fatalf("sets=%d want 0 on hit", cache.sets)
	}
}

func TestGetStatsMiss(t *testing.T) {
	st := newFakeStore()
	want := store.CampaignStats{Total: 5, Reported: 1}
	st.stats["c1"] = want
	cache := newFakeCache()
	svc := newService(st, &fakePublisher{}, cache, "http://localhost")

	got, err := svc.GetStats(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
	if cache.sets != 1 {
		t.Fatalf("sets=%d want 1", cache.sets)
	}
	if cache.stats["c1"] != want {
		t.Fatal("cache not populated on miss")
	}
}

func TestGetStatsCacheErrorFallsThrough(t *testing.T) {
	st := newFakeStore()
	want := store.CampaignStats{Total: 9}
	st.stats["c1"] = want
	cache := newFakeCache()
	cache.getErr = errors.New("redis down")
	svc := newService(st, &fakePublisher{}, cache, "http://localhost")

	got, err := svc.GetStats(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
