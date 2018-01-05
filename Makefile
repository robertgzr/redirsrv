all: build

build:
	@cargo build

docker:
	docker build -t robertgzr/redirsrv:latest .

run:
	docker run -d \
		--name redirsrv \
		-e ROCKET_SECRET_KEY=$(shell openssl rand -base64 32) \
		-v $(shell pwd)/Rocket.toml:/Rocket.toml \
		-v $(shell pwd)/linkfile.json:/etc/redirsrv/linkfile.json \
		-p 8080:80 \
		robertgzr/redirsrv:latest

clean:
	@cargo clean

.PHONY: clean build linux
