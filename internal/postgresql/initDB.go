package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func InitDB(connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	// Сначала создаем таблицу users, так как другие таблицы ссылаются на нее
	queries := map[string]string{
		"users": `CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL RIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(20),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			)`,

		"passwords": `CREATE TABLE IF NOT EXISTS passwords (
			user_id INTEGER PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			CONSTRAINT fk_passwords_user 
				FOREIGN KEY(user_id) 
				REFERENCES users(user_id)
				ON DELETE CASCADE)`,

		"habits": `CREATE TABLE IF NOT EXISTS habits (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			text VARCHAR(63) NOT NULL,
			note VARCHAR(255),
			good BOOLEAN DEFAULT TRUE NOT NULL,
			bad BOOLEAN DEFAULT FALSE NOT NULL,
			difficulty INT NOT NULL,
			count_reset_after INT DEFAULT 0 NOT NULL,
			good_count INT DEFAULT 0 NOT NULL,
			bad_count INT DEFAULT 0 NOT NULL,
			CONSTRAINT fk_habits_user 
				FOREIGN KEY(user_id) 
				REFERENCES users(user_id)
				ON DELETE CASCADE)`,

		"dailies": `CREATE TABLE IF NOT EXISTS dailies (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			text VARCHAR(63) NOT NULL,
			note VARCHAR(255),
			difficulty INT NOT NULL,
			start_date DATE NOT NULL,
			repeat_every INT DEFAULT 0 NOT NULL,
			repeat_every_x INT NOT NULL,
			dayweeks VARCHAR(14) DEFAULT NULL,
			streak INT DEFAULT 0 NOT NULL,
			CONSTRAINT fk_dailies_user 
				FOREIGN KEY(user_id) 
				REFERENCES users(user_id)
				ON DELETE CASCADE)`,

		"tasks": `CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			name VARCHAR(63) NOT NULL,
			note VARCHAR(255),
			difficulty INT NOT NULL,
			deadline DATE NOT NULL,
			CONSTRAINT fk_tasks_user 
				FOREIGN KEY(user_id) 
				REFERENCES users(user_id)
				ON DELETE CASCADE)`,
	}

	// Создаем таблицы в правильном порядке
	creationOrder := []string{"users", "passwords", "habits", "dailies", "tasks"}
	for _, table := range creationOrder {
		query := queries[table]
		_, err = conn.Exec(context.Background(), query)
		if err != nil {
			return nil, fmt.Errorf("unable to create table %s: %v", table, err)
		}
	}

	log.Println("Database initialized successfully with foreign key constraints")
	return conn, nil
}
