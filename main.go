package main

import (
	"flag"
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

func main() {
	flag.Parse()

	var h http.Handler

	file, err := os.Open(*assets)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		h = http.StripPrefix(*assets, http.FileServer(http.Dir(fi.Name())))
		http.Handle("/", h)
	case mode.IsRegular():
		http.HandleFunc("/", fileHandle)

	}

	log.Fatal(http.ListenAndServe(*port, nil))

}
