-- +goose Up
CREATE TABLE group_chats (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL,
    sender_id INTEGER NOT NULL,
    content TEXT,
    image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (sender_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE group_chats;
