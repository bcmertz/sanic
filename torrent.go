package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	start := time.Now()
	var file_path string
	var chunks int
	var err error

	custom_arguments := os.Args[1:]
	if len(custom_arguments) < 2 {
		// if arguments improperly given then use default arguments
		file_path = "http://download.blender.org/peach/bigbuckbunny_movies/BigBuckBunny_320x180.mp4"
		chunks = 10
	} else {
		file_path = custom_arguments[0]
		chunks, err = strconv.Atoi(custom_arguments[1])
		if err != nil {
			panic(err)
		}
	}

	size := check_size(file_path)
	chans := make([]chan []byte, chunks)
	for i := range chans {
		chans[i] = make(chan []byte, 1)
		wg.Add(1)
		start := i * size / chunks
		end := (i + 1) * size / chunks
		go download(file_path, start, end, chans[i])
	}

	wg.Wait()

	file, _ := os.Create("downloaded_video.mp4")
	defer file.Close()

	for i := range chans {
		start := i * size / chunks
		file.WriteAt(<-chans[i], int64(start))
	}
	end := time.Now()
	elapsed := end.Sub(start)
	print("time elapsed: %v \n", elapsed)
}

func download(file_path string, start, end int, channel chan []byte) {
	req, err := http.NewRequest("GET", file_path, nil)

	download_range := "bytes=" + strconv.Itoa(start) + "-" + strconv.Itoa(end)

	req.Header.Set("Range", download_range)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("download failed")
	}
	defer resp.Body.Close()
	defer wg.Done()

	body, err := ioutil.ReadAll(resp.Body)
	channel <- body
	close(channel)
}

func check_size(url string) (length int) {
	resp, err := http.Head(url)
	if err != nil {
		panic(err)
	}
	length = int(resp.ContentLength)
	return
}
