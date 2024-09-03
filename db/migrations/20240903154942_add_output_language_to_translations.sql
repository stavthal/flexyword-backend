-- +goose Up
ALTER TABLE translations ADD COLUMN output_languages TEXT DEFAULT '[]' NOT NULL;

-- +goose Down
ALTER TABLE translations DROP COLUMN output_languages;
