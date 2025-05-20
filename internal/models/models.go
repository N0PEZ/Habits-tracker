package models

import (
	"time"
)

type RegisterUser struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone,omitempty" db:"phone"`
	Password string `json:"password" db:"password"`
}

type User struct {
	UserID    int       `json:"id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone,omitempty" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Password struct {
	ID       int    `json:"id" db:"id"`
	UserID   int    `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type Habit struct {
	ID              int    `json:"id" db:"id"`
	UserID          int    `json:"user_id" db:"user_id"`
	Text            string `json:"text" db:"text"`
	Note            string `json:"note,omitempty" db:"note"`
	Good            bool   `json:"good" db:"good"`
	Bad             bool   `json:"bad" db:"bad"`
	Difficulty      int    `json:"difficulty" db:"difficulty"`
	CountResetAfter int    `json:"count_reset_after" db:"count_reset_after"`
	GoodCount       int    `json:"good_count" db:"good_count"`
	BadCount        int    `json:"bad_count" db:"bad_count"`
}

type Daily struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	Text         string    `json:"text" db:"text"`
	Note         string    `json:"note,omitempty" db:"note"`
	Difficulty   int       `json:"difficulty" db:"difficulty"`
	StartDate    time.Time `json:"start_date" db:"start_date"`
	RepeatEvery  int       `json:"repeat_every" db:"repeat_every"`
	RepeatEveryX int       `json:"repeat_every_x" db:"repeat_every_x"`
	DayWeeks     string    `json:"day_weeks,omitempty" db:"dayweeks"`
	Streak       int       `json:"streak" db:"streak"`
}

type Task struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	Name       string    `json:"name" db:"name"`
	Note       string    `json:"note,omitempty" db:"note"`
	Difficulty int       `json:"difficulty" db:"difficulty"`
	Deadline   time.Time `json:"deadline" db:"deadline"`
	Completed  bool      `json:"completed" db:"completed"`
}
