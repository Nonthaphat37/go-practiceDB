package main

import (
	"database/sql"
)

type user struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

const (
	querySql = "SELECT firstname, lastname FROM users WHERE id=$1";
)

func (u *user) getUser(db *sql.DB) error{
	err := db.QueryRow(querySql, u.ID).Scan(&u.Firstname, &u.Lastname);
	return err;
}