package campaign

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/phishguard/phishguard/internal/queue"
	"github.com/phishguard/phishguard/internal/store"
)

type fakeStore struct {
	mu sync.Mutex

	templates   map[string]store.Template
	campaigns   map[string]store.Campaign
	recipients  map[string]store.Recipient
	users       map[string]store.User
	stats       map[string]store.CampaignStats
	inbox       map[string]store.InboxItem
	events      []string
	delivered   map[string]bool
	opened      map[string]bool
	reported    map[string]bool
	clicked     map[string]bool
	byToken     map[string]string // token -> recipientID
	templatesBy map[string]store.Template
	createSeq   int
	failCreate  error
	failMarkAct error
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		templates:  map[string]store.Template{},
		campaigns:  map[string]store.Campaign{},
		recipients: map[string]store.Recipient{},
		users:      map[string]store.User{},
		stats:      map[string]store.CampaignStats{},
		inbox:      map[string]store.InboxItem{},
		delivered:  map[string]bool{},
		opened:     map[string]bool{},
		reported:   map[string]bool{},
		clicked:    map[string]bool{},
		byToken:    map[string]string{},
	}
}

func (f *fakeStore) GetTemplate(ctx context.Context, id, orgID string) (store.Template, error) {
	t, ok := f.templates[id]
	if !ok || t.OrgID != orgID {
		return store.Template{}, store.ErrNotFound
	}
	return t, nil
}

func (f *fakeStore) CreateCampaignWithRecipients(ctx context.Context, orgID, createdBy, templateID, name string, userIDs, tokens []string) (store.Campaign, []string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failCreate != nil {
		return store.Campaign{}, nil, f.failCreate
	}
	f.createSeq++
	cid := "camp-" + strconv.Itoa(f.createSeq)
	c := store.Campaign{
		ID:         cid,
		OrgID:      orgID,
		CreatedBy:  createdBy,
		TemplateID: templateID,
		Name:       name,
		Status:     store.StatusLaunching,
	}
	f.campaigns[cid] = c
	ids := make([]string, len(userIDs))
	for i, uid := range userIDs {
		rid := cid + "-r" + strconv.Itoa(i+1)
		ids[i] = rid
		tok := ""
		if i < len(tokens) {
			tok = tokens[i]
		}
		f.recipients[rid] = store.Recipient{
			ID:         rid,
			CampaignID: cid,
			UserID:     uid,
			ClickToken: tok,
		}
		if tok != "" {
			f.byToken[tok] = rid
		}
	}
	return c, ids, nil
}

func (f *fakeStore) MarkCampaignActive(ctx context.Context, campaignID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failMarkAct != nil {
		return f.failMarkAct
	}
	c, ok := f.campaigns[campaignID]
	if !ok {
		return store.ErrNotFound
	}
	c.Status = store.StatusActive
	f.campaigns[campaignID] = c
	return nil
}

func (f *fakeStore) GetCampaign(ctx context.Context, id, orgID string) (store.Campaign, error) {
	c, ok := f.campaigns[id]
	if !ok || c.OrgID != orgID {
		return store.Campaign{}, store.ErrNotFound
	}
	return c, nil
}

func (f *fakeStore) ListUndeliveredRecipientIDs(ctx context.Context, campaignID string) ([]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var ids []string
	for id, r := range f.recipients {
		if r.CampaignID == campaignID && !f.delivered[id] {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (f *fakeStore) MarkDelivered(ctx context.Context, recipientID string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.recipients[recipientID]; !ok {
		return false, store.ErrNotFound
	}
	if f.delivered[recipientID] {
		return false, nil
	}
	f.delivered[recipientID] = true
	return true, nil
}

func (f *fakeStore) GetRecipient(ctx context.Context, recipientID string) (store.Recipient, error) {
	r, ok := f.recipients[recipientID]
	if !ok {
		return store.Recipient{}, store.ErrNotFound
	}
	return r, nil
}

func (f *fakeStore) GetUserByID(ctx context.Context, id string) (store.User, error) {
	u, ok := f.users[id]
	if !ok {
		return store.User{}, store.ErrNotFound
	}
	return u, nil
}

func (f *fakeStore) InsertEvent(ctx context.Context, orgID, campaignID, userID, kind, detail string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = append(f.events, kind+":"+detail)
	return nil
}

func (f *fakeStore) GetCampaignStats(ctx context.Context, campaignID string) (store.CampaignStats, error) {
	st, ok := f.stats[campaignID]
	if !ok {
		return store.CampaignStats{}, store.ErrNotFound
	}
	return st, nil
}

func (f *fakeStore) GetInboxItem(ctx context.Context, recipientID, userID string) (store.InboxItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	item, ok := f.inbox[recipientID]
	if !ok {
		return store.InboxItem{}, store.ErrNotFound
	}
	if f.opened[recipientID] && item.OpenedAt == nil {
		now := time.Now().UTC()
		item.OpenedAt = &now
		f.inbox[recipientID] = item
	}
	if f.reported[recipientID] && item.ReportedAt == nil {
		now := time.Now().UTC()
		item.ReportedAt = &now
		f.inbox[recipientID] = item
	}
	return item, nil
}

func (f *fakeStore) MarkReported(ctx context.Context, recipientID, userID string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.inbox[recipientID]; !ok {
		if _, ok := f.recipients[recipientID]; !ok {
			return false, store.ErrNotFound
		}
	}
	if f.reported[recipientID] {
		return false, nil
	}
	f.reported[recipientID] = true
	f.opened[recipientID] = true
	return true, nil
}

func (f *fakeStore) MarkOpened(ctx context.Context, recipientID, userID string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.inbox[recipientID]; !ok {
		if _, ok := f.recipients[recipientID]; !ok {
			return false, store.ErrNotFound
		}
	}
	if f.opened[recipientID] {
		return false, nil
	}
	f.opened[recipientID] = true
	return true, nil
}

func (f *fakeStore) MarkClickedByToken(ctx context.Context, token string) (store.Recipient, store.Template, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	rid, ok := f.byToken[token]
	if !ok {
		return store.Recipient{}, store.Template{}, false, store.ErrNotFound
	}
	r, ok := f.recipients[rid]
	if !ok {
		return store.Recipient{}, store.Template{}, false, store.ErrNotFound
	}
	c, ok := f.campaigns[r.CampaignID]
	if !ok {
		return store.Recipient{}, store.Template{}, false, store.ErrNotFound
	}
	tmpl, ok := f.templates[c.TemplateID]
	if !ok {
		return store.Recipient{}, store.Template{}, false, store.ErrNotFound
	}
	first := !f.clicked[rid]
	f.clicked[rid] = true
	f.opened[rid] = true
	now := time.Now().UTC()
	if r.ClickedAt == nil {
		r.ClickedAt = &now
	}
	if r.OpenedAt == nil {
		r.OpenedAt = &now
	}
	f.recipients[rid] = r
	return r, tmpl, first, nil
}

type fakePublisher struct {
	mu       sync.Mutex
	msgs     []queue.DeliverMessage
	failOn   map[string]error // recipientID -> error
	failNext int              // fail after N successful publishes (0 = never by count)
	count    int
}

func (p *fakePublisher) PublishDeliver(ctx context.Context, msg queue.DeliverMessage) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.count++
	if p.failOn != nil {
		if err, ok := p.failOn[msg.RecipientID]; ok {
			return err
		}
	}
	if p.failNext > 0 && p.count > p.failNext {
		return errors.New("publish failed")
	}
	p.msgs = append(p.msgs, msg)
	return nil
}

func (p *fakePublisher) publishedIDs() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	ids := make([]string, len(p.msgs))
	for i, m := range p.msgs {
		ids[i] = m.RecipientID
	}
	return ids
}

type fakeCache struct {
	mu     sync.Mutex
	stats  map[string]store.CampaignStats
	getErr error
	setErr error
	hits   int
	sets   int
	inv    int
}

func newFakeCache() *fakeCache {
	return &fakeCache{stats: map[string]store.CampaignStats{}}
}

func (c *fakeCache) GetStats(ctx context.Context, campaignID string) (store.CampaignStats, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.getErr != nil {
		return store.CampaignStats{}, false, c.getErr
	}
	st, ok := c.stats[campaignID]
	if ok {
		c.hits++
	}
	return st, ok, nil
}

func (c *fakeCache) SetStats(ctx context.Context, campaignID string, st store.CampaignStats) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.setErr != nil {
		return c.setErr
	}
	c.sets++
	c.stats[campaignID] = st
	return nil
}

func (c *fakeCache) Invalidate(ctx context.Context, campaignID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.inv++
	delete(c.stats, campaignID)
	return nil
}
