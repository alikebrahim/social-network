-- +goose Up
CREATE TABLE event_responses (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    response TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE event_responses;
