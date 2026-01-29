-- Create card_recipients junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS card_recipients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL,
    recipient_id VARCHAR(255) NOT NULL,
    recipient_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_card_recipients_card_id ON card_recipients(card_id);
CREATE INDEX IF NOT EXISTS idx_card_recipients_recipient_id ON card_recipients(recipient_id);

-- Ensure unique card-recipient combinations
CREATE UNIQUE INDEX IF NOT EXISTS idx_card_recipients_unique ON card_recipients(card_id, recipient_id);

-- Add comment
COMMENT ON TABLE card_recipients IS 'Junction table linking cards to their recipients (many-to-many)';
COMMENT ON COLUMN card_recipients.recipient_id IS 'Employee ID of the card recipient';
