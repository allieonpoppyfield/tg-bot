-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN photo_url VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN photo_url;
-- +goose StatementEnd
