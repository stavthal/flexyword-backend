-- +goose Up
-- Enable UUID extension (if not already enabled)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Add a new UUID column (temporary)
ALTER TABLE translations ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();

-- Copy data from the old ID to the new UUID column
UPDATE translations SET id_new = uuid_generate_v4();

-- Drop the old ID column and rename the new one
ALTER TABLE translations DROP COLUMN id;
ALTER TABLE translations RENAME COLUMN id_new TO id;

-- Set the new UUID column as the primary key
ALTER TABLE translations ADD PRIMARY KEY (id);

-- +goose Down
-- Reverse the changes by adding the old id back as an integer primary key
ALTER TABLE translations ADD COLUMN id SERIAL PRIMARY KEY;
