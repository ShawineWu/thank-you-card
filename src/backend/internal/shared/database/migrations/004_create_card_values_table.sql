-- Create card_values junction table for card-value associations
CREATE TABLE IF NOT EXISTS card_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    value_id UUID NOT NULL REFERENCES company_values(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_card_values_card_id ON card_values(card_id);
CREATE INDEX IF NOT EXISTS idx_card_values_value_id ON card_values(value_id);

-- Ensure unique card-value combinations
CREATE UNIQUE INDEX IF NOT EXISTS idx_card_values_unique ON card_values(card_id, value_id);

-- Add comment
COMMENT ON TABLE card_values IS 'Junction table linking cards to selected company values/credos (many-to-many)';
COMMENT ON COLUMN card_values.value_id IS 'Reference to company value or credo';
