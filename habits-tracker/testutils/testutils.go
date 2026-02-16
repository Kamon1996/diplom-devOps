package testutils

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupTestDB() (*pgxpool.Pool, func()) {
	connStr := "postgres://postgres:postgres@localhost:5433/habit_test?sslmode=disable"

	var pool *pgxpool.Pool
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			break
		}
		log.Printf("Waiting for test DB: %v", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Unable to connect to test DB: %v", err)
	}

	if err := initTestDBSchema(ctx, pool); err != nil {
		log.Fatalf("Failed to init test DB schema: %v", err)
	}

	cleanup := func() {
		if _, err := pool.Exec(ctx, `
			TRUNCATE TABLE habit_records, habits, users RESTART IDENTITY CASCADE;
		`); err != nil {
			log.Printf("Failed to cleanup test DB: %v", err)
		}
		pool.Close()
	}

	return pool, cleanup
}

func initTestDBSchema(ctx context.Context, pool *pgxpool.Pool) error {
	schemaSQL := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(50) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS habits (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			description TEXT NOT NULL,
			frequency INT NOT NULL,
			target_percent INT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS habit_records (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			habit_id UUID REFERENCES habits(id) ON DELETE CASCADE,
			date DATE NOT NULL,
			done BOOLEAN NOT NULL DEFAULT FALSE,
			UNIQUE (habit_id, date)
		);

		CREATE TABLE IF NOT EXISTS advice (
			id SERIAL PRIMARY KEY,
			message TEXT NOT NULL
		);

		INSERT INTO advice (message) VALUES
		('Не сдавайся, привычки требуют времени!'),
		('Каждый день — новый шанс стать лучше.'),
		('Прогресс важнее перфекционизма.'),
		('Возьми паузу, а потом продолжай.'),
		('Ты на правильном пути!')
		ON CONFLICT DO NOTHING;
	`

	_, err := pool.Exec(ctx, schemaSQL)
	return err
}
