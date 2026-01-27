# Database setup script for Thank You Card system

# Create database
CREATE DATABASE thankyoucard;

# Create user
CREATE USER thankyoucard WITH ENCRYPTED PASSWORD 'password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE thankyoucard TO thankyoucard;

# Connect to database
\c thankyoucard

# Grant schema privileges
GRANT ALL ON SCHEMA public TO thankyoucard;
