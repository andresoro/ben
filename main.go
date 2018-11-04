package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/russross/blackfriday"

	"github.com/gin-gonic/gin"
)

// Post represents a blog post
type Post struct {
	Title   string
	Content string
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("./templates/*.tmpl.html")
	r.Delims("{{", "}}")

	r.GET("/", indexHandler())
	r.GET("/:post", postHandler())

	r.Run()
}

func indexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var posts []string
		files, err := ioutil.ReadDir("./posts")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fmt.Println(file.Name())
			posts = append(posts, file.Name())
		}
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"posts": posts,
		})
	}
}

func postHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		postName := c.Param("post")

		md, err := ioutil.ReadFile("./posts/" + postName)
		if err != nil {
			//TODO: handler error with error page
			log.Fatal(err)
		}
		s := blackfriday.MarkdownCommon([]byte(md))
		html := template.HTML(s)

		c.HTML(http.StatusOK, "post.tmpl.html", gin.H{
			"Title":   postName,
			"Content": html,
		})
	}
}
