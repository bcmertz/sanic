package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	chunks := flag.Int("n", 25, "number of goroutines to downlaod from")
	url := flag.String("u", "https://s3-us-west-2.amazonaws.com/getlantern-test/downloaded_video.mp4", "url to download from")
	file_name := flag.String("o", "download.mp4", "name of downloaded file")
	verify := flag.Bool("v", true, "verify md5 hash of download against etag")

	flag.Parse() // parse flags from os.Args[1:]

	file, _ := os.Create(*file_name)

	size, etag := check_size(*url)
	etag = etag[1 : len(etag)-1]   // remove quotes from etag
	for i := 0; i < *chunks; i++ { // create all our goroutines to download over range [start:end]
		wg.Add(1)
		start := i * size / *chunks
		end := (i + 1) * size / *chunks
		client := &http.Client{} // client has internal state (cached TCP connections) and so shoul dbe reused as needed [from go http/client documentation]
		go download(*url, start, end, file, client)
	}

	wg.Wait()    // wait until all our goroutines close
	file.Close() // able to be closed now that we've written to it

	if *verify {
		verify_download(etag, *file_name) // compare etag hash and md5 hash of downloaded file
	} else {
		print("download finished \n")
	}
}

func download(url string, start, end int, file *os.File, client *http.Client) {
	req, err := http.NewRequest("GET", url, nil) // create a new request so we can specify the download range header for the request
	if err != nil {
		panic(err)
	}

	download_range := "bytes=" + strconv.Itoa(start) + "-" + strconv.Itoa(end) // format: "bytes=XXX-YYY"
	req.Header.Set("Range", download_range)

	resp, err := client.Do(req) // execute req (http request) over defined range
	if err != nil {
		print("download over byte range: ", start, " - ", end, " failed - resulting file will be incomplete \n")
	}

	defer resp.Body.Close() // must close response
	defer wg.Done()         // once goroutine has finished executing we can inform the waitgroup it has finished - prevent early waitgroup finish

	body, err := ioutil.ReadAll(resp.Body) // get body []byte which we can write into file
	if err != nil {
		print("failed to extract part of file data from server response - output file will be incomplete \n")
	}
	_, err = file.WriteAt(body, int64(start))
	if err != nil {
		print("failed to write part of downloaded data - output file will be incomplete \n")
	}
}

func check_size(url string) (length int, etag string) {
	resp, err := http.Head(url)
	if err != nil {
		panic(err) // we can't download the file properly if we don't know it's content length
	}
	etag = resp.Header.Get("ETag")
	length = int(resp.ContentLength)
	return
}

func verify_download(etag string, file_name string) {
	var hash_sum string
	var sum_md5 [16]byte

	downloaded_bytes, err := ioutil.ReadFile(file_name) // get downloaded_bytes []byte from newly created file
	if err != nil {
		print(file_name, " downloaded - could not open file to verify etag hash")
	}

	sum_md5 = md5.Sum(downloaded_bytes)
	hash_sum = hex.EncodeToString(sum_md5[:]) // convert [16]byte to []byte, then convert to string

	if hash_sum == etag {
		print(file_name, " downloaded - hash verified", "\n")
	} else {
		print(file_name, " downloaded - not hash verified", "\n")
	}
}

func depricated_verify_download(etag string, file_name string) {
	// this function was meant to be able to handle aws s3 multiple chunk uploads however they don't have any comprehensible hashing method
	// online posts suggest that each upload chunk is hashed, all are concatenated, and then the string of concatenated hashes are hashed
	// however this approach only works if you know the chunk size of your upload which must be guessed as there is no way to deduce that from etag
	// you can know the number of chunks, but for a file of 61 mb for example, logical chunking of [15, 15, 15, 16] or [20, 20, 20, 1] both fail
	// ultimately i decided to depricate this because there is no transparency as to how s3 hashes etags and the approach appears to vary for files of
	// different size and I don't want to hardcode behavior based off of deduced behavior in 5 upvote stackoverflow posts... sorry no multichunk s3 upload support

	format_etag := strings.Split(etag, "-") //commented out to remove strings import
	var hash_sum string
	var sum_md5 [16]byte
	downloaded_bytes, err := ioutil.ReadFile(file_name)

	if err != nil {
		panic(err)
	}

	if len(format_etag) == 1 {
		sum_md5 = md5.Sum(downloaded_bytes)
		hash_sum = hex.EncodeToString(sum_md5[:])
	} else {
		num_chunks_upload, _ := strconv.Atoi(format_etag[len(format_etag)-1])
		size := len(downloaded_bytes)
		chunk_size := 15*10 ^ 6 // test s3 chunk size of 15mb -- failed [as did many other tests]
		remainder := size - ((num_chunks_upload - 1) * chunk_size)
		partial_hash_sum := ""

		for i := 0; i < num_chunks_upload; i++ {
			start := i * chunk_size
			end := (i + 1) * chunk_size
			if num_chunks_upload-1 == i {
				sum_md5 = md5.Sum(downloaded_bytes[len(downloaded_bytes)-remainder:])
				partial_hash_sum += string(sum_md5[:])
			} else {
				sum_md5 = md5.Sum(downloaded_bytes[start:end])
				partial_hash_sum += string(sum_md5[:])
			}
		}
		sum_md5 = md5.Sum([]byte(partial_hash_sum))
		hash_sum = hex.EncodeToString(sum_md5[:]) + "-" + strconv.Itoa(num_chunks_upload)
	}

	print("downloaded file hash sum: ", hash_sum, " etag: ", etag, "\n")
	if hash_sum == etag {
		print("successful download of ", file_name, "\n")
	} else {
		print("could not verify successful download of ", file_name, "\n")
	}
}
