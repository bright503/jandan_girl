package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() (err error) {
	dbFile := "data/db/sqlite.db"
	if !pathExists(dbFile) {
		err := os.MkdirAll(path.Dir(dbFile), os.ModePerm)
		if err != nil {
			return err
		}
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("打开数据库错误: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("连接数据库错误: %v", err)
	}

	createTableSQL := `create table if not exists posts (
    	id            int          not null primary key,
    	post_id       int          null,
    	author        varchar(255) null,
    	author_type   int          null,
    	date          datetime     null,
    	content       text         null,
    	user_id       int          null,
    	vote_positive int          null,
    	vote_negative int          null,
    	ip_location   varchar(255) null
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("建表错误: %v", err)
	}

	createTableSQL = `create table if not exists images (
	    id         varchar(32)  not null primary key,
	    post_id    int          null,
	    url        varchar(255) null,
	    full_url   varchar(255) null,
	    host       varchar(255) null,
	    thumb_size varchar(50)  null,
	    ext        varchar(10)  null,
	    file_name  varchar(255) null,
	    path       varchar(64)  null
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("建表错误: %v", err)
	}

	return err
}

func SelectPostByPage(offset int, limit int) []Post {
	queryPost := `SELECT id ,post_id, author, author_type, date, content, user_id, vote_positive, vote_negative,
ip_location FROM posts order by date desc limit ?,?`
	stmt, _ := db.Prepare(queryPost)
	rows, _ := stmt.Query(offset, limit)
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		post := Post{}
		_ = rows.Scan(&post.ID, &post.PostID, &post.Author, &post.AuthorType, &post.Date, &post.Content, &post.UserID, &post.VotePositive, &post.VoteNegative, &post.IPLocation)
		post.Images = SelectImageByPostId(strconv.Itoa(post.ID))
		posts = append(posts, post)
	}
	return posts
}

func SelectAllImage() []Image {
	queryImg := `SELECT url, full_url, host, thumb_size, ext, file_name, path FROM images `
	rows, _ := db.Query(queryImg)
	defer rows.Close()
	var images []Image
	for rows.Next() {
		var img Image
		_ = rows.Scan(&img.URL, &img.FullURL, &img.Host, &img.ThumbSize, &img.Ext, &img.FileName, &img.Path)
		images = append(images, img)
	}
	return images
}

func SelectImageByPostId(postId string) []Image {
	queryImg := `SELECT url, full_url, host, thumb_size, ext, file_name, path FROM images where post_id = ?`
	stmt, _ := db.Prepare(queryImg)
	rows, _ := stmt.Query(postId)
	defer rows.Close()
	var images []Image
	for rows.Next() {
		var img Image
		_ = rows.Scan(&img.URL, &img.FullURL, &img.Host, &img.ThumbSize, &img.Ext, &img.FileName, &img.Path)
		images = append(images, img)
	}
	return images
}

func SelectPostById(postId string) Post {
	queryPost := `SELECT id ,post_id, author, author_type, date, content, user_id, vote_positive, vote_negative,
ip_location, FROM posts WHERE id = ?`
	stmt, _ := db.Prepare(queryPost)
	defer stmt.Close()
	row := stmt.QueryRow(postId)

	post := Post{}
	_ = row.Scan(&post.ID, &post.PostID, &post.Author, &post.AuthorType, &post.Date, &post.Content, &post.UserID, &post.VotePositive, &post.VoteNegative, &post.IPLocation)

	queryImg := `SELECT url, full_url, host, thumb_size, ext, file_name FROM images where post_id=` + postId

	rows, _ := db.Query(queryImg)
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		_ = rows.Scan(&img.URL, &img.FullURL, &img.Host, &img.ThumbSize, &img.Ext, &img.FileName)
		images = append(images, img)
	}
	post.Images = images
	return post
}

func ReplaceInto(posts []Post) {
	length := len(posts)
	if length <= 0 {
		return
	}

	postInsert, _ := db.Prepare(`REPLACE INTO posts (id, post_id, author, author_type, date, content, user_id, vote_positive, vote_negative, ip_location)
		              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	defer postInsert.Close()

	imgInsert, _ := db.Prepare(`REPLACE INTO images (id, post_id, url, full_url, host, thumb_size, ext, file_name, path)
			               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	defer imgInsert.Close()

	// Insert data into MySQL
	//for _, post := range response.Data {
	for i := length - 1; i >= 0; i-- {
		post := posts[i]

		// Insert post data
		_, err := postInsert.Exec(post.ID, post.PostID, post.Author, post.AuthorType, post.Date, post.Content, post.UserID, post.VotePositive, post.VoteNegative, post.IPLocation)
		if err != nil {
			log.Printf("插入文章错误: %v", err)
			continue
		}

		parsedTime, _ := time.Parse(time.RFC3339, post.Date)
		formattedDate := parsedTime.Format("20060102")
		// Insert images data
		for i, image := range post.Images {
			id := fmt.Sprintf("%d-%d", post.ID, i)
			post.Images[i].Path = formattedDate
			_, err := imgInsert.Exec(id, post.ID, image.URL, image.FullURL, image.Host, image.ThumbSize, image.Ext, image.FileName, formattedDate)
			if err != nil {
				log.Printf("插入图片错误 %v", err)
			}
		}
	}
}
