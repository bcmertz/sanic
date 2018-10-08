#!/bin/bash

echo "Beginning Test"
$(touch results.csv)
echo "goroutines,time" > results.csv
for i in 1 2 5 25 125
do
    start=$(date +%s.%N)
    $(go run torrent.go -n=$i)
    elapsed=$(echo "$(date +%s.%N) - $start"|bc) #from online resource
    echo "$i,$elapsed" >> results.csv
done
