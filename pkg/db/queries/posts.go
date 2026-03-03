package queries

const CreatePostsTableSQL = `
CREATE TABLE IF NOT EXISTS posts (
    id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    author_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_post_author
        FOREIGN KEY(author_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts (author_id);
`
