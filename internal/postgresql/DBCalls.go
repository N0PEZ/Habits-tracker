package postgresql

import (
	"context"
	"errors"
	"fmt"
	"huibitica/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUser(user models.RegisterUserRequest, conn *pgxpool.Pool) error {
	tx, err := conn.Begin(context.Background()) // Начинаем транзакцию
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background()) // Откатываем в случае ошибки

	// 1. Вставляем пользователя
	var userID int
	err = tx.QueryRow(context.Background(),
		`INSERT INTO users (username, email, phone, created_at)
         VALUES ($1, $2, $3, $4)
         RETURNING user_id`,
		user.Username,
		user.Email,
		user.Phone,
		time.Now(),
	).Scan(&userID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "users_username_key":
				return fmt.Errorf("username already exists")
			case "users_email_key":
				return fmt.Errorf("email already exists")
			}
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// 2. Вставляем пароль
	_, err = tx.Exec(context.Background(),
		`INSERT INTO passwords (user_id, username, password)
         VALUES ($1, $2, $3)`,
		userID,
		user.Username,
		user.Password,
	)

	if err != nil {
		return fmt.Errorf("failed to insert password: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func AddHabit(habit models.Habit, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`INSERT INTO habits (
			user_id, text, note, good, bad, difficulty,
			count_reset_after, good_count, bad_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		habit.UserID,
		habit.Text,
		habit.Note,
		habit.Good,
		habit.Bad,
		habit.Difficulty,
		habit.CountResetAfter,
		habit.GoodCount,
		habit.BadCount,
	)
	if err != nil {
		return fmt.Errorf("failed to insert habit: %w", err)
	}
	return nil
}

func AddDaily(daily models.Daily, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`INSERT INTO dailies (
			user_id, text, note, difficulty, start_date,
			repeat_every, repeat_every_x, dayweeks, streak)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		daily.UserID,
		daily.Text,
		daily.Note,
		daily.Difficulty,
		daily.StartDate,
		daily.RepeatEvery,
		daily.RepeatEveryX,
		daily.DayWeeks,
		daily.Streak,
	)
	if err != nil {
		return fmt.Errorf("failed to insert daily: %w", err)
	}
	return nil
}

func AddTask(task models.Task, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`INSERT INTO tasks (
			user_id, name, note, difficulty, deadline)
		VALUES ($1, $2, $3, $4, $5)`,
		task.UserID,
		task.Name,
		task.Note,
		task.Difficulty,
		task.Deadline,
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}
	return nil
}

func EditUserUsername(r models.EditUserData, conn *pgxpool.Pool) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`UPDATE users
		SET username = $1
		WHERE user_id = $2`,
		r.NewString,
		r.UserID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_key" {
				return fmt.Errorf("username already exists")
			}
		}
		return fmt.Errorf("failed to update username: %w", err)
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE passwords
		SET username = $1
		WHERE user_id = $2`,
		r.NewString,
		r.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func EditUserEmail(r models.EditUserData, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`UPDATE users
		SET email = $1
		WHERE user_id = $2`,
		r.NewString,
		r.UserID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				return fmt.Errorf("account with this email already exists")
			}
		}
		return fmt.Errorf("failed to update username: %w", err)
	}
	return nil
}

func EditUserPhone(r models.EditUserData, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`UPDATE users
		SET phone = $1
		WHERE user_id = $2`,
		r.NewString,
		r.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update phone: %w", err)
	}
	return nil
}

func EditPassword(password models.Password, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`UPDATE passwords
		SET password = $1
		WHERE user_id = $2`,
		password.Password,
		password.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}
func EditHabit(habit models.Habit, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`UPDATE habits
		SET text = $1, note = $2, good = $3, bad = $4,
			difficulty = $5, count_reset_after = $6,
			good_count = $7, bad_count = $8
		WHERE id = $9`,
		habit.Text,
		habit.Note,
		habit.Good,
		habit.Bad,
		habit.Difficulty,
		habit.CountResetAfter,
		habit.GoodCount,
		habit.BadCount,
		habit.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}
	return nil
}

func EditDaily(daily models.Daily, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`UPDATE dailies
		SET text = $1, note = $2, difficulty = $3,
			start_date = $4, repeat_every = $5,
			repeat_every_x = $6, dayweeks = $7,
			streak = $8
		WHERE id = $9`,
		daily.Text,
		daily.Note,
		daily.Difficulty,
		daily.StartDate,
		daily.RepeatEvery,
		daily.RepeatEveryX,
		daily.DayWeeks,
		daily.Streak,
		daily.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update daily: %w", err)
	}
	return nil
}

func EditTask(task models.Task, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`UPDATE tasks
		SET name = $1, note = $2, difficulty = $3, deadline = $4
		WHERE id = $5`,
		task.Name,
		task.Note,
		task.Difficulty,
		task.Deadline,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

func DeleteUser(userID int, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`DELETE FROM users
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func DeleteHabit(id int, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`DELETE FROM habits
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}
	return nil
}

func DeleteDaily(id int, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`DELETE FROM dailies
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete daily: %w", err)
	}
	return nil
}

func DeleteTask(id int, conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(),
		`DELETE FROM tasks
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

func GetUserByID(userID int, conn *pgxpool.Pool) (*models.User, error) {
	var user models.User
	err := conn.QueryRow(context.Background(),
		`SELECT user_id, username, email, phone
		FROM users
		WHERE user_id = $1`,
		userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Phone,
	)

	if err != nil {
		return &models.User{}, err
	}
	return &models.User{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

func GetUserByUsername(username string, conn *pgxpool.Pool) (*models.User, error) {
	var user models.User
	err := conn.QueryRow(context.Background(),
		`SELECT user_id, username, email, phone
		FROM users
		WHERE username = $1`,
		username).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Phone,
	)

	if err != nil {
		return &models.User{}, err
	}

	return &models.User{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}, nil
}

func GetUserByEmail(email string, conn *pgxpool.Pool) (*models.User, error) {
	var user models.User
	err := conn.QueryRow(context.Background(),
		`SELECT user_id, username, email, phone
		FROM users
		WHERE email = $1`,
		email).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Phone,
	)

	if err != nil {
		return &models.User{}, err
	}

	return &models.User{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}, nil
}

func GetHabits(userID int, conn *pgxpool.Pool) ([]models.Habit, error) {
	var habits []models.Habit
	rows, err := conn.Query(context.Background(),
		`SELECT id, user_id, text, note, good, bad,
			difficulty, count_reset_after, good_count, bad_count
		FROM habits
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var habit models.Habit
		err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Text,
			&habit.Note,
			&habit.Good,
			&habit.Bad,
			&habit.Difficulty,
			&habit.CountResetAfter,
			&habit.GoodCount,
			&habit.BadCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}
		habits = append(habits, habit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return habits, nil
}

func GetDailies(userID int, conn *pgxpool.Pool) ([]models.Daily, error) {
	var dailies []models.Daily
	rows, err := conn.Query(context.Background(),
		`SELECT id, user_id, text, note, difficulty,
			start_date, repeat_every, repeat_every_x,
			dayweeks, streak
		FROM dailies
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get dailies: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var daily models.Daily
		err := rows.Scan(
			&daily.ID,
			&daily.UserID,
			&daily.Text,
			&daily.Note,
			&daily.Difficulty,
			&daily.StartDate,
			&daily.RepeatEvery,
			&daily.RepeatEveryX,
			&daily.DayWeeks,
			&daily.Streak,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan daily: %w", err)
		}
		dailies = append(dailies, daily)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return dailies, nil
}

func GetTasks(userID int, conn *pgxpool.Pool) ([]models.Task, error) {
	var tasks []models.Task
	rows, err := conn.Query(context.Background(),
		`SELECT id, user_id, name, note, difficulty,
			deadline
		FROM tasks
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Name,
			&task.Note,
			&task.Difficulty,
			&task.Deadline,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tasks, nil
}
