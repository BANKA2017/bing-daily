#!/usr/bin/bash

rm -rf build

cd crawler
go build crawler.go
cd ../server
go build server.go
cd ../

mkdir build/
mv crawler/crawler build/
mv server/server build/
cp .env.example build/.env
echo "[]" > build/bing.json
