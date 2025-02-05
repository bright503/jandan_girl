package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"
)

const PageSize = 300

func init() {
	_ = InitDB()
}

func main() {
	log.Printf("下载数据库中现有图片...")
	for i := 0; true; i++ {
		posts := SelectPostByPage(i*PageSize, PageSize)
		log.Printf("查询到%d条文章", len(posts))
		DownloadPosts(posts)
		if len(posts) < PageSize {
			break
		}
	}
	log.Printf("下载数据库中图片完成！")
	Refresh(false)
	// 启动后每天全量刷新 更新XX和OO
	fullRefreshTicker := time.NewTicker(1 * time.Hour)
	// 半分钟检查最新的图片
	checkTicker := time.NewTicker(20 * time.Second)
	defer fullRefreshTicker.Stop()
	defer checkTicker.Stop()

	for {
		select {
		case <-fullRefreshTicker.C:
			Refresh(false)
		case <-checkTicker.C:
			Refresh(true)
		}
	}
}

// Refresh 更新OO和XX
// check 检查到已下载的图片时就停止
func Refresh(check bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic: %v\nStack Trace:\n%s\n", err, debug.Stack())
		}
	}()
	if !check {
		log.Printf("全量更新煎蛋数据...")
	}
	startId := ""
	for {
		response, err := GetPosts(startId)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		startId = response.PrintLog()
		if startId == "" {
			break
		}
		ReplaceInto(response.Data)
		if DownloadPosts(response.Data) && check {
			log.Printf("检查到已下载图片，更新完成")
			return
		}
		time.Sleep(1 * time.Second)
	}
	log.Printf("更新成功!")
}
