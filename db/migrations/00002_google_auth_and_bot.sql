-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ALTER COLUMN password DROP NOT NULL;

ALTER TABLE users
    ADD COLUMN google_sub VARCHAR(255) UNIQUE,
    ADD COLUMN email VARCHAR(255),
    ADD COLUMN is_bot BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_users_google_sub ON users (google_sub);

-- Fixed UUID for the AI critic bot account (idempotent seed).
INSERT INTO users (id, created_at, updated_at, username, password, google_sub, email, is_bot)
VALUES (
    'a0000000-0000-4000-8000-000000000001',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    'vibe_critic',
    NULL,
    NULL,
    NULL,
    TRUE
)
ON CONFLICT (id) DO NOTHING;

CREATE TABLE bot_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    kind VARCHAR(32) NOT NULL,
    target_id UUID NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT,
    available_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bot_jobs_status_available ON bot_jobs (status, available_at);
CREATE INDEX idx_bot_jobs_target ON bot_jobs (kind, target_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bot_jobs;

DELETE FROM users WHERE id = 'a0000000-0000-4000-8000-000000000001';

DROP INDEX IF EXISTS idx_users_google_sub;

ALTER TABLE users
    DROP COLUMN IF EXISTS is_bot,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS google_sub;

ALTER TABLE users
    ALTER COLUMN password SET NOT NULL;
-- +goose StatementEnd
