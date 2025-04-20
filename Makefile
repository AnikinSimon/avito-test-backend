.PHONY: lint

lint:
	golangci-lint run

api-gen:
	oapi-codegen --config=./configs/oapi.yaml api/swagger.yaml

PROTO_DIR = api
PROTO_OUT = internal/grpc/pvz/v1

.PHONY: gen-proto
gen-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_OUT) \
		--go-grpc_out=$(PROTO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/pvz.proto

sqlc:
	sqlc generate -f ./database/sqlc.yaml

gogen:
	go generate ./...

run:
	docker compose up -d --build

stop:
	docker compose stop

unit:
	go test -v ./...

.PHONY: test

test:
	go test -v -tags=integrations -count=1 ./...


