-- +goose Up
ALTER TABLE translations ADD COLUMN input_language TEXT DEFAULT 'English' NOT NULL;

-- +goose Down
ALTER TABLE translations DROP COLUMN input_language;
