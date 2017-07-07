NAME = vape

BUILD_NUMBER=$(shell git rev-parse --short HEAD)

build:
	go build -o $(NAME)

clean:
	rm -f vape
	rm -f bin/vape*
	rm -f dist/vape*

release:
	GOOS=darwin GOARCH=amd64 go build -o "bin/$(NAME)_darwin_amd64"
	GOOS=darwin GOARCH=386   go build -o "bin/$(NAME)_darwin_386"
	GOOS=linux  GOARCH=amd64 go build -o "bin/$(NAME)_linux_amd64"
	GOOS=linux  GOARCH=386   go build -o "bin/$(NAME)_linux_386"
	tar cvfz dist/vape-$(BUILD_NUMBER).tar.gz bin/vape*
	gpg --sign dist/vape-$(BUILD_NUMBER).tar.gz

build-docker:
	docker build -t symm/vape:latest .

config:
	cp Vapefile.example Vapefile

run-docker:
	docker run \
		--rm \
		-t \
		-v $(PWD)/Vapefile.example:/Vapefile \
		symm/vape:latest \
		https://httpbin.org

test:
	go test ./... -cover -coverprofile=coverage.out -covermode=atomic

coverage: test
	go tool cover -html=coverage.out

push:
	docker login -u="$(DOCKER_USERNAME)" -p="$(DOCKER_PASSWORD)";
	docker push symm/vape:latest
