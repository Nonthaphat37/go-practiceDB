package main

import (
	"strconv"
	"database/sql"
	"github.com/go-redis/redis"
)

type user struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

const (
	getUserCommand = "SELECT firstname, lastname FROM users WHERE id=$1";
)

func (u *user) getUserRedis(Redis *redis.Client) (string, error){
	val, err := Redis.Get(strconv.Itoa(u.ID)).Result()
	return val, err
}

func (u *user) getUserDB(db *sql.DB) error{
	err := db.QueryRow(getUserCommand, u.ID).Scan(&u.Firstname, &u.Lastname);
	return err;
}