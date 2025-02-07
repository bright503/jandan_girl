package main

import (
	"jandan_girl/db"
	"jandan_girl/jandan"
)

func init() {
	_ = db.InitDB()
}

func main() {
	go jandan.Run()
	select {}
}
