.PHONY: build clean deploy gomodgen

build: gomodgen
	scripts/build.sh

clean:
	rm -rf ./bin ./vendor go.sum

deploy: build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
