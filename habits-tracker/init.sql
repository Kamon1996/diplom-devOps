CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE habits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    frequency INT NOT NULL, -- планируемая частота выполнения в неделю (пример)
    target_percent INT NOT NULL, -- целевая частота выполнения в процентах (0-100)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE habit_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    habit_id UUID REFERENCES habits(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    done BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (habit_id, date)
);

CREATE TABLE advice (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL
);

INSERT INTO advice (message) VALUES
('Не сдавайся, привычки требуют времени!'),
('Каждый день — новый шанс стать лучше.'),
('Прогресс важнее перфекционизма.'),
('Возьми паузу, а потом продолжай.'),
('Ты на правильном пути!');