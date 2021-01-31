package main

import (
	// "fmt"
	"encoding/json"
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
	createUserCommand = "INSERT INTO users(id, firstname, lastname) VALUES($1, $2, $3) RETURNING id";
	updateUserCommand = "UPDATE users SET firstname=$1, lastname=$2 WHERE id=$3";
	CACHE_TTL = 0;
)

func (u *user) getUserRedis(Redis *redis.Client) (string, error){
	val, err := Redis.Get(strconv.Itoa(u.ID)).Result()
	return val, err
}

func (u *user) delUserRedis(Redis *redis.Client) error{
	err := Redis.Del(strconv.Itoa(u.ID)).Err();
	return err;
}

func (u *user) setUserRedis(Redis *redis.Client) error{
	json, _ := json.Marshal(u)
	err := Redis.Set(strconv.Itoa(u.ID), json, CACHE_TTL).Err()
	return err;
}

func (u *user) getUserDB(db *sql.DB) error{
	err := db.QueryRow(getUserCommand, u.ID).Scan(&u.Firstname, &u.Lastname);
	return err;
}

func (u *user) createUser(db *sql.DB) error{
	err := db.QueryRow(createUserCommand, u.ID, u.Firstname, u.Lastname).Scan(&u.ID)
	return err;
}

func (u *user) updateUser(db *sql.DB) error{
	_, err := db.Exec(updateUserCommand, u.Firstname, u.Lastname, u.ID)
	return err
}
