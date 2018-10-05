#!/bin/bash

echo "Begining Test"
$(touch results.csv)
echo "goroutines,time" > results.csv
for i in 5 10 25 125
do
    start=$(date +%s.%N)
    $(go run torrent.go "http://download.blender.org/peach/bigbuckbunny_movies/BigBuckBunny_320x180.mp4" $i)
    elapsed=$(echo "$(date +%s.%N) - $start"|bc)
    echo "$i,$elapsed" >> results.csv
done
