-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

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

-- Create User table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    pricing_plan_id INT,
    used_tokens INT DEFAULT 0,
    last_reset TIMESTAMP,
    billing_address VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pricing_plan_id) REFERENCES pricing_plans(id)
);

-- Create Translation table
CREATE TABLE IF NOT EXISTS translations (
    id SERIAL PRIMARY KEY,
    phrase TEXT NOT NULL,
    input_language TEXT NOT NULL,
    output_languages JSONB NOT NULL,
    translation_result JSONB NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
-- Drop tables in reverse order
DROP TABLE IF EXISTS translations;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS pricing_plans;
