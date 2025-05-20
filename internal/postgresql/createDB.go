package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

const (
	dbName     = "huibitica"
	pgUser     = "postgres"
	pgPassword = "8968"
)

// InitDB инициализирует PostgreSQL и создаёт БД с таблицами
func CreateDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:8968@localhost:5432")
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %v", err)
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s OWNER %s", dbName, pgUser))

	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, fmt.Errorf("ошибка создания БД: %v", err)
	}
	fmt.Println("✅ База данных успешно инициализирована")
	return conn, nil
}
