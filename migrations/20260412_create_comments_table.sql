-- UP
CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_comment_post
        FOREIGN KEY(post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_comment_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

-- DOWN
DROP INDEX IF EXISTS idx_comments_post_id;
DROP TABLE IF EXISTS comments;
