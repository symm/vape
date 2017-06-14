build:
	docker build -t symm/vape:latest .
run:
	docker run \
		--rm \
		-t \
		-v $(PWD)/smoke.example:/.smoke \
		symm/vape:latest
