-- PhishGuard schema. Applied once by Postgres on first boot via docker-entrypoint-initdb.d.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE orgs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id        UUID NOT NULL REFERENCES orgs(id),
    email         TEXT NOT NULL UNIQUE,
    display_name  TEXT NOT NULL,
    role          TEXT NOT NULL CHECK (role IN ('admin', 'employee')),
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX users_org_id_idx ON users(org_id);

CREATE TABLE sessions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX sessions_user_id_idx ON sessions(user_id);
CREATE INDEX sessions_expires_at_idx ON sessions(expires_at);

CREATE TABLE templates (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id       UUID NOT NULL REFERENCES orgs(id),
    name         TEXT NOT NULL,
    subject      TEXT NOT NULL,
    body_text    TEXT NOT NULL,
    lesson_text  TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX templates_org_id_idx ON templates(org_id);

CREATE TABLE campaigns (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id       UUID NOT NULL REFERENCES orgs(id),
    created_by   UUID NOT NULL REFERENCES users(id),
    template_id  UUID NOT NULL REFERENCES templates(id),
    name         TEXT NOT NULL,
    status       TEXT NOT NULL CHECK (status IN ('draft', 'launching', 'active')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    launched_at  TIMESTAMPTZ
);

CREATE INDEX campaigns_org_id_idx ON campaigns(org_id);

CREATE TABLE campaign_recipients (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id   UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    user_id       UUID NOT NULL REFERENCES users(id),
    click_token   TEXT NOT NULL UNIQUE,
    delivered_at  TIMESTAMPTZ,
    opened_at     TIMESTAMPTZ,
    clicked_at    TIMESTAMPTZ,
    reported_at   TIMESTAMPTZ,
    UNIQUE (campaign_id, user_id)
);

CREATE INDEX campaign_recipients_user_id_idx ON campaign_recipients(user_id);
CREATE INDEX campaign_recipients_campaign_id_idx ON campaign_recipients(campaign_id);

CREATE TABLE events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id       UUID NOT NULL REFERENCES orgs(id),
    campaign_id  UUID REFERENCES campaigns(id) ON DELETE SET NULL,
    user_id      UUID REFERENCES users(id) ON DELETE SET NULL,
    kind         TEXT NOT NULL,
    detail       TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX events_campaign_id_idx ON events(campaign_id);
CREATE INDEX events_created_at_idx ON events(created_at);
