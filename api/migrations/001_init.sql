-- +goose Up

-- user_info (основная таблица пользователей)
CREATE TABLE IF NOT EXISTS user_info (
    id BIGSERIAL PRIMARY KEY,
    max_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    last_login TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- user_stat (статистика и прогресс)
CREATE TABLE IF NOT EXISTS user_stat (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE UNIQUE,
    experience BIGINT DEFAULT 0,
    gold BIGINT DEFAULT 0,
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    total_plant_harvested BIGINT DEFAULT 0,
    total_tasks_completed BIGINT DEFAULT 0,
    is_drought BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- seed (семена/растения)
CREATE TABLE IF NOT EXISTS seed (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    icon VARCHAR(100),
    level_required INT DEFAULT 1,
    target_growth INT NOT NULL,
    rarity VARCHAR(20) CHECK (rarity IN ('common', 'uncommon', 'rare', 'legendary', 'unique')) DEFAULT 'common',
    modification DECIMAL(5,2) DEFAULT 1.0,
    gold_reward INT DEFAULT 0,
    xp_reward INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- user_seed (семена пользователей)
CREATE TABLE IF NOT EXISTS user_seed (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE,
    seed_id INT NOT NULL REFERENCES seed(id),
    quantity BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, seed_id)
);

-- tag (категории/теги)
CREATE TABLE IF NOT EXISTS tag (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, name)
);

-- task (задачи)
CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    difficulty VARCHAR(20) CHECK (difficulty IN ('trifle', 'easy', 'normal', 'hard')) DEFAULT 'normal',
    tag_id INT REFERENCES tag(id),
    done BOOLEAN DEFAULT FALSE,
    date DATE,
    xp_reward INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- habit (привычки)
CREATE TABLE IF NOT EXISTS habit (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    difficulty VARCHAR(20) CHECK (difficulty IN ('trifle', 'easy', 'normal', 'hard')) DEFAULT 'normal',
    tag_id INT REFERENCES tag(id),
    done BOOLEAN DEFAULT FALSE,
    count INTEGER DEFAULT 0,
    period VARCHAR(20) CHECK (period IN ('day', 'week', 'month')),
    every INTEGER DEFAULT 1,
    start_date DATE,
    xp_reward INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- progress_log (лог прогресса)
CREATE TABLE IF NOT EXISTS progress_log (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE,
    task_id INT REFERENCES task(id),
    habit_id INT REFERENCES habit(id),
    xp_earned INT DEFAULT 0,
    gold_earned INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    CHECK (
        (task_id IS NOT NULL AND habit_id IS NULL) OR 
        (task_id IS NULL AND habit_id IS NOT NULL)
    )
);

-- bed (грядки пользователей)
CREATE TABLE IF NOT EXISTS bed (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE,
    cell_number INT NOT NULL CHECK (cell_number >= 1),
    is_locked BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, cell_number)
);

-- user_plant (растения пользователей)
CREATE TABLE IF NOT EXISTS user_plant (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE,
    seed_id INT NOT NULL REFERENCES seed(id),
    bed_id INT NOT NULL REFERENCES bed(id),
    current_growth INT DEFAULT 0,
    is_withered BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(bed_id) -- Одно растение на грядку
);


-- Таблица товаров магазина (индивидуальная для каждого пользователя)
CREATE TABLE IF NOT EXISTS good (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_info(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('seed', 'bed')),
    id_good INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    cost INTEGER NOT NULL CHECK (cost >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, type, id_good) -- Уникальность комбинации пользователь-тип-id
);


-- Создание индексов для производительности
CREATE INDEX IF NOT EXISTS idx_good_user_id ON good(user_id);
CREATE INDEX IF NOT EXISTS idx_good_type ON good(type);
CREATE INDEX IF NOT EXISTS idx_good_id_good ON good(id_good);
CREATE INDEX IF NOT EXISTS idx_good_user_type ON good(user_id, type);
CREATE INDEX IF NOT EXISTS idx_user_stat_user_id ON user_stat(user_id);
CREATE INDEX IF NOT EXISTS idx_user_seed_user_id ON user_seed(user_id);
CREATE INDEX IF NOT EXISTS idx_user_seed_seed_id ON user_seed(seed_id);
CREATE INDEX IF NOT EXISTS idx_bed_user_id ON bed(user_id);
CREATE INDEX IF NOT EXISTS idx_bed_cell_number ON bed(cell_number);
CREATE INDEX IF NOT EXISTS idx_bed_locked ON bed(is_locked);
CREATE INDEX IF NOT EXISTS idx_tag_user_id ON tag(user_id);
CREATE INDEX IF NOT EXISTS idx_task_user_id ON task(user_id);
CREATE INDEX IF NOT EXISTS idx_task_user_done ON task(user_id, done);
CREATE INDEX IF NOT EXISTS idx_task_date ON task(date);
CREATE INDEX IF NOT EXISTS idx_task_tag_id ON task(tag_id);
CREATE INDEX IF NOT EXISTS idx_habit_user_id ON habit(user_id);
CREATE INDEX IF NOT EXISTS idx_habit_user_done ON habit(user_id, done);
CREATE INDEX IF NOT EXISTS idx_habit_period ON habit(period);
CREATE INDEX IF NOT EXISTS idx_habit_start_date ON habit(start_date);
CREATE INDEX IF NOT EXISTS idx_habit_tag_id ON habit(tag_id);
CREATE INDEX IF NOT EXISTS idx_progress_log_user_id ON progress_log(user_id);
CREATE INDEX IF NOT EXISTS idx_progress_log_task_id ON progress_log(task_id);
CREATE INDEX IF NOT EXISTS idx_progress_log_habit_id ON progress_log(habit_id);
CREATE INDEX IF NOT EXISTS idx_progress_log_created_at ON progress_log(created_at);
CREATE INDEX IF NOT EXISTS idx_user_plant_user_id ON user_plant(user_id);
CREATE INDEX IF NOT EXISTS idx_user_plant_seed_id ON user_plant(seed_id);
CREATE INDEX IF NOT EXISTS idx_user_plant_bed_id ON user_plant(bed_id);

-- Дополнительные индексы для новых полей
CREATE INDEX IF NOT EXISTS idx_user_plant_withered ON user_plant(user_id, is_withered);
CREATE INDEX IF NOT EXISTS idx_progress_log_user_date ON progress_log(user_id, created_at);

-- +goose Down

-- Удаление индексов
DROP INDEX IF EXISTS idx_good_user_id;
DROP INDEX IF EXISTS idx_good_type;
DROP INDEX IF EXISTS idx_good_id_good;
DROP INDEX IF EXISTS idx_good_user_type;
DROP INDEX IF EXISTS idx_progress_log_user_date;
DROP INDEX IF EXISTS idx_user_plant_withered;
DROP INDEX IF EXISTS idx_user_plant_bed_id;
DROP INDEX IF EXISTS idx_user_plant_seed_id;
DROP INDEX IF EXISTS idx_user_plant_user_id;
DROP INDEX IF EXISTS idx_progress_log_created_at;
DROP INDEX IF EXISTS idx_progress_log_habit_id;
DROP INDEX IF EXISTS idx_progress_log_task_id;
DROP INDEX IF EXISTS idx_progress_log_user_id;
DROP INDEX IF EXISTS idx_habit_tag_id;
DROP INDEX IF EXISTS idx_habit_start_date;
DROP INDEX IF EXISTS idx_habit_period;
DROP INDEX IF EXISTS idx_habit_user_done;
DROP INDEX IF EXISTS idx_habit_user_id;
DROP INDEX IF EXISTS idx_task_tag_id;
DROP INDEX IF EXISTS idx_task_date;
DROP INDEX IF EXISTS idx_task_user_done;
DROP INDEX IF EXISTS idx_task_user_id;
DROP INDEX IF EXISTS idx_tag_user_id;
DROP INDEX IF EXISTS idx_bed_locked;
DROP INDEX IF EXISTS idx_bed_cell_number;
DROP INDEX IF EXISTS idx_bed_user_id;
DROP INDEX IF EXISTS idx_user_seed_seed_id;
DROP INDEX IF EXISTS idx_user_seed_user_id;
DROP INDEX IF EXISTS idx_user_stat_user_id;

-- Удаление таблиц в правильном порядке (из-за foreign keys)
DROP TABLE IF EXISTS good;
DROP TABLE IF EXISTS user_plant;
DROP TABLE IF EXISTS progress_log;
DROP TABLE IF EXISTS habit;
DROP TABLE IF EXISTS task;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS user_seed;
DROP TABLE IF EXISTS bed;
DROP TABLE IF EXISTS seed;
DROP TABLE IF EXISTS user_stat;
DROP TABLE IF EXISTS user_info;