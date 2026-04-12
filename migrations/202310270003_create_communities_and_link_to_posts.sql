-- UP
CREATE TABLE communities (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE posts ADD COLUMN community_id UUID REFERENCES communities(id) ON DELETE SET NULL;

-- DOWN
ALTER TABLE posts DROP COLUMN IF EXISTS community_id;
DROP TABLE IF EXISTS communities;
