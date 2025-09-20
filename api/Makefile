.PHONY: init
init: gitsm proto

.PHONY: gitsm
gitsm:
	git submodule update --init --recursive && \
	cd pb && \
	git checkout main && \
	cd ./../..


.PHONY: proto
proto:
	mkdir -p pb/gen/go

	protoc -I pb/proto -I pb/lib \
	--go_out pb/gen/go \
	--go_opt paths=source_relative \
	--go-grpc_out pb/gen/go \
	--go-grpc_opt paths=source_relative \
	  pb/proto/*/*.proto

.PHONY: swag
swag:
	@echo "swagger generating..."
	@swag init -g ./cmd/main.go -q -o ./docs
	@echo "swagger generated"

