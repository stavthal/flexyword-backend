-- +goose Up
-- Create PricingPlan table
CREATE TABLE IF NOT EXISTS pricing_plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    translation_limit INT,
    languages_limit INT,
    phrase_length_limit INT,
    token_limit INT,
    advanced_features BOOLEAN,
    priority_support BOOLEAN,
    price_per_month DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default pricing plans
INSERT INTO pricing_plans (name, translation_limit, languages_limit, phrase_length_limit, token_limit, advanced_features, priority_support, price_per_month)
VALUES
    ('Freemium', 10, 2, 100, 1000, FALSE, FALSE, 0.00),
    ('Standard', 100, 5, 250, 5000, TRUE, FALSE, 9.99),
    ('Premium', 1000, 10, 500, 50000, TRUE, TRUE, 29.99);

-- +goose Down
-- Drop PricingPlan table if rolling back
DROP TABLE IF EXISTS pricing_plans;
