-- +goose Up
-- +goose StatementBegin
-- Add deleted_at column to users, pricing_plans, and translations tables for soft delete support

ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

ALTER TABLE pricing_plans ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

ALTER TABLE translations ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the deleted_at column if rolling back
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE pricing_plans DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE translations DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd
