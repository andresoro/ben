package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/russross/blackfriday"

	"github.com/gin-gonic/gin"
)

var (
	posts     = flag.String("p", "posts", "Folder where .md posts are held")
	templates = flag.String("t", "templates", "folder where templates are held")
)

// map md file names to raw html bytes
var html map[string][]byte

// Post represents a blog post
type Post struct {
	Title   string
	Content string
}

func main() {
	flag.Parse()

	r := gin.Default()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("./templates/*.tmpl.html")
	r.Delims("{{", "}}")

	r.GET("/", indexHandler())
	r.GET("/:post", postHandler())

	r.Run()
}

// Handler index
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

// Handle individual posts
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

// convert markdown into html map
func processMarkdown() error {
	files, err := ioutil.ReadDir("./posts")
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()
		md, err := ioutil.ReadFile("./posts" + name)
		if err != nil {
			return err
		}
		html[name] = blackfriday.MarkdownCommon([]byte(md))
	}
	return nil
}
