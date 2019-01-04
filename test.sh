#!/bin/bash

echo "Beginning Test"
$(touch results.csv)
echo "goroutines,time" > results.csv
for i in 2 10 25 50 125 
do
    start=$(date +%s.%N)
    $(go run torrent.go --n=$i --v=false --u=http://localhost:8080/mp4)
    elapsed=$(echo "$(date +%s.%N) - $start"|bc) #from online resource
    echo "$i,$elapsed" >> results.csv
done
