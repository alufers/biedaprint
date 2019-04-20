#!/bin/bash
go run ./gocode_docs_extractor
cd frontend
node_modules/.bin/vue-cli-service build --dest ../static 
cd ..
GOARCH=arm GOARM=7 ~/go/bin/packr build .
scp ./biedaprint root@192.168.1.213:/home/steam/biedaprint
ssh -t  root@192.168.1.213 './biedaprint'