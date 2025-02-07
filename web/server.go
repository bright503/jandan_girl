package web

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"jandan_girl/db"
	"net/http"
	"strconv"
)

const pageSize = 5

//go:embed templates/*
var templatesFS embed.FS // 用于嵌入 templates 目录下的所有文件
func StartServer() {
	r := gin.Default()

	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(templatesFS, "templates/*")))

	// 首页路由，渲染页面
	r.GET("/", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		posts := db.SelectPostByPage((page-1)*pageSize, pageSize)

		totalPosts := db.GetTotalPostsCount() // 你需要根据数据库实际情况来获取总文章数
		totalPages := (totalPosts + pageSize - 1) / pageSize

		prevPage := page
		nextPage := page

		if page > 1 {
			prevPage = page - 1
		}
		if page < totalPages {
			nextPage = page + 1
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Posts":    posts,
			"PrevPage": prevPage,
			"NextPage": nextPage,
		})
	})

	// 图片文件服务
	r.Static("/img", "data/img")
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
