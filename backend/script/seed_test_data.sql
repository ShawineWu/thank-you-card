-- Test Data Seed Script
-- Run this AFTER the Go application has created all tables via AutoMigrate
-- Usage: psql -h localhost -U postgres -d thankyoucard -f seed_test_data.sql

-- Additional test employees
DO $$
BEGIN
    INSERT INTO employees (email, name, department, is_hr_admin, created_at, updated_at)
    VALUES
        ('alice.smith@example.com', 'Alice Smith', 'Engineering', false, NOW(), NOW()),
        ('bob.jones@example.com', 'Bob Jones', 'Engineering', false, NOW(), NOW()),
        ('carol.white@example.com', 'Carol White', 'Design', false, NOW(), NOW()),
        ('david.brown@example.com', 'David Brown', 'Product', false, NOW(), NOW()),
        ('eve.green@example.com', 'Eve Green', 'Marketing', false, NOW(), NOW()),
        ('frank.taylor@example.com', 'Frank Taylor', 'Sales', false, NOW(), NOW()),
        ('grace.lee@example.com', 'Grace Lee', 'Engineering', false, NOW(), NOW()),
        ('henry.wang@example.com', 'Henry Wang', 'Design', false, NOW(), NOW())
    ON CONFLICT (email) DO NOTHING;
END $$;

-- Card 1: Employee 1 sends to Alice and Bob (2 days ago)
DO $$
DECLARE
    card1_id INTEGER;
    emp1_id INTEGER;
    alice_id INTEGER;
    bob_id INTEGER;
    value1_id INTEGER;
    value2_id INTEGER;
BEGIN
    SELECT id INTO emp1_id FROM employees WHERE email = 'employee@example.com' LIMIT 1;
    SELECT id INTO alice_id FROM employees WHERE email = 'alice.smith@example.com' LIMIT 1;
    SELECT id INTO bob_id FROM employees WHERE email = 'bob.jones@example.com' LIMIT 1;
    SELECT id INTO value1_id FROM company_values WHERE code = 'STAND_TOGETHER' LIMIT 1;
    SELECT id INTO value2_id FROM company_values WHERE code = 'STRIVE_EXCELLENCE' LIMIT 1;
    
    IF emp1_id IS NOT NULL AND alice_id IS NOT NULL AND bob_id IS NOT NULL THEN
        INSERT INTO cards (sender_id, reason, created_at, updated_at)
        VALUES (emp1_id, 'Thank you Alice and Bob for your excellent collaboration on the recent project! Your teamwork and dedication made all the difference.', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days')
        RETURNING id INTO card1_id;
        
        INSERT INTO card_recipients (card_id, recipient_id, created_at, updated_at) VALUES
            (card1_id, alice_id, NOW(), NOW()),
            (card1_id, bob_id, NOW(), NOW());
        
        IF value1_id IS NOT NULL AND value2_id IS NOT NULL THEN
            INSERT INTO card_values (card_id, company_value_id, created_at, updated_at) VALUES
                (card1_id, value1_id, NOW(), NOW()),
                (card1_id, value2_id, NOW(), NOW());
        END IF;
    END IF;
END $$;

-- Card 2: Alice sends to Employee 1 (1 day ago)
DO $$
DECLARE
    card2_id INTEGER;
    alice_id INTEGER;
    emp1_id INTEGER;
    value3_id INTEGER;
    value4_id INTEGER;
BEGIN
    SELECT id INTO alice_id FROM employees WHERE email = 'alice.smith@example.com' LIMIT 1;
    SELECT id INTO emp1_id FROM employees WHERE email = 'employee@example.com' LIMIT 1;
    SELECT id INTO value3_id FROM company_values WHERE code = 'EARN_TRUST' LIMIT 1;
    SELECT id INTO value4_id FROM company_values WHERE code = 'DEEP_DIVE' LIMIT 1;
    
    IF alice_id IS NOT NULL AND emp1_id IS NOT NULL THEN
        INSERT INTO cards (sender_id, reason, created_at, updated_at)
        VALUES (alice_id, 'Thanks for the great code review and helpful feedback! Really appreciate your attention to detail.', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day')
        RETURNING id INTO card2_id;
        
        INSERT INTO card_recipients (card_id, recipient_id, created_at, updated_at) VALUES
            (card2_id, emp1_id, NOW(), NOW());
        
        IF value3_id IS NOT NULL AND value4_id IS NOT NULL THEN
            INSERT INTO card_values (card_id, company_value_id, created_at, updated_at) VALUES
                (card2_id, value3_id, NOW(), NOW()),
                (card2_id, value4_id, NOW(), NOW());
        END IF;
    END IF;
END $$;

-- Card 3: HR Admin sends to Engineering team (5 days ago)
DO $$
DECLARE
    card3_id INTEGER;
    hr_admin_id INTEGER;
    value5_id INTEGER;
    value6_id INTEGER;
    value7_id INTEGER;
    eng_emp RECORD;
BEGIN
    SELECT id INTO hr_admin_id FROM employees WHERE email = 'hradmin@example.com' LIMIT 1;
    SELECT id INTO value5_id FROM company_values WHERE code = 'STRIVE_EXCELLENCE' LIMIT 1;
    SELECT id INTO value6_id FROM company_values WHERE code = 'TAKE_OWNERSHIP' LIMIT 1;
    SELECT id INTO value7_id FROM company_values WHERE code = 'BIAS_FOR_ACTION' LIMIT 1;
    
    IF hr_admin_id IS NOT NULL THEN
        INSERT INTO cards (sender_id, reason, created_at, updated_at)
        VALUES (hr_admin_id, 'Thank you to the entire Engineering team for delivering the Q4 features on time! Your hard work and commitment to excellence is truly appreciated.', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days')
        RETURNING id INTO card3_id;
        
        FOR eng_emp IN SELECT id FROM employees WHERE department = 'Engineering' AND id != hr_admin_id
        LOOP
            INSERT INTO card_recipients (card_id, recipient_id, created_at, updated_at) VALUES
                (card3_id, eng_emp.id, NOW(), NOW())
            ON CONFLICT DO NOTHING;
        END LOOP;
        
        IF value5_id IS NOT NULL AND value6_id IS NOT NULL AND value7_id IS NOT NULL THEN
            INSERT INTO card_values (card_id, company_value_id, created_at, updated_at) VALUES
                (card3_id, value5_id, NOW(), NOW()),
                (card3_id, value6_id, NOW(), NOW()),
                (card3_id, value7_id, NOW(), NOW());
        END IF;
    END IF;
END $$;

-- Card 4: Employee sends to HR Admin (7 days ago)
DO $$
DECLARE
    card4_id INTEGER;
    emp1_id INTEGER;
    hr_admin_id INTEGER;
    value8_id INTEGER;
    value9_id INTEGER;
BEGIN
    SELECT id INTO emp1_id FROM employees WHERE email = 'employee@example.com' LIMIT 1;
    SELECT id INTO hr_admin_id FROM employees WHERE email = 'hradmin@example.com' LIMIT 1;
    SELECT id INTO value8_id FROM company_values WHERE code = 'STAND_TOGETHER' LIMIT 1;
    SELECT id INTO value9_id FROM company_values WHERE code = 'CUSTOMER_CENTRIC' LIMIT 1;
    
    IF emp1_id IS NOT NULL AND hr_admin_id IS NOT NULL THEN
        INSERT INTO cards (sender_id, reason, created_at, updated_at)
        VALUES (emp1_id, 'Thank you for organizing the team building event! It was a great opportunity to connect with colleagues.', NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days')
        RETURNING id INTO card4_id;
        
        INSERT INTO card_recipients (card_id, recipient_id, created_at, updated_at) VALUES
            (card4_id, hr_admin_id, NOW(), NOW());
        
        IF value8_id IS NOT NULL AND value9_id IS NOT NULL THEN
            INSERT INTO card_values (card_id, company_value_id, created_at, updated_at) VALUES
                (card4_id, value8_id, NOW(), NOW()),
                (card4_id, value9_id, NOW(), NOW());
        END IF;
    END IF;
END $$;

-- Card 5: Bob sends to Carol and Henry (3 days ago)
DO $$
DECLARE
    card5_id INTEGER;
    bob_id INTEGER;
    carol_id INTEGER;
    henry_id INTEGER;
    value10_id INTEGER;
    value11_id INTEGER;
    value12_id INTEGER;
BEGIN
    SELECT id INTO bob_id FROM employees WHERE email = 'bob.jones@example.com' LIMIT 1;
    SELECT id INTO carol_id FROM employees WHERE email = 'carol.white@example.com' LIMIT 1;
    SELECT id INTO henry_id FROM employees WHERE email = 'henry.wang@example.com' LIMIT 1;
    SELECT id INTO value10_id FROM company_values WHERE code = 'INVENT_SIMPLIFY' LIMIT 1;
    SELECT id INTO value11_id FROM company_values WHERE code = 'STRIVE_EXCELLENCE' LIMIT 1;
    SELECT id INTO value12_id FROM company_values WHERE code = 'DEEP_DIVE' LIMIT 1;
    
    IF bob_id IS NOT NULL AND carol_id IS NOT NULL AND henry_id IS NOT NULL THEN
        INSERT INTO cards (sender_id, reason, created_at, updated_at)
        VALUES (bob_id, 'Amazing work on the design system! The new components are beautiful and well-documented. Thank you Carol and Henry!', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days')
        RETURNING id INTO card5_id;
        
        INSERT INTO card_recipients (card_id, recipient_id, created_at, updated_at) VALUES
            (card5_id, carol_id, NOW(), NOW()),
            (card5_id, henry_id, NOW(), NOW());
        
        IF value10_id IS NOT NULL AND value11_id IS NOT NULL AND value12_id IS NOT NULL THEN
            INSERT INTO card_values (card_id, company_value_id, created_at, updated_at) VALUES
                (card5_id, value10_id, NOW(), NOW()),
                (card5_id, value11_id, NOW(), NOW()),
                (card5_id, value12_id, NOW(), NOW());
        END IF;
    END IF;
END $$;

-- Add more cards for better testing
-- Card 6: Carol sends to Eve (4 days ago)
DO $$
DECLARE
    card6_id INTEGER;
    carol_id INTEGER;
    eve_id INTEGER;
    value13_id INTEGER;
BEGIN
    SELECT id INTO carol_id FROM employees WHERE email = 'carol.white@example.com' LIMIT 1;
    SELECT id INTO eve_id FROM employees WHERE email = 'eve.green@example.com' LIMIT 1;
    SELECT id INTO value13_id FROM company_values WHERE code = 'BE_OPEN_MINDED' LIMIT 1;
    
    IF carol_id IS NOT NULL AND eve_id IS NOT NULL THEN
        INSERT INTO cards (sender_id, reason, created_at, updated_at)
        VALUES (carol_id, 'Thank you Eve for the amazing marketing campaign! Your creativity and strategic thinking really made it shine.', NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days')
        RETURNING id INTO card6_id;
        
        INSERT INTO card_recipients (card_id, recipient_id, created_at, updated_at) VALUES
            (card6_id, eve_id, NOW(), NOW());
        
        IF value13_id IS NOT NULL THEN
            INSERT INTO card_values (card_id, company_value_id, created_at, updated_at) VALUES
                (card6_id, value13_id, NOW(), NOW());
        END IF;
    END IF;
END $$;

-- Card 7: David sends to Frank (6 days ago)
DO $$
DECLARE
    card7_id INTEGER;
    david_id INTEGER;
    frank_id INTEGER;
    value14_id INTEGER;
    value15_id INTEGER;
BEGIN
    SELECT id INTO david_id FROM employees WHERE email = 'david.brown@example.com' LIMIT 1;
    SELECT id INTO frank_id FROM employees WHERE email = 'frank.taylor@example.com' LIMIT 1;
    SELECT id INTO value14_id FROM company_values WHERE code = 'CUSTOMER_CENTRIC' LIMIT 1;
    SELECT id INTO value15_id FROM company_values WHERE code = 'EARN_TRUST' LIMIT 1;
    
    IF david_id IS NOT NULL AND frank_id IS NOT NULL THEN
        INSERT INTO cards (sender_id, reason, created_at, updated_at)
        VALUES (david_id, 'Great job on closing that big deal! Your persistence and customer focus really paid off.', NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days')
        RETURNING id INTO card7_id;
        
        INSERT INTO card_recipients (card_id, recipient_id, created_at, updated_at) VALUES
            (card7_id, frank_id, NOW(), NOW());
        
        IF value14_id IS NOT NULL AND value15_id IS NOT NULL THEN
            INSERT INTO card_values (card_id, company_value_id, created_at, updated_at) VALUES
                (card7_id, value14_id, NOW(), NOW()),
                (card7_id, value15_id, NOW(), NOW());
        END IF;
    END IF;
END $$;

-- Add emoji reactions
DO $$
DECLARE
    card1_id INTEGER;
    card2_id INTEGER;
    card3_id INTEGER;
    alice_id INTEGER;
    bob_id INTEGER;
    emp1_id INTEGER;
    eng_emp RECORD;
    reaction_count INTEGER := 0;
BEGIN
    -- Get card IDs (most recent cards)
    SELECT id INTO card1_id FROM cards WHERE sender_id = (SELECT id FROM employees WHERE email = 'employee@example.com' LIMIT 1) ORDER BY created_at DESC LIMIT 1;
    SELECT id INTO card2_id FROM cards WHERE sender_id = (SELECT id FROM employees WHERE email = 'alice.smith@example.com' LIMIT 1) ORDER BY created_at DESC LIMIT 1;
    SELECT id INTO card3_id FROM cards WHERE sender_id = (SELECT id FROM employees WHERE email = 'hradmin@example.com' LIMIT 1) ORDER BY created_at DESC LIMIT 1;
    
    SELECT id INTO alice_id FROM employees WHERE email = 'alice.smith@example.com' LIMIT 1;
    SELECT id INTO bob_id FROM employees WHERE email = 'bob.jones@example.com' LIMIT 1;
    SELECT id INTO emp1_id FROM employees WHERE email = 'employee@example.com' LIMIT 1;
    
    -- Reactions on card 1
    IF card1_id IS NOT NULL AND alice_id IS NOT NULL AND bob_id IS NOT NULL THEN
        INSERT INTO emoji_reactions (card_id, user_id, emoji_code, created_at, updated_at) VALUES
            (card1_id, alice_id, '👍', NOW(), NOW()),
            (card1_id, bob_id, '❤️', NOW(), NOW())
        ON CONFLICT (card_id, user_id) DO NOTHING;
    END IF;
    
    -- Reaction on card 2
    IF card2_id IS NOT NULL AND emp1_id IS NOT NULL THEN
        INSERT INTO emoji_reactions (card_id, user_id, emoji_code, created_at, updated_at) VALUES
            (card2_id, emp1_id, '👍', NOW(), NOW())
        ON CONFLICT (card_id, user_id) DO NOTHING;
    END IF;
    
    -- Reactions on card 3 (Engineering team card)
    IF card3_id IS NOT NULL THEN
        FOR eng_emp IN SELECT id FROM employees WHERE department = 'Engineering' AND id != (SELECT id FROM employees WHERE email = 'hradmin@example.com' LIMIT 1) LIMIT 3
        LOOP
            INSERT INTO emoji_reactions (card_id, user_id, emoji_code, created_at, updated_at) VALUES
                (card3_id, eng_emp.id, '🎉', NOW(), NOW())
            ON CONFLICT (card_id, user_id) DO NOTHING;
            reaction_count := reaction_count + 1;
            EXIT WHEN reaction_count >= 3;
        END LOOP;
    END IF;
END $$;

-- Summary
SELECT 'Test data seeded successfully!' AS message;
SELECT COUNT(*) AS total_employees FROM employees;
SELECT COUNT(*) AS total_cards FROM cards;
SELECT COUNT(*) AS total_reactions FROM emoji_reactions;
