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
	redisAddr := os.Getenv("RedisAddr")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPwd := os.Getenv("RedisPwd")
	if redisPwd == "" {
		redisPwd = ""
	}
	maxIdle := os.Getenv("MaxIdle")
	if maxIdle == "" {
		maxIdle = "2000"
	}
	idle, err := strconv.Atoi(maxIdle)
	if err != nil {
		log.Fatal(err)
	}
	maxActive := os.Getenv("MaxActive")
	if maxActive == "" {
		maxActive = "2000"
	}
	active, err := strconv.Atoi(maxActive)
	if err != nil {
		log.Fatal(err)
	}
	redisDb := os.Getenv("RedisDb")
	if redisDb == "" {
		redisDb = "0"
	}
	db, err := strconv.Atoi(redisDb)
	if err != nil {
		log.Fatal(err)
	}
	cli := NewRedisCli(redisAddr, redisPwd, idle, active, db)
	return &Env{S: cli}
}
