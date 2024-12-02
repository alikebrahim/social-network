-- 000005_create_notifications_table.up.sql
CREATE TABLE notifications (
    id TEXT PRIMARY KEY,               -- Unique notification ID (UUID)
    user_id TEXT NOT NULL,             -- User ID (foreign key to users table)
    message TEXT NOT NULL,             -- The content of the notification
    is_read BOOLEAN NOT NULL DEFAULT 0, -- Flag to indicate if the notification has been read (default: false)
    created_at DATETIME NOT NULL,      -- Timestamp when the notification was created
    updated_at DATETIME,               -- Timestamp when the notification was last updated
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE  -- Foreign key constraint to the users table
);
