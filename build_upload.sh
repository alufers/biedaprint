#!/bin/bash
cd frontend
node_modules/.bin/vue-cli-service build --dest ../static 
cd ..
~/go/bin/packr build .
scp ./biedaprint alufers@192.168.1.15:/home/alufers/biedaprint
ssh -t  alufers@192.168.1.15 './biedaprint'