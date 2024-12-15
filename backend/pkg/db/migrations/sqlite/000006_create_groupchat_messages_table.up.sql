CREATE TABLE groupchat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    groupchat_id INTEGER NOT NULL,
    sender_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    image BLOB,
    FOREIGN KEY (groupchat_id) REFERENCES groupchat(id),
    FOREIGN KEY (sender_id) REFERENCES users(id)
);