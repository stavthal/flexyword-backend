-- +goose Up
-- +goose StatementBegin
-- Add missing fields to the users table

-- Add FirstName column if it does not exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS first_name VARCHAR(255) NOT NULL DEFAULT '';

-- Add LastName column if it does not exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_name VARCHAR(255) NOT NULL DEFAULT '';

-- Add PricingPlanID column if it does not exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS pricing_plan_id INT;

-- Add UsedTokens column if it does not exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS used_tokens INT DEFAULT 0;

-- Add LastReset column if it does not exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_reset TIMESTAMP;

-- Add BillingAddress column if it does not exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS billing_address VARCHAR(255);

-- Create foreign key constraint for PricingPlanID if it does not exist
DO $$ 
BEGIN 
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_type = 'FOREIGN KEY' 
        AND table_name = 'users' 
        AND constraint_name = 'fk_users_pricing_plan_id'
    ) THEN 
        ALTER TABLE users ADD CONSTRAINT fk_users_pricing_plan_id 
        FOREIGN KEY (pricing_plan_id) 
        REFERENCES pricing_plans(id);
    END IF; 
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the columns that were added

-- Remove FirstName column
ALTER TABLE users DROP COLUMN IF EXISTS first_name;

-- Remove LastName column
ALTER TABLE users DROP COLUMN IF EXISTS last_name;

-- Remove PricingPlanID column
ALTER TABLE users DROP COLUMN IF EXISTS pricing_plan_id;

-- Remove UsedTokens column
ALTER TABLE users DROP COLUMN IF EXISTS used_tokens;

-- Remove LastReset column
ALTER TABLE users DROP COLUMN IF EXISTS last_reset;

-- Remove BillingAddress column
ALTER TABLE users DROP COLUMN IF EXISTS billing_address;

-- Drop foreign key constraint for PricingPlanID if exists
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_pricing_plan_id;
-- +goose StatementEnd
