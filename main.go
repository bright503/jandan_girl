package main

import (
	"jandan_girl/db"
	"jandan_girl/jandan"
	"jandan_girl/web"
)

func init() {
	_ = db.InitDB()
}

func main() {
	go jandan.Run()
	web.StartServer()
}
