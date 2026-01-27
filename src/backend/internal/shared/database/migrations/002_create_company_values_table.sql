-- Create company_values table for storing company values and credos
CREATE TABLE IF NOT EXISTS company_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('Value', 'Credo')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_company_values_type ON company_values(type);
CREATE INDEX IF NOT EXISTS idx_company_values_name ON company_values(name);

-- Insert company values
INSERT INTO company_values (name, description, type) VALUES
    ('Make an Impact', 'Be driven by the desire to build something that can touch millions of lives.', 'Value'),
    ('Strive for Excellence', 'Today''s great is not good enough for tomorrow.', 'Value'),
    ('Stand Together', 'Embrace the spirit of solidarity & commitment through thick and thin.', 'Value'),
    ('Be Open-Minded', 'Diversity is strength. Magic often happens at the intersection of two different worlds.', 'Value'),
    ('Stay Grounded', 'Always remember our roots and higher purpose. Plus, life is too short to be around jerks.', 'Value')
ON CONFLICT (name) DO NOTHING;

-- Insert company credos (10 credos as per inception documents)
INSERT INTO company_values (name, description, type) VALUES
    ('Bias for Action', 'Take initiative and move forward decisively', 'Credo'),
    ('Customer Centric', 'Put customer needs at the heart of everything we do', 'Credo'),
    ('Think Strategically', 'Consider long-term implications and broader context', 'Credo'),
    ('Data-Driven Decisions', 'Base decisions on evidence and analytics', 'Credo'),
    ('Continuous Learning', 'Commit to ongoing personal and professional growth', 'Credo'),
    ('Ownership Mindset', 'Take responsibility for outcomes and results', 'Credo'),
    ('Collaborate Effectively', 'Work together to achieve shared goals', 'Credo'),
    ('Communicate Clearly', 'Express ideas and information transparently', 'Credo'),
    ('Innovate Constantly', 'Seek new and better ways to solve problems', 'Credo'),
    ('Deliver Quality', 'Ensure excellence in all work outputs', 'Credo')
ON CONFLICT (name) DO NOTHING;

-- Add comment
COMMENT ON TABLE company_values IS 'Stores company values and credos that can be associated with cards';
COMMENT ON COLUMN company_values.type IS 'Either "Value" for core company values or "Credo" for operational principles';
