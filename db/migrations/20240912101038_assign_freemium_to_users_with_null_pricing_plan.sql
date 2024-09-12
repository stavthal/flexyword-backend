-- +goose Up
-- +goose StatementBegin
-- Ensure there's a Freemium plan available
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pricing_plans WHERE name = 'Freemium') THEN
        INSERT INTO pricing_plans (name, translation_limit, languages_limit, phrase_length_limit, token_limit, advanced_features, priority_support, price_per_month, created_at, updated_at)
        VALUES ('Freemium', 10, 2, 100, 1000, FALSE, FALSE, 0.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
    END IF;
END $$;

-- Set users with null pricing_plan_id to Freemium plan
WITH freemium_plan AS (
    SELECT id FROM pricing_plans WHERE name = 'Freemium'
)
UPDATE users 
SET pricing_plan_id = (SELECT id FROM freemium_plan) 
WHERE pricing_plan_id IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Optionally, reverse the update by setting the pricing_plan_id to NULL again
UPDATE users SET pricing_plan_id = NULL WHERE pricing_plan_id = (SELECT id FROM pricing_plans WHERE name = 'Freemium');
-- +goose StatementEnd
