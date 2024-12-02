-- 000002_create_posts_table.up.sql
CREATE TABLE posts (
    id TEXT PRIMARY KEY,                -- Unique post ID (can be UUID)
    user_id TEXT NOT NULL,              -- User ID (foreign key to users table)
    content TEXT NOT NULL,              -- Content of the post
    image_url TEXT,                     -- URL to an image associated with the post (optional)
    privacy TEXT NOT NULL CHECK (privacy IN ('public', 'private', 'almost_private')), -- Privacy of the post
    created_at DATETIME NOT NULL,       -- Timestamp of when the post was created
    updated_at DATETIME,                -- Timestamp of when the post was last updated
    FOREIGN KEY (user_id) REFERENCES users(id) -- Foreign key constraint to the users table
);
