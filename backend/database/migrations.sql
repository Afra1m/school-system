-- Создание таблицы subjects с ограничением уникальности на поле name
CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    teacher_id INTEGER REFERENCES teachers(id)
);

-- Добавление индекса для ускорения поиска по имени
CREATE INDEX IF NOT EXISTS idx_subjects_name ON subjects(name); 