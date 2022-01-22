#!/bin/bash

PROJECT_NAME=jrctl
goos=("linux" "windows" "darwin")
arch=("386" "amd64" "arm64")

if [ -d ./artifacts ];then
  rm -fr ./artifacts
fi
if [ -d ./bin ];then
  rm -fr ./bin
fi

mkdir -p ./artifacts ./bin

for a in ${goos[@]}; do
	for i in ${arch[@]}; do
    if [[ "$a" == "darwin" && "$i" == "386" ]] || [[ "$a" == "windows" && "$i" == "arm64" ]];then
	  		echo "Not allowed: $a $i"
		else
	  		echo "Building: $a $i"
 			CGO_ENABLED=0 GOOS=$a GOARCH=$i go build -o ./bin/$a/$i/${PROJECT_NAME} &&\
			tar -czvf ./artifacts/${PROJECT_NAME}-$a-$i.tar.gz ./bin/$a/$i/${PROJECT_NAME} &
		fi
	done
done; wait
