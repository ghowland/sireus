#!/bin/bash

if [ ! -d "./build" ]; then
	echo "Creating ./build/ directory"
	mkdir ./build
fi


go build -o build/sireus sireus.go

if [ $? -eq 0 ] ; then
	echo "Running"
	./build/sireus
else
	echo "Failed to build, will not run."
fi

