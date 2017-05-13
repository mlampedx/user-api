package main

import (
	"database/sql"
	"time"
)

type user struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Password     []byte    `json:"password"`
	Data         []string  `json:"data"`
	AdminGroupID int       `json:"adminGroupId"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
}

// To-Do: Implement authentication for db actions

func retrieveUsers(db *sql.DB, start, count int) ([]user, error) {

}

func (u *user) retrieveUser(db *sql.DB) error {
	return db.QueryRow("SELECT username, password, data, adminGroupId, dateCreated, dateModified FROM users WHERE id=$1",
		u.ID).Scan(&u.Username, &u.Password, &u.Data, &u.AdminGroupID, &u.DateCreated, &u.DateModified)
}

func (u *user) updateUser(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET username=$1, password=$2, data=$3, adminGroupId=$4, dateCreated=$5, dateModified=$6 FROM users WHERE id=$7",
		u.Username, u.Password, u.Data, u.AdminGroupID, u.DateCreated, u.DateModified, u.ID)

	return err
}

func (u *user) deleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)

	return err
}

func (u *user) createUser(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO users(username, password, data, adminGroupId, dateCreated, dateModified) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		u.Username, u.Password, u.Data, u.AdminGroupID, u.DateCreated, u.DateModified).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}
