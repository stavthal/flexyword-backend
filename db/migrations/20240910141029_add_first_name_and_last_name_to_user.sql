-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN first_name VARCHAR(50);
ALTER TABLE users ADD COLUMN last_name VARCHAR(50);
UPDATE users SET first_name = '' WHERE first_name IS NULL;
UPDATE users SET last_name = '' WHERE last_name IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN first_name;
ALTER TABLE users DROP COLUMN last_name;
-- +goose StatementEnd