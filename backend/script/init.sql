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