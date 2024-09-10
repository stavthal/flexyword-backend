
check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	GO111MODULE=on go mod vendor  && GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml


# Docker related
DOCKER_IMAGE=flexword-api
CONTAINER_NAME=flexyword-container
PORT=8080

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run --name $(CONTAINER_NAME) -p $(PORT):8080  -d --env-file .env $(DOCKER_IMAGE)

docker-stop:
	docker stop $(CONTAINER_NAME)

docker-start:
	docker start $(CONTAINER_NAME)

docker-rm:
	docker rm $(CONTAINER_NAME)

docker-image-prune:
	docker image prune 

docker-image-prune-sudo:
	docker image prune -a

# Goose Related
migrate-up:
	@source .env && goose -dir ./db/migrations postgres $$SUPABASE_DB_URL up

migrate-down:
	@source .env && goose -dir ./db/migrations postgres $$SUPABASE_DB_URL down

migrate-status:
	@source .env && goose -dir ./db/migrations postgres $$SUPABASE_DB_URL status

migrate-create:
	@source .env && goose -dir ./db/migrations create $(name) sql
