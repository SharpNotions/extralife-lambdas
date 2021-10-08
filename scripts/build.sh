#!/bin/sh

export GO111MODULE=on
go mod download
env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/update-data lambdas/db.go lambdas/shared.go lambdas/update-data.go
env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-participant lambdas/db.go lambdas/shared.go lambdas/get-participant.go
env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-team lambdas/db.go lambdas/shared.go lambdas/get-team.go