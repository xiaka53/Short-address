package main

import "os"

func main() {
	setEnv()
	a := App{}
	a.Initialize(getEnv())
	a.Run(":8000")
}

func setEnv() {
	_ = os.Setenv("RedisAddr", "127.0.0.1:6379")
	_ = os.Setenv("RedisPwd", "")
	_ = os.Setenv("RedisDb", "1")
	_ = os.Setenv("MaxActive", "2000")
	_ = os.Setenv("MaxIdle", "2000")
}