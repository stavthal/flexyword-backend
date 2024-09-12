-- +goose Up
-- +goose StatementBegin
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_attribute 
                   WHERE attrelid = 'users'::regclass 
                   AND attname = 'first_name' 
                   AND NOT attisdropped) THEN
        ALTER TABLE users ADD COLUMN first_name VARCHAR(50);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_attribute 
                   WHERE attrelid = 'users'::regclass 
                   AND attname = 'last_name' 
                   AND NOT attisdropped) THEN
        ALTER TABLE users ADD COLUMN last_name VARCHAR(50);
    END IF;
    
    -- Update existing NULL fields with default empty string
    UPDATE users SET first_name = '' WHERE first_name IS NULL;
    UPDATE users SET last_name = '' WHERE last_name IS NULL;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS first_name;
ALTER TABLE users DROP COLUMN IF EXISTS last_name;
-- +goose StatementEnd
