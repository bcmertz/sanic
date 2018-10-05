package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func main() {
	var file_path string
	var chunks int
	var err error

	custom_arguments := os.Args[1:]
	if len(custom_arguments) < 2 {
		file_path = "https://s3-us-west-1.amazonaws.com/mybucket-bennettmertz/myvideo.mp4"
		chunks = 25
	} else {
		file_path = custom_arguments[0]
		chunks, err = strconv.Atoi(custom_arguments[1])
		if err != nil {
			panic(err)
		}
	}

	file, _ := os.Create("downloaded_video.mp4")
	defer file.Close()

	size := check_size(file_path)
	for i := 0; i < chunks; i++ {
		wg.Add(1)
		start := i * size / chunks
		end := (i + 1) * size / chunks
		go download(file_path, start, end, file)
	}

	wg.Wait()
}

func download(file_path string, start, end int, file *os.File) {
	req, err := http.NewRequest("GET", file_path, nil)
	if err != nil {
		panic(err)
	}

	download_range := "bytes=" + strconv.Itoa(start) + "-" + strconv.Itoa(end)

	req.Header.Set("Range", download_range)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	defer wg.Done()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteAt(body, int64(start))
	if err != nil {
		panic(err)
	}
}

func check_size(url string) (length int) {
	resp, err := http.Head(url)
	if err != nil {
		panic(err)
	}
	length = int(resp.ContentLength)
	return
}
