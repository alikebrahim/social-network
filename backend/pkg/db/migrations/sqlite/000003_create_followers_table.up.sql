-- 000003_create_followers_table.up.sql
CREATE TABLE followers (
    id TEXT PRIMARY KEY,               -- Unique ID for the follow relationship (UUID)
    follower_id TEXT NOT NULL,         -- The user who is following
    following_id TEXT NOT NULL,        -- The user being followed
    status TEXT NOT NULL CHECK (status IN ('pending', 'accepted', 'declined')),  -- Follow request status
    created_at DATETIME NOT NULL,      -- Timestamp of when the follow request was created
    updated_at DATETIME,               -- Timestamp of when the follow request status was last updated
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,  -- Foreign key constraint to the users table
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE  -- Foreign key constraint to the users table
);
