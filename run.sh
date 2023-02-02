#!/bin/bash

if [ ! -d "./build" ]; then
	echo "Creating ./build/ directory"
	mkdir ./build
fi


go build -o build/sireus sireus.go
./build/sireus

