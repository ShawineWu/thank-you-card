-- Create company_values table for storing company values and credos
CREATE TABLE IF NOT EXISTS company_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('Value', 'Credo')),
    examples JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_company_values_type ON company_values(type);
CREATE INDEX IF NOT EXISTS idx_company_values_name ON company_values(name);
CREATE INDEX IF NOT EXISTS idx_company_values_code ON company_values(code);

-- Insert company values
INSERT INTO company_values (code, name, description, type) VALUES
    ('MAKE_IMPACT', 'Make an Impact', 'Be driven by the desire to build something that can touch millions of lives.', 'Value'),
    ('STRIVE_EXCELLENCE', 'Strive for Excellence', 'Today''s great is not good enough for tomorrow.', 'Value'),
    ('STAND_TOGETHER', 'Stand Together', 'Embrace the spirit of solidarity & commitment through thick and thin.', 'Value'),
    ('BE_OPEN_MINDED', 'Be Open-Minded', 'Diversity is strength. Magic often happens at the intersection of two different worlds.', 'Value'),
    ('STAY_GROUNDED', 'Stay Grounded', 'Always remember our roots and higher purpose. Plus, life is too short to be around jerks.', 'Value')
ON CONFLICT (code) DO NOTHING;

-- Insert company credos (15 credos as per seed data)
INSERT INTO company_values (code, name, description, type) VALUES
    ('BIAS_FOR_ACTION', 'Bias for Action', 'Speed matters, take action to deliver a high quality result with calculated risk taking. Be time conscious and set deadlines on initiatives. Promptly ask for help when faced with roadblocks, be the roadblock-remover where we can. It is better to be moving than not, most consequences are manageable.', 'Credo'),
    ('CUSTOMER_CENTRIC', 'Customer Centric', 'We start with the customer experience, centring our processes and decisions around it. Never stick to what has been done for the sake of tradition. Advocate for our customers, add value wherever possible and build trust with sincere interactions.', 'Credo'),
    ('THINK_STRATEGICALLY', 'Think Strategically', 'We know our business inside out, acutely aware of the market dynamics and crystal clear of our positioning and what we offer. We stay ahead of the curve by being flexible enough to pivot quickly, being nimble in problem solving and being obsessed with value creation for our customers, employees and partners. We connect the dots where others don''t.', 'Credo'),
    ('DEEP_DIVE', 'Deep Dive', 'No detail is too small, no task too unimportant. Always ask why, validate with data, slicing for insights and looking for opportunities. Only when we know the nuts and bolts of everything we work on, can we troubleshoot effectively and innovate creatively. Identify root causes. When in doubt, keep questioning constructively.', 'Credo'),
    ('INVENT_SIMPLIFY', 'Invent and Simplify', 'There is no need to over-complicate. Get creative with solutions, streamline where we can and tap on the expertise of different teams. Innovation and invention are expected and achieved by leveraging all available resources.', 'Credo'),
    ('EARN_TRUST', 'Earn Trust', 'Always deliver on promises. We listen attentively, speak candidly and maintain the highest standards of ethics - towards our customers and our team.', 'Credo'),
    ('TAKE_OWNERSHIP', 'Take Ownership', 'Each of us represents the company, beyond just ourselves or our team. When times are good, we celebrate together. When times are bad, we stay and face the adversity as one. Act with a view of the long term, today''s great piece of work will have a lasting impact on the bigger picture. The work we produce speaks for the company, be proud of it, own it.', 'Credo'),
    ('CHALLENGE_DISAGREE_COMMIT', 'Challenge Disagree and Commit', 'When in doubt, challenge respectfully and objectively, even if it is uncomfortable. Unemotionally review the objective of the things we are doing and whether how we are doing it is the best way to do it. Present facts instead of succumbing to feelings. Ideas evolve and improve when scrutinised. Once a decision is determined, commit wholeheartedly and own the consequences together.', 'Credo'),
    ('LEARN_BE_CURIOUS', 'Learn and Be Curious', 'Learn from team mates, customers and competitors. Always be curious about the whys and how things are done, be eager to go deeper and bring our learnings back to our work.', 'Credo'),
    ('DO_MORE_WITH_LESS', 'Do More with Less', 'Know where we should be investing, and invest it better. Everyone is empowered to implement efficient processes and reduce spending. Do so wisely because we cut costs but we don''t cut corners.', 'Credo')
ON CONFLICT (code) DO NOTHING;

-- Add comment
COMMENT ON TABLE company_values IS 'Stores company values and credos that can be associated with cards';
COMMENT ON COLUMN company_values.type IS 'Either "Value" for core company values or "Credo" for operational principles';
