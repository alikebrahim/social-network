-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,                -- Unique user ID (can be UUID)
    email TEXT UNIQUE NOT NULL,          -- User's email (unique constraint)
    password TEXT NOT NULL,             -- User's hashed password (not stored in plain text)
    first_name TEXT NOT NULL,           -- User's first name
    last_name TEXT NOT NULL,            -- User's last name
    date_of_birth TEXT NOT NULL,        -- User's date of birth (format: YYYY-MM-DD)
    avatar TEXT,                        -- URL to the user's avatar image (optional)
    nickname TEXT,                      -- User's nickname (optional)
    about_me TEXT,                      -- Short bio or description about the user (optional)
    is_public BOOLEAN NOT NULL DEFAULT 1, -- Flag to indicate if the user's profile is public (default: true)
    created_at DATETIME NOT NULL        -- Timestamp when the user account was created
);
