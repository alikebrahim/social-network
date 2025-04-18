-- +goose Up
ALTER TABLE notifications ADD COLUMN related_id INTEGER;
ALTER TABLE notifications ADD COLUMN sender_id INTEGER;

-- +goose Down
CREATE TABLE notifications_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type TEXT NOT NULL,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO notifications_backup (id, user_id, type, content, is_read, created_at)
SELECT id, user_id, type, content, is_read, created_at FROM notifications;

DROP TABLE notifications;

ALTER TABLE notifications_backup RENAME TO notifications;