package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(dbAddress string, dbName string) (*pgxpool.Pool, error) {
	tempConn, err := pgx.Connect(context.Background(), dbAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	err = tempConn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	_, err = tempConn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s OWNER postgres", dbName))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if !(pgErr.Code == "42P04") {
				return nil, fmt.Errorf("ошибка создания БД: %v", err)
			}
		} else {
			return nil, fmt.Errorf("ошибка создания БД: %v", err)
		}
	}
	tempConn.Close(context.Background())

	pool, err := pgxpool.New(context.Background(), dbAddress+"/"+dbName)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s database: %v", dbName, err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping %s database: %v", dbName, err)
	}

	// Сначала создаем таблицу users, так как другие таблицы ссылаются на нее
	queries := map[string]string{
		"users": `CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(20),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
			difficulty INT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
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
			difficulty INT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
			start_date DATE NOT NULL,
			repeat_every INT DEFAULT 0 NOT NULL,
			repeat_every_x INT NOT NULL,
			dayweeks VARCHAR(32) DEFAULT NULL,
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
			difficulty INT NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
			deadline DATE NOT NULL,
			CONSTRAINT fk_tasks_user 
				FOREIGN KEY(user_id) 
				REFERENCES users(user_id)
				ON DELETE CASCADE)`,
	}

	creationOrder := [5]string{"users", "passwords", "habits", "dailies", "tasks"}
	for _, table := range creationOrder {
		query := queries[table]
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			return nil, fmt.Errorf("unable to create table %s: %v", table, err)
		}
	}

	return pool, nil
}
