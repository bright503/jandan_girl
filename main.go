package main

import (
	"log"
	"time"
)

func main() {
	_ = InitDB()

	log.Printf("下载数据库中现有图片...")
	posts := SelectAllPost()
	log.Printf("查询到%d条文章", len(posts))
	DownloadPosts(posts)
	log.Printf("下载完成")

	for {
		refresh()
		time.Sleep(300 * time.Second)
	}
}

func refresh() {
	startId := ""
	for i := 0; i < 10; i++ {
		log.Printf("查询 start_id: %s", startId)
		response := getPosts(startId)
		n := len(response.Data)
		log.Printf("%d条", n)
		if n <= 0 {
			break
		}
		startId = replaceInto(response)
		time.Sleep(1 * time.Second)
	}
	log.Printf("更新成功!")
}
