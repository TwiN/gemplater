#!/bin/bash

APPLICATION_NAME="gemplater"

env GOOS=darwin GOARCH=amd64 go build -mod vendor -o bin/${APPLICATION_NAME}-darwin64
env GOOS=linux GOARCH=amd64 go build -mod vendor -o bin/${APPLICATION_NAME}-linux64
env GOOS=windows GOARCH=amd64 go build -mod vendor -o bin/${APPLICATION_NAME}-windows64.exe

chmod +x bin/${APPLICATION_NAME}-darwin64
chmod +x bin/${APPLICATION_NAME}-linux64
chmod +x bin/${APPLICATION_NAME}-windows64.exe
