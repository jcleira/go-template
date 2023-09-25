test:
	go test -race ./...

lint:
	golangci-lint run

migration-create:
	migrate create -ext sql -dir data/db/migrations -seq $(name)

build:
	docker-compose -f ~/.code.envs/docker-compose.yml \
		-f ./docker-compose.yml \
		--project-directory . \
		--project-name local \
		build

run: build
	docker-compose -f ~/.code.envs/docker-compose.yml \
		-f ./docker-compose.yml \
		--project-directory . \
		--project-name local \
		up -d

down:
	docker-compose -f ~/.code.envs/docker-compose.yml \
		-f ./docker-compose.yml \
		--project-directory . \
		--project-name local \
	  down --remove-orphans
