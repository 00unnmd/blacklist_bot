package database

import (
	"blacklist_bot/internal/models"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Database struct {
	db *sql.DB
}

func loadDBConfig() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable client_encoding=UTF8",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS banned_users (
			id SERIAL PRIMARY KEY,
			phone_number TEXT UNIQUE NOT NULL,
			full_name TEXT NOT NULL,
			description TEXT NOT NULL,
			birth_day TEXT,
			city TEXT,
			school_format TEXT
		);
		
		CREATE TABLE IF NOT EXISTS appeals (
			id SERIAL PRIMARY KEY,
			question TEXT NOT NULL,
			initiator TEXT NOT NULL,
			is_answered BOOLEAN NOT NULL
		);
	`)

	return err
}

func New() (*Database, error) {
	connStr := loadDBConfig()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) AddBannedUser(user models.BannedUser) error {
	_, err := d.db.Exec(
		"INSERT INTO banned_users (phone_number, full_name, description, birth_day, city, school_format) VALUES ($1, $2, $3, $4, $5, $6)",
		user.PhoneNumber, user.FullName, user.Description, user.BirthDay, user.City, user.SchoolFormat,
	)

	return err
}

func (d *Database) FindBannedUser(phoneNumber string) (*models.BannedUser, error) {
	row := d.db.QueryRow(
		"SELECT id, phone_number, full_name, description FROM banned_users WHERE phone_number = $1",
		phoneNumber,
	)

	var user models.BannedUser
	err := row.Scan(&user.ID, &user.PhoneNumber, &user.FullName, &user.Description)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (d *Database) AddAppeal(appeal models.Appeal) error {
	_, err := d.db.Exec(
		"INSERT INTO appeals (question, initiator, is_answered) VALUES ($1, $2, $3)",
		appeal.Question, appeal.Initiator, appeal.IsAnswered,
	)

	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetUnansweredAppeals() ([]models.Appeal, error) {
	rows, err := d.db.Query(
		"SELECT id, question, initiator FROM appeals WHERE is_answered = FALSE",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appeals []models.Appeal
	for rows.Next() {
		var a models.Appeal
		if err := rows.Scan(&a.ID, &a.Question, &a.Initiator); err != nil {
			return nil, err
		}
		appeals = append(appeals, a)
	}

	return appeals, nil
}

func (d *Database) MarkAppealAsAnswered(id int) error {
	_, err := d.db.Exec(
		"UPDATE appeals SET is_answered = TRUE WHERE id = $1",
		id,
	)

	return err
}
