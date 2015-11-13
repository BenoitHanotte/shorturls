all: start

start: build run

build:
	docker-compose build

run:
	 docker-compose up -d

stop:
	docker-compose stop
