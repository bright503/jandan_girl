package main

import (
	"log"
	"strconv"
	"time"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []Post `json:"data"`
}

func (res Response) PrintLog() string {
	n := len(res.Data)

	log.Printf("%d条", n)
	if n <= 0 {
		return ""
	}

	startTime, _ := time.Parse(time.RFC3339, res.Data[n-1].Date)
	start := startTime.Format("2006-01-02 15:04:05")
	endTime, _ := time.Parse(time.RFC3339, res.Data[0].Date)
	end := endTime.Format("2006-01-02 15:04:05")

	log.Printf("保存文章 %d 条, %s-%s", n, start, end)
	return strconv.Itoa(res.Data[n-1].ID)
}

type Post struct {
	ID           int     `json:"id"`
	PostID       int     `json:"post_id"`
	Author       string  `json:"author"`
	AuthorType   int     `json:"author_type"`
	Date         string  `json:"date"`
	DateUnix     int64   `json:"date_unix"`
	Content      string  `json:"content"`
	UserID       int     `json:"user_id"`
	VotePositive int     `json:"vote_positive"`
	VoteNegative int     `json:"vote_negative"`
	Images       []Image `json:"images"`
	IPLocation   string  `json:"ip_location"`
}
type Image struct {
	URL       string `json:"url"`
	FullURL   string `json:"full_url"`
	Host      string `json:"host"`
	ThumbSize string `json:"thumb_size"`
	Ext       string `json:"ext"`
	FileName  string `json:"file_name"`
	Path      string
}
