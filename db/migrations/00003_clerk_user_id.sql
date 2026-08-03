-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN clerk_user_id VARCHAR(255) UNIQUE;

CREATE INDEX idx_users_clerk_user_id ON users (clerk_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_clerk_user_id;

ALTER TABLE users
    DROP COLUMN IF EXISTS clerk_user_id;
-- +goose StatementEnd
