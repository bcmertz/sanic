# go-download
## How to use:
### Usage:

`go run torrent.go [--n int] [--u string [--o string]] [--o string] [--v bool]`

### Options:

	--n	Number of goroutines to download from
	--o 	Output filename of specified resource
	--u	Url of resource to download [requires -o]
	--v	Verify etag == md5 hash of output file

## Examples:
`go run torrent.go --v=true --o=bigfile.mp4 --n=10`

`go run torrent.go`

`go run torrent.go --v=false --u=http://download.blender.org/peach/bigbuckbunny_movies/BigBuckBunny_320x180.mp4 --o=bigbuchbunny.mp4`

# Testing Perfomance - # of Goroutines

`go run server.go && ./test.sh`

## Options:

	--r Rate of throttling for the server
