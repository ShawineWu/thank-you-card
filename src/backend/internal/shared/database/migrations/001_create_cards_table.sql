-- Create cards table for storing thank you cards
CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id VARCHAR(255) NOT NULL,
    recognition_reason TEXT NOT NULL CHECK (
        length(recognition_reason) >= 10 AND 
        length(recognition_reason) <= 1000
    ),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_cards_sender_id ON cards(sender_id);
CREATE INDEX IF NOT EXISTS idx_cards_created_at ON cards(created_at DESC);

-- Add comment
COMMENT ON TABLE cards IS 'Stores recognition/thank you cards sent between employees';
COMMENT ON COLUMN cards.sender_id IS 'Employee ID of the card sender';
COMMENT ON COLUMN cards.recognition_reason IS 'Reason for recognition, must be between 10 and 1000 characters';
