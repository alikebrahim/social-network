CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    type TEXT NOT NULL CHECK (type IN ('public', 'almost private', 'private')),
    image BLOB,
    FOREIGN KEY (user_id) REFERENCES users(id)
);