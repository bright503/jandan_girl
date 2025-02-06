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

// DownloadPosts 下载多个Post
// return 是否有图片已经被下载
func DownloadPosts(posts []Post) bool {
	result := false
	for _, post := range posts {
		result = DownloadPost(post) || result
	}
	return result
}

// DownloadPost 下载单个Post
// return 是否有图片已经被下载
func DownloadPost(post Post) bool {
	result := false
	for _, img := range post.Images {
		parsedTime, _ := time.Parse(time.RFC3339, post.Date)
		result = downloadImageSetTime(img, parsedTime) || result
	}
	return result
}

// JanDanHost jandan图床
const (
	JanDanHost = "https://img.toto.im/"
	RegStr     = `https?://[^/]+/`
)

func downloadImageSetTime(image Image, time time.Time) bool {
	folderName := filepath.Join("data", "img", image.Path)
	err := os.MkdirAll(folderName, os.ModePerm)
	if err != nil {
		log.Printf("创建文件夹失败: %v", err)
		return false
	}
	fileName := image.FileName + "." + image.Ext
	savePath := filepath.Join(folderName, fileName)

	if pathExists(savePath) {
		return true
	}
	log.Printf("下载图片 %s", savePath)

	// 创建文件
	file, err := os.Create(savePath)
	if err != nil {
		log.Println(err)
		return false
	}
	defer file.Close()

	reg := regexp.MustCompile(RegStr)

	// 尺寸优先，新浪图床优先
	urls := []string{
		image.FullURL,
		reg.ReplaceAllString(image.FullURL, JanDanHost),
		image.URL,
		reg.ReplaceAllString(image.URL, JanDanHost),
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
	return false
}

func pathExists(path string) bool {
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
