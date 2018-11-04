package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	port   = flag.String("port", ":8080", "port to host server on")
	assets = flag.String("assets", "", "directory of files to serve over http")
)

func fileHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, *assets)
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("./templates/*.tmpl.html")
	r.Delims("{{", "}}")

	r.GET("/", indexHandler())

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
