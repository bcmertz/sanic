package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	etag := `"` + etag("test.mp4") + `"` // because this is what aws does
	mp4Handler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "HEAD" {
			w.Header().Set("ETag", etag)
			w.Header().Set("Content-Length", "3496607")
			w.WriteHeader(http.StatusOK)
			// fmt.Println("ETAG HEAD", etag)
		} else if req.Method == "GET" {
			// byteRange := req.Header.Get("Range")
			// fmt.Println("range:", byteRange)
			http.ServeFile(w, req, "./test.mp4")
		}
	}

	http.HandleFunc("/mp4", mp4Handler)

	//headHandler := func(w http.ResponseWriter, req *http.Request) {
	//
	//}

	// handler == nil (DefaultServeMux), Handle and HandleFunc add handlers to this
	log.Print("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func etag(file_name string) string {
	var sum_md5 [16]byte

	downloaded_bytes, _ := ioutil.ReadFile(file_name) // get downloaded_bytes []byte from newly created file

	sum_md5 = md5.Sum(downloaded_bytes)
	return hex.EncodeToString(sum_md5[:]) // convert [16]byte to []byte, then convert to string

}
