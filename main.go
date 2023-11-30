package main

import (
	"blog/database"
	"blog/generator"
	"blog/post"
	"blog/web"
	"blog/web/components"
	"blog/web/pages"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.GetDatabaseConnection()
	defer db.Close()
	generator.SaveTemplComponent("index", pages.LandingPage())
	r := gin.Default()
	r.HTMLRender = &web.TemplRender{}

	username := os.Getenv("GOTH_USER")
	password := os.Getenv("GOTH_PASSWORD")

	if username == "" {
		log.Println("Environment Variable GOTH_USER not set.")
		log.Println("Using default user admin:admin")
		username = "admin"
		password = "admin"
	}

	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Next()
			return
		}
		if len(c.Request.URL.Path) == 1 {
			http.ServeFile(c.Writer, c.Request, "static/index.html")
			c.Abort()
			return
		}

		requestedFile := "static" + c.Request.URL.Path
		if _, err := os.Stat(requestedFile); err == nil {
			http.ServeFile(c.Writer, c.Request, requestedFile)
			c.Abort()
			return
		}
		c.Next()
	})
	api := r.Group("/api")

	authorizedApi := r.Group("/api", gin.BasicAuth(gin.Accounts{
		username: password,
	}))

	authorizedApi.POST("/posts", func(c *gin.Context) {
		content, _ := io.ReadAll(c.Request.Body)
		filename, err := post.GeneratePost(string(content), db)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Can't generate post! %v", err))
			return
		}
		message := fmt.Sprintf("Created new blog post %s", filename)
		log.Println(message)
		c.String(http.StatusOK, message)
	})

	authorizedApi.DELETE("/posts", func(c *gin.Context) {
		id := c.Query("id")
		err := post.DeletePost(id, db)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Can't delete post! %v", err))
			return
		}
		message := fmt.Sprintf("Deleted post %s", id)
		log.Println(message)
		c.String(http.StatusOK, message)
	})

	api.GET("/posts", func(c *gin.Context) {
		posts, err := post.GetPosts(db)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Can't load posts %v", err))
		}
		c.HTML(http.StatusOK, "", components.PostGrid(posts))
	})

	authorizedApi.POST("/asset", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["files"]
		paths := []string{}
		for _, file := range files {
			path := fmt.Sprintf("assets/%s", file.Filename)
			c.SaveUploadedFile(file, "static/"+path)
			paths = append(paths, path)
		}
		message := fmt.Sprintf("Files uploaded: %v", paths)
		log.Println(message)
		c.String(http.StatusOK, message)
	})

	authorizedApi.DELETE("/asset", func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.String(http.StatusBadRequest, "Can't delete asset without filename given!")
			return
		}
		path := fmt.Sprintf("static/assets/%s", filename)
		err := os.Remove(path)
		if err != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("Asset %s not found on server!", path))
			return
		}
		message := fmt.Sprintf("Removed asset %s from the server!", path)
		log.Println(message)
		c.String(http.StatusNotFound, message)
	})

	r.Run(":8080")
}
