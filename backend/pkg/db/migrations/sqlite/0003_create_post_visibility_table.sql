-- +goose Up
CREATE TABLE post_visibility (
    post_id INTEGER NOT NULL,
    viewer_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, viewer_id),
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (viewer_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
--DROP TABLE post_visibility;
