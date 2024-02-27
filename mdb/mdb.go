package mdb

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
)

type EmailEntry struct {
	Id          int64
	Email       string
	ConfirmedAt *time.Time
	OptOut      bool
}

func TryCreate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS emails (
			id INTEGER PRIMARY KEY,
			email TEXT UNIQUE,
			confirmed_at INTEGER,
			opt_out BOOL
		);
	`)
	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			// code = 1 - table is already exists
			if sqlError.Code != 1 {
				log.Fatal(sqlError)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func emailEntryFromRow(row *sql.Rows) (*EmailEntry, error) {
	var id int64
	var email string
	var confirmedAt int64
	var optOut bool

	err := row.Scan(&id, &email, &confirmedAt, &optOut)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	t := time.Unix(confirmedAt, 0)
	return &EmailEntry{
		Id:          id,
		Email:       email,
		ConfirmedAt: &t,
		OptOut:      optOut,
	}, nil
}

func CreateEmailEntry(db *sql.DB, email string) error {
	_, err := db.Exec(`INSERT INTO 
			emails (email, confirmed_at, opt_out) 
			VALUES (?, 0, false);`, email)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetEmailEntry(db *sql.DB, email string) (*EmailEntry, error) {
	rows, err := db.Query(`
		SELECT id, email, confirmed_at, opt_out
		FROM emails 
		WHERE email = ?;
	`, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return emailEntryFromRow(rows) //taking first matched row as email should unique
	}
	return nil, nil
}

func UpdateEmailEntry(db *sql.DB, emailEntry *EmailEntry) error {
	if emailEntry.ConfirmedAt == nil || emailEntry.OptOut {
		errMsg := "confirmedAt and optOut fields should not be empty"
		log.Printf("%s %v", errMsg, emailEntry)
		return errors.New(errMsg)
	}
	t := emailEntry.ConfirmedAt.Unix()
	_, err := db.Exec(`
		INSERT INTO 
		emails (email, confirmed_at, opt_out)
		VALUES (?, ?, ?)
		ON CONFLICT (email) 
		DO UPDATE SET 
		confirmed_at = ?, opt_out = ?;
	`, emailEntry.Email, t, emailEntry.OptOut, t, emailEntry.OptOut)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteEmailEntry(db *sql.DB, email string) error {
	_, err := db.Exec(`
		UPDATE emails
		SET opt_out=true
		WHERE email = ?;
	`, email)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type GetEmailBatchQueryParams struct {
	Page  int
	Count int
}

func GetEmailBatch(db *sql.DB, params GetEmailBatchQueryParams) ([]EmailEntry, error) {
	rows, err := db.Query(`
		SELECT id, email, confirmed_at, opt_out
		FROM emails
		WHERE opt_out=false
		ORDER BY id ASC
		LIMIT ? OFFSET ?;
	`, params.Count, (params.Page-1)*params.Count)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	emailEntries := make([]EmailEntry, 0, params.Count)
	for rows.Next() {
		emailEntry, err := emailEntryFromRow(rows)
		if err != nil {
			return nil, err
		}
		emailEntries = append(emailEntries, *emailEntry)
	}
	return emailEntries, nil
}
