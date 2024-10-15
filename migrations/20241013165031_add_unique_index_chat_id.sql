-- +goose Up
-- +goose NO TRANSACTION
-- +goose StatementBegin
CREATE UNIQUE INDEX CONCURRENTLY idx_chat_id_users
ON users (chat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX CONCURRENTLY IF EXISTS idx_chat_id_users;
-- +goose StatementEnd