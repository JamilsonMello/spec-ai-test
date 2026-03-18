-- UP
ALTER TABLE users ADD COLUMN profile_picture_url TEXT;

-- DOWN
ALTER TABLE users DROP COLUMN profile_picture_url;
