CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_role_check CHECK (role IN ('admin', 'user'))
);

CREATE TABLE keyboard_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    symbols TEXT NOT NULL
);

CREATE TABLE difficulty_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    keyboard_zone_id UUID NOT NULL,
    allowed_mistakes INT NOT NULL,
    key_press_time NUMERIC NOT NULL,
    min_exercise_length INT NOT NULL,
    max_exercise_length INT NOT NULL,

    CONSTRAINT fk_difficulty_zone
        FOREIGN KEY (keyboard_zone_id)
        REFERENCES keyboard_zones(id)
        ON DELETE CASCADE
);

CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text TEXT NOT NULL,
    level_id UUID NOT NULL,

    CONSTRAINT fk_exercise_level
        FOREIGN KEY (level_id)
        REFERENCES difficulty_levels(id)
        ON DELETE CASCADE
);

CREATE TABLE statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    level_id UUID NOT NULL,
    exercise_id UUID NOT NULL,
    mistakes_percent NUMERIC NOT NULL,
    execution_time NUMERIC NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    speed NUMERIC NOT NULL,

    CONSTRAINT fk_stat_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_stat_level
        FOREIGN KEY (level_id)
        REFERENCES difficulty_levels(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_stat_exercise
        FOREIGN KEY (exercise_id)
        REFERENCES exercises(id)
        ON DELETE CASCADE
);
