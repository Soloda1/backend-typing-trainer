CREATE UNIQUE INDEX IF NOT EXISTS users_single_admin_idx
    ON users (role)
    WHERE role = 'admin';

INSERT INTO users (id, login, password_hash, role, created_at)
SELECT
    '00000000-0000-0000-0000-000000000001'::uuid,
    'admin',
    '$2a$10$.PkgksFS2VzhGAqjmIkXMOE9fYYlrtLaSm/4qGhCdVaO2lrzTvvoO',
    'admin',
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE role = 'admin'
);

