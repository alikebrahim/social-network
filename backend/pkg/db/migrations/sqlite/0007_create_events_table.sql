-- +goose Up
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL,
    creator_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    event_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (creator_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE events;
