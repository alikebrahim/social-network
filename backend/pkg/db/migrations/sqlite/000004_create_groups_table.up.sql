-- 000004_create_groups_table.up.sql
CREATE TABLE groups (
    id TEXT PRIMARY KEY,                 -- Unique group ID (UUID)
    title TEXT NOT NULL,                  -- The title of the group
    description TEXT NOT NULL,            -- A brief description of the group
    owner_id TEXT NOT NULL,               -- The ID of the user who created the group (foreign key to users table)
    created_at DATETIME NOT NULL,         -- Timestamp when the group was created
    updated_at DATETIME,                  -- Timestamp when the group was last updated
    is_public BOOLEAN NOT NULL DEFAULT 1,  -- Flag to indicate if the group is public (default: true)
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE  -- Foreign key reference to the users table
);
