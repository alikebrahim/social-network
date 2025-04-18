-- +goose Up
ALTER TABLE posts ADD COLUMN privacy_level TEXT NOT NULL DEFAULT 'public';
ALTER TABLE posts ADD COLUMN allowed_followers TEXT DEFAULT NULL;

-- +goose Down
CREATE TABLE posts_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO posts_backup (id, user_id, content, image, created_at, updated_at)
SELECT id, user_id, content, image, created_at, updated_at FROM posts;

DROP TABLE posts;

ALTER TABLE posts_backup RENAME TO posts;