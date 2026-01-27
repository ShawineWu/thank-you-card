-- Create employee_milestones table for tracking recognition achievements
CREATE TABLE IF NOT EXISTS employee_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id VARCHAR(255) NOT NULL,
    milestone_type VARCHAR(50) NOT NULL CHECK (milestone_type IN ('CardsSent', 'CardsReceived')),
    threshold INTEGER NOT NULL,
    achieved_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL
);

-- Create indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_employee_milestones_employee_id ON employee_milestones(employee_id);
CREATE INDEX IF NOT EXISTS idx_employee_milestones_type ON employee_milestones(milestone_type);

-- Ensure each milestone is achieved only once per employee
CREATE UNIQUE INDEX IF NOT EXISTS idx_employee_milestones_unique 
    ON employee_milestones(employee_id, milestone_type, threshold);

-- Add comment
COMMENT ON TABLE employee_milestones IS 'Tracks milestone achievements for employee recognition activity';
COMMENT ON COLUMN employee_milestones.milestone_type IS 'Either "CardsSent" or "CardsReceived"';
COMMENT ON COLUMN employee_milestones.threshold IS 'Number of cards that triggered this milestone (e.g., 1, 5, 10, 25, 50, 100)';
