package queries

const CreatePasswordRecoveriesTableSQL = `
CREATE TABLE IF NOT EXISTS password_recoveries (
    id TEXT PRIMARY KEY,
    token TEXT UNIQUE NOT NULL,
    user_id TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_recovery_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_password_recoveries_token ON password_recoveries (token);
CREATE INDEX IF NOT EXISTS idx_password_recoveries_user_id ON password_recoveries (user_id);
`
