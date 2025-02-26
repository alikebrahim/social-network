-- +goose Up
CREATE TABLE followers (
    id SERIAL PRIMARY KEY,
    follower_id INTEGER NOT NULL,
    followed_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users (id),
    FOREIGN KEY (followed_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE followers;
