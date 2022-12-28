
.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handleNowPlaying ./cmd/spotify/handleNowPlaying/ 

clean:
	# rm -rf ./bin
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose