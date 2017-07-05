NAME = vape

build:
	go build -o $(NAME)

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
	go test ./... -cover -coverprofile=coverage.out

coverage: test
	go tool cover -html=coverage.out

push:
	docker login -u="$(DOCKER_USERNAME)" -p="$(DOCKER_PASSWORD)";
	docker push symm/vape:latest
