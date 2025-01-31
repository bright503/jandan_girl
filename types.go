package main

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []Post `json:"data"`
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
