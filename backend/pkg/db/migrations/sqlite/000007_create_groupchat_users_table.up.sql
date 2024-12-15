CREATE TABLE groupchat_users (
    groupchat_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (groupchat_id) REFERENCES groupchat(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (groupchat_id, user_id)
);