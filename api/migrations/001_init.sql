-- +goose Up

-- users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(50),
    level INT DEFAULT 1,
    experience INT DEFAULT 0,
    gold INT DEFAULT 0,
    seeds INT DEFAULT 0,
    streak INT DEFAULT 0,
    growth_stage INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- categories
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    name VARCHAR(50),
    color VARCHAR(20),
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- tasks
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(20) DEFAULT 'task',
    title TEXT,
    description TEXT,
    importance INT,
    category_id INT REFERENCES categories(id),
    due_date DATE,
    repeat_interval VARCHAR(20),
    reminder_time TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    xp_reward INT DEFAULT 0,
    gold_reward INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- progress_log
CREATE TABLE IF NOT EXISTS progress_log (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    task_id INT REFERENCES tasks(id),
    xp_earned INT DEFAULT 0,
    gold_earned INT DEFAULT 0,
    logged_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- farms
CREATE TABLE IF NOT EXISTS farms (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    season VARCHAR(10) DEFAULT 'spring',
    is_day BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- plants
CREATE TABLE IF NOT EXISTS plants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    image_url TEXT,
    level_required INT DEFAULT 1,
    growth_modifier FLOAT DEFAULT 1.0,
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- user_plants
CREATE TABLE IF NOT EXISTS user_plants (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    plant_id INT REFERENCES plants(id),
    status VARCHAR(20) DEFAULT 'growing',
    planted_at TIMESTAMP DEFAULT now(),
    growth_progress INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- inventory_items
CREATE TABLE IF NOT EXISTS inventory_items (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    item_type VARCHAR(20),
    item_name VARCHAR(100),
    quantity INT DEFAULT 1,
    is_equipped BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- reminders
CREATE TABLE IF NOT EXISTS reminders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    task_id INT REFERENCES tasks(id),
    remind_at TIMESTAMP,
    is_sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- onboarding
CREATE TABLE IF NOT EXISTS onboarding (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    step VARCHAR(50),
    is_completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    created_by INT,
    updated_at TIMESTAMP,
    updated_by INT,
    deleted_at TIMESTAMP,
    deleted_by INT
);

-- +goose Down

DROP TABLE IF EXISTS onboarding;
DROP TABLE IF EXISTS reminders;
DROP TABLE IF EXISTS inventory_items;
DROP TABLE IF EXISTS user_plants;
DROP TABLE IF EXISTS plants;
DROP TABLE IF EXISTS farms;
DROP TABLE IF EXISTS progress_log;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
