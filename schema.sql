-- ============================================================
-- PERSONAL FINANCE GOAL API - DATABASE SCHEMA
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    user_id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username    VARCHAR(100) NOT NULL UNIQUE,
    email       VARCHAR(150) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'),
    is_deleted  BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE goals (
    goal_id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(user_id),
    name            VARCHAR(150) NOT NULL,
    description     TEXT,
    target_amount   NUMERIC(15, 2) NOT NULL,
    current_amount  NUMERIC(15, 2) NOT NULL DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'on_track' CHECK (status IN ('on_track', 'completed', 'cancelled')),
    deadline        DATE NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'),
    updated_at      TIMESTAMP,
    is_deleted      BOOLEAN NOT NULL DEFAULT false,
    CHECK (current_amount >= 0),
    CHECK (target_amount > 0)
);

CREATE TABLE saving_logs (
    log_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    goal_id     UUID NOT NULL REFERENCES goals(goal_id),
    user_id     UUID NOT NULL REFERENCES users(user_id),
    type        VARCHAR(10) NOT NULL CHECK (type IN ('deposit', 'withdraw')),
    amount      NUMERIC(15, 2) NOT NULL,
    note        TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'),
    CHECK (amount > 0)
);

CREATE INDEX idx_goals_user         ON goals(user_id) WHERE is_deleted = false;
CREATE INDEX idx_goals_deadline     ON goals(deadline ASC) WHERE is_deleted = false;
CREATE INDEX idx_saving_logs_goal   ON saving_logs(goal_id, created_at DESC);
CREATE INDEX idx_saving_logs_user   ON saving_logs(user_id);
