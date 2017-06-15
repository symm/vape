build:
	docker build -t symm/vape:latest .
run:
	docker run \
		--rm \
		-t \
		-v $(PWD)/Vapefile.example:/Vapefile \
		symm/vape:latest \
		https://httpbin.org
