#!/bin/bash

echo "Begining Test"

for i in 1 5 25 125
do
    $(go run torrent.go "https://s3-us-west-1.amazonaws.com/mybucket-bennettmertz/myvideo.mp4" $i)
done
