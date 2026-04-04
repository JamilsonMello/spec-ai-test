-- UP
ALTER TABLE posts ADD COLUMN video_url TEXT;

-- DOWN
ALTER TABLE posts DROP COLUMN video_url;
