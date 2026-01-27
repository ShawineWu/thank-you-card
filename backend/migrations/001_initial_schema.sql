-- Initial schema for Thank You Card system
-- This file contains the complete database schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Company Values table
CREATE TABLE company_values (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('Value', 'Credo')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for company_values
CREATE INDEX idx_company_values_type ON company_values(type);
CREATE INDEX idx_company_values_name ON company_values(name);

-- Cards table
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id VARCHAR(255) NOT NULL,
    recognition_reason TEXT NOT NULL CHECK (length(recognition_reason) >= 10 AND length(recognition_reason) <= 1000),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for cards
CREATE INDEX idx_cards_sender_id ON cards(sender_id);
CREATE INDEX idx_cards_created_at ON cards(created_at DESC);

-- Card Recipients table
CREATE TABLE card_recipients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    recipient_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for card_recipients
CREATE INDEX idx_card_recipients_card_id ON card_recipients(card_id);
CREATE INDEX idx_card_recipients_recipient_id ON card_recipients(recipient_id);
CREATE UNIQUE INDEX idx_card_recipients_unique ON card_recipients(card_id, recipient_id);

-- Card Values table (many-to-many relationship between cards and company values)
CREATE TABLE card_values (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    value_id UUID NOT NULL REFERENCES company_values(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for card_values
CREATE INDEX idx_card_values_card_id ON card_values(card_id);
CREATE INDEX idx_card_values_value_id ON card_values(value_id);
CREATE UNIQUE INDEX idx_card_values_unique ON card_values(card_id, value_id);

-- Employee Milestones table
CREATE TABLE employee_milestones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id VARCHAR(255) NOT NULL,
    milestone_type VARCHAR(50) NOT NULL CHECK (milestone_type IN ('CardsSent', 'CardsReceived')),
    threshold INTEGER NOT NULL,
    achieved_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL
);

-- Create indexes for employee_milestones
CREATE INDEX idx_employee_milestones_employee_id ON employee_milestones(employee_id);
CREATE INDEX idx_employee_milestones_type ON employee_milestones(milestone_type);
CREATE UNIQUE INDEX idx_employee_milestones_unique ON employee_milestones(employee_id, milestone_type, threshold);

-- Create a function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_cards_updated_at BEFORE UPDATE ON cards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_company_values_updated_at BEFORE UPDATE ON company_values
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();