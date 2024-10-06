
build:
	go build -o runner main.go

run:
	./runner

dc-build:
	templ generate
	docker-compose build

dc-run:
	docker-compose up

dc-rund:
	docker-compose up -d

dc-stop:
	docker-compose kill
	docker-compose rm -f