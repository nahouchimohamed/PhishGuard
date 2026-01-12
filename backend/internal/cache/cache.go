package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/phishguard/phishguard/internal/store"
	"github.com/redis/go-redis/v9"
)

// Dashboard caches campaign stats for a few seconds.
// Mutations call Invalidate so the next read hits Postgres.
type Dashboard struct {
	rdb *redis.Client
	ttl time.Duration
}

func New(redisURL string) (*Dashboard, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return &Dashboard{rdb: rdb, ttl: 5 * time.Second}, nil
}

func (d *Dashboard) Close() error {
	return d.rdb.Close()
}

func (d *Dashboard) Ping(ctx context.Context) error {
	return d.rdb.Ping(ctx).Err()
}

func (d *Dashboard) key(campaignID string) string {
	return "campaign:stats:" + campaignID
}

func (d *Dashboard) GetStats(ctx context.Context, campaignID string) (store.CampaignStats, bool, error) {
	raw, err := d.rdb.Get(ctx, d.key(campaignID)).Bytes()
	if err == redis.Nil {
		return store.CampaignStats{}, false, nil
	}
	if err != nil {
		return store.CampaignStats{}, false, err
	}
	var st store.CampaignStats
	if err := json.Unmarshal(raw, &st); err != nil {
		return store.CampaignStats{}, false, err
	}
	return st, true, nil
}

func (d *Dashboard) SetStats(ctx context.Context, campaignID string, st store.CampaignStats) error {
	raw, err := json.Marshal(st)
	if err != nil {
		return err
	}
	return d.rdb.Set(ctx, d.key(campaignID), raw, d.ttl).Err()
}

func (d *Dashboard) Invalidate(ctx context.Context, campaignID string) error {
	return d.rdb.Del(ctx, d.key(campaignID)).Err()
}
