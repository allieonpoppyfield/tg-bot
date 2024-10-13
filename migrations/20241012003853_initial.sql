-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  chat_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  age INT NOT NULL,
  gender INT2 NOT NULL,
  description TEXT NOT NULL,
  created_at DATE NOT NULL DEFAULT NOW(),
  updated_at DATE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
