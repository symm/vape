NAME = vape

build:
	docker build -t symm/vape:latest .

run:
	docker run \
		--rm \
		-t \
		-v $(PWD)/Vapefile.example:/Vapefile \
		symm/vape:latest \
		https://httpbin.org

test:
	docker run \
		--rm \
		-t \
		-v $(PWD)/:/usr/src/myapp \
		-w /usr/src/myapp \
		golang:alpine \
		go test ./... -cover
