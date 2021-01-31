package main

import (
	"fmt"
	"errors"
	"database/sql"
)

type user struct {
	id          int     `json:"id"`
	firstname   string  `json:"firstname"`
	lastname    string  `json:"lastname"`
}

func (u *user) getUser(db *sql.DB) error{
	fmt.Println("abc");
	return errors.New("abc");
}