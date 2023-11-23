package main

import (
	"blog/model"
	"blog/web"
	"blog/web/pages"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.HTMLRender = &web.TemplRender{}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", pages.HomePage())
	})

	r.GET("/test", func(c *gin.Context) {
		blogPost := examplePost()
		if blogPost == nil {
			c.HTML(http.StatusNotFound, "", nil)
			return
		}
		c.HTML(http.StatusOK, "", pages.BlogPost(*blogPost))
	})

	r.Static("/static", "./static")

	r.Run(":8080")
}

func examplePost() *model.BlogPost {
	return &model.BlogPost{
		Title:        "Delving into the GOTH Stack: My Experience Building a Blog with GoLang, Templ, and HTMX",
		HeroImageURL: "",
		HTML:         "<h1>Delving into the GOTH Stack: My Experience Building a Blog with GoLang, Templ, and HTMX</h1><h2>Why did i do this?</h2>",
		Author: model.Author{
			Name:     "David Caudill",
			ImageURL: "/static/david.jpg",
			Socials: []model.Social{{
				Name:  "X",
				Image: "",
				URL:   "https://x.com/kirboyyy",
			}},
		},
	}
}
