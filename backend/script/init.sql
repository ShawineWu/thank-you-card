-- Thank You Card Backend - Database Initialization Script
-- This script creates the database and sets up basic configuration
-- Run this script as a PostgreSQL superuser (e.g., postgres)
-- 
-- Database Configuration:
--   DBName: thankyoucard
--   User: postgres
--   Host: localhost
--   Port: 5432
--   Password: 123456

-- Create database if it doesn't exist
SELECT 'CREATE DATABASE thankyoucard'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'thankyoucard')\gexec

-- Connect to the database
\c thankyoucard

-- Create extension if needed (for UUID, etc.)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Note: Table creation is handled by Gorm AutoMigrate
-- This script only creates the database itself
-- Run the Go application (make backend-dev) to execute migrations and seed data
-- The migration will automatically:
--   1. Create all tables (employees, company_values, cards, etc.)
--   2. Seed company values and credos
--   3. Seed mock employees (ID: 1 = regular, ID: 2 = HR admin)
--   4. Seed test user accounts (employee, hr, admin)

-- ============================================
-- TEST DATA - Execute after tables are created
-- ============================================
-- IMPORTANT: Run the seed script AFTER running the Go application once
-- to ensure all tables are created by Gorm AutoMigrate
-- 
-- To seed test data:
-- 1. Start the backend once: make backend-dev (then stop it)
-- 2. Run the seed script: psql -h localhost -U postgres -d thankyoucard -f backend/script/seed_test_data.sql
-- 
-- The seed script will create:
-- - 8 additional test employees (Alice, Bob, Carol, David, Eve, Frank, Grace, Henry)
-- - 7 test cards with various scenarios (single/multiple recipients, different dates)
-- - Card recipients and value associations
-- - Emoji reactions on cards
