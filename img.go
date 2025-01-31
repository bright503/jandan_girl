package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func DownloadPosts(posts []Post) {
	for _, post := range posts {
		DownloadPost(post)
	}
}

func DownloadPost(post Post) {
	for _, img := range post.Images {
		parsedTime, _ := time.Parse(time.RFC3339, post.Date)
		downloadImageSetTime(img, parsedTime)
	}
}

func DownloadImages(images []Image) {
	for _, img := range images {
		downloadImage(img)
	}
}

func downloadImageSetTime(image Image, time time.Time) {

	folderName := filepath.Join("img", image.Path)
	err := os.MkdirAll(folderName, os.ModePerm)
	if err != nil {
		log.Printf("创建文件夹失败: %v", err)
		return
	}
	fileName := image.FileName + "." + image.Ext
	savePath := filepath.Join(folderName, fileName)

	if PathExists(savePath) {
		return
	}
	log.Printf("下载图片 %s", savePath)

	// 创建文件
	file, err := os.Create(savePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// jandan图床
	host := "https://img.toto.im/"

	reg := regexp.MustCompile(`https?://[a-zA-Z0-9._-]+/`)

	// 尺寸优先，新浪图床优先
	urls := []string{
		image.FullURL,
		reg.ReplaceAllString(image.FullURL, host),
		image.URL,
		reg.ReplaceAllString(image.URL, host),
	}

	for _, url := range urls {
		err := saveImage(url, file)
		if err != nil {
			log.Printf("%v 重新下载", err)
			continue
		}
		break
	}
	file.Close()
	if !time.IsZero() {
		_ = os.Chtimes(savePath, time, time)
	}
}

func downloadImage(image Image) {
	downloadImageSetTime(image, time.Time{})
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func saveImage(url string, file *os.File) error {

	c := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	log.Printf("下载: %s", url)
	resp, err := c.Get(url)
	if err != nil {
		return fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()
	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("响应错误: %s", resp.Status)
	}
	_, _ = io.Copy(file, resp.Body)
	return nil
}
