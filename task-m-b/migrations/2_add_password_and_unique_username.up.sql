ALTER TABLE users
    ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL DEFAULT '';

-- Ensure usernames are unique for reliable login
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM   pg_constraint
        WHERE  conname = 'unique_user_name'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT unique_user_name UNIQUE (user_name);
    END IF;
END$$;

