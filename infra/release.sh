#!/bin/bash

goos=(
"linux"
"windows"
"darwin"
)

arch=(
"386"
"amd64"
"arm64"
)

for a in ${goos[@]}; do
	for i in ${arch[@]}; do
		if [[ $a == "darwin" && $i == "386" ]]; then
	  		echo "Not allowed: $a $i"
		else
	  		echo "Building: $a $i"
 			CGO_ENABLED=0 GOOS=$a GOARCH=$i go build -o ./bin/$a/$i/piktoctl ./ &
		fi
	done
done; wait

