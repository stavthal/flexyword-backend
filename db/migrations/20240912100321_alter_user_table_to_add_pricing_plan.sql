-- +goose Up
-- Add pricing_plan_id foreign key to the users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS pricing_plan_id INT,
ADD CONSTRAINT fk_pricing_plan FOREIGN KEY (pricing_plan_id) REFERENCES pricing_plans(id);

-- +goose Down
-- Remove the pricing_plan_id foreign key from the users table
ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_pricing_plan,
DROP COLUMN IF EXISTS pricing_plan_id;
