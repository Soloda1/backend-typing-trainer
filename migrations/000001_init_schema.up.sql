-- Расширение для генерации UUID (если не включено)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Пользователи
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_role_check CHECK (role IN ('admin', 'user'))
);

-- Клавиатурные зоны
CREATE TABLE keyboard_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    symbols TEXT NOT NULL
);

-- Уровни сложности (теперь без привязки к одной зоне)
CREATE TABLE difficulty_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    allowed_mistakes INT NOT NULL,
    key_press_time NUMERIC NOT NULL,
    min_exercise_length INT NOT NULL,
    max_exercise_length INT NOT NULL
);

-- Связующая таблица: уровни <-> зоны (многие ко многим)
CREATE TABLE level_keyboard_zones (
    level_id UUID NOT NULL,
    keyboard_zone_id UUID NOT NULL,
    PRIMARY KEY (level_id, keyboard_zone_id),
    FOREIGN KEY (level_id) REFERENCES difficulty_levels(id) ON DELETE CASCADE,
    FOREIGN KEY (keyboard_zone_id) REFERENCES keyboard_zones(id) ON DELETE CASCADE
);

-- Упражнения (ссылаются на уровень)
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text TEXT NOT NULL,
    level_id UUID NOT NULL,
    CONSTRAINT fk_exercise_level
        FOREIGN KEY (level_id)
        REFERENCES difficulty_levels(id)
        ON DELETE CASCADE
);

-- Статистика (ссылается на пользователя, уровень, упражнение)
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