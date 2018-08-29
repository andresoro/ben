package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	port   = flag.String("port", ":8080", "port to host server on")
	assets = flag.String("assets", "", "directory of files to serve over http")
)

func fileHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, *assets)
}

func folderHandle(w http.ResponseWriter, r *http.Request) {

}

func main() {
	flag.Parse()

	var name string

	file, err := os.Open(*assets)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	name = fi.Name()

	switch mode := fi.Mode(); {
	case mode.IsDir():
		//handle
		fs := http.FileServer(http.Dir(name))
		http.Handle("/", fs)

		//log
		files, err := ioutil.ReadDir(name)
		if err != nil {
			log.Printf("Error looping files: %s", err)
		}
		log.Printf("Serving following files in %s", name)
		for _, f := range files {
			fmt.Println(f.Name())
		}

	case mode.IsRegular():
		log.Printf("Handling %s", name)
		http.HandleFunc("/", fileHandle)

	}
	log.Printf("Hosting server on port %s", *port)
	log.Fatal(http.ListenAndServe(*port, nil))

}
