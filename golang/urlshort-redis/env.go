package main

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	S Storage
}

func getEnv() *Env {
	addr := os.Getenv("APP_REDIS_ADDR")
	if addr == "" {
		addr = "locahost:6379"
	}

	passwd := os.Getenv("APP_REDIS_PASSWD")
	if passwd == "" {
		passwd = ""
	}

	dbS := os.Getenv("APP_REDIS_DB")
	if dbS == "" {
		dbS = ""
	}
	db, err := strconv.Atoi(dbS)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connect to redis (addr: %s password: %s db: %d", addr, passwd, db)
	r := NewRedicCli(addr, passwd, db)
	return &Env{S: Storage(r)}

}
