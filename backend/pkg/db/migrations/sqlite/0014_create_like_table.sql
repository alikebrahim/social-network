-- +goose Up
CREATE TABLE likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    is_like BOOLEAN NOT NULL, 
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);


-- +goose Down
DROP TABLE likes;
