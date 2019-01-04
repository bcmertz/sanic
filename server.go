package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	rate := flag.Float64("r", 200, "rate kb/s the server throttles traffic")
	flag.Parse()

	etag := `"` + etag("test.mp4") + `"` // because this is what aws does
	mp4Handler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "HEAD" {
			w.Header().Set("ETag", etag)
			w.Header().Set("Content-Length", "3496607")
			w.WriteHeader(http.StatusOK)
			// fmt.Println("ETAG HEAD", etag)
		} else if req.Method == "GET" {
			rate_limiter(req, rate)
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

func rate_limiter(req *http.Request, rate *float64) {
	RATE_KBS := *rate                                   // kb/s
	reqrange := req.Header.Get("Range")[len("bytes="):] //remove bytes=
	byteRange := strings.Split(reqrange, "-")
	br1, _ := strconv.Atoi(byteRange[0])
	br2, _ := strconv.Atoi(byteRange[1])
	rangeint := br2 - br1
	microsec := (float64(rangeint) * math.Pow10(6)) / (RATE_KBS * 1000.0)
	sleep_for := time.Duration(float64(microsec)) * time.Microsecond
	time.Sleep(sleep_for)
}

func etag(file_name string) string {
	var sum_md5 [16]byte

	downloaded_bytes, _ := ioutil.ReadFile(file_name) // get downloaded_bytes []byte from newly created file

	sum_md5 = md5.Sum(downloaded_bytes)
	return hex.EncodeToString(sum_md5[:]) // convert [16]byte to []byte, then convert to string

}
