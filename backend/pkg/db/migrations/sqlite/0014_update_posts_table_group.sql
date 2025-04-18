-- +goose Up
ALTER TABLE posts ADD COLUMN group_id INTEGER DEFAULT NULL;
CREATE INDEX idx_posts_group_id ON posts(group_id);

-- +goose Down
DROP INDEX idx_posts_group_id;
CREATE TABLE posts_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    privacy_level TEXT NOT NULL DEFAULT 'public',
    allowed_followers TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO posts_backup (id, user_id, content, image, privacy_level, allowed_followers, created_at, updated_at)
SELECT id, user_id, content, image, privacy_level, allowed_followers, created_at, updated_at FROM posts;

DROP TABLE posts;

ALTER TABLE posts_backup RENAME TO posts;