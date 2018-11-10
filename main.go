package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"

	"github.com/gin-gonic/gin"
)

var (
	data      = flag.String("p", "posts", "Folder where .md posts are held")
	templates = flag.String("t", "templates", "folder where templates are held")
)

type posts map[string][]byte

// Post represents a blog post
type Post struct {
	Title   string
	Content string
}

func main() {
	flag.Parse()

	// process all posts into memory
	posts, err := processMarkdown()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("./templates/*.tmpl.html")
	r.Delims("{{", "}}")

	r.GET("/", indexHandler())
	r.GET("/:post", postHandler(posts))

	r.Run()
}

// Handler index
func indexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		posts := make(map[string]string)
		files, err := ioutil.ReadDir(*data)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fmt.Println(file.Name())
			path := file.Name()
			name := strings.TrimSuffix(path, filepath.Ext(path))

			posts[path] = name
		}
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"posts": posts,
		})
	}
}

// Handle individual posts
func postHandler(p posts) gin.HandlerFunc {
	return func(c *gin.Context) {
		postName := c.Param("post")

		md, ok := p[postName]
		if !ok {
			c.HTML(http.StatusOK, "post.tmpl.html", gin.H{
				"Title":   "Not found",
				"Content": "This post does not exist",
			})
			return
		}

		html := template.HTML(md)

		c.HTML(http.StatusOK, "post.tmpl.html", gin.H{
			"Title":   postName,
			"Content": html,
		})
	}
}

// convert markdown into html map
func processMarkdown() (map[string][]byte, error) {
	// map md file names to raw html bytes
	html := make(map[string][]byte)

	files, err := ioutil.ReadDir("./posts")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		md, err := ioutil.ReadFile("./posts/" + file.Name())
		if err != nil {
			return nil, err
		}
		html[name] = blackfriday.MarkdownCommon([]byte(md))
	}
	return html, nil
}
