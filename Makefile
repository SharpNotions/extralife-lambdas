.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	go mod download
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/update-data lambdas/db.go lambdas/shared.go lambdas/update-data.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get lambdas/db.go lambdas/shared.go lambdas/get.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
