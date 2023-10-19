CREATE TABLE session (
    sid text PRIMARY KEY,
    user_id text not null,
    refresh_token text not null default '',
    is_valid boolean not null default false,
    created_at timestamptz,
    updated_at timestamptz,
    expired_at timestamptz,
    last_attempts_to_log_at timestamptz
);

-- for now, every check session assume attempt to login
CREATE OR REPLACE FUNCTION update_last_attempts_to_log()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE "session"
    SET last_accessed_at = NOW()
    WHERE id = NEW.id; -- Assuming you have a unique identifier like 'id'
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;