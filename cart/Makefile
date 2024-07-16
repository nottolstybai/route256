COVERAGE_DIR := artifacts
REPO_PKG := ./internal/repository/...
SERVICE_PKG := ./internal/service/...
CMD := ./cmd/cart
SERVER := ./cmd/cart/main.go

.PHONY: lint coverage build run clean test

lint:
	@echo "Checking cyclomatic complexity"
	@gocyclo -over 10 -ignore "_test/" .
	@echo "Checking cognitive complexity"
	@gocognit -over 10 .

test:
	go test -parallel 50 -race $(REPO_PKG) && \
	go test -race $(SERVICE_PKG)

clean:
	rm -rf $(COVERAGE_DIR) bin

coverage:
	mkdir $(COVERAGE_DIR) && \
	go test -coverprofile $(COVERAGE_DIR)/repository.cover.out $(REPO_PKG) && \
	go tool cover -html $(COVERAGE_DIR)/repository.cover.out -o $(COVERAGE_DIR)/repository.cover.html && \
	go test -coverprofile $(COVERAGE_DIR)/service.cover.out $(SERVICE_PKG) && \
    go tool cover -html $(COVERAGE_DIR)/service.cover.out -o $(COVERAGE_DIR)/service.cover.html

build:
	go build -o bin/cart $(CMD)

run:
	@go run $(SERVER)



LOCAL_BIN:=$(CURDIR)/bin

PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto


# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp


# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem


# Вендоринг внешних proto файлов
.vendor-proto: vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options


.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0


LOMS_PROTO_PATH:="/api/loms/v1"

PHONY: .protoc-generate
.protoc-generate: .bin-deps .vendor-proto
	protoc \
	-I $(PWD)/../loms/api/loms/v1 \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go.exe \
	--go_out pkg/${LOMS_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc.exe \
	--go-grpc_out pkg/${LOMS_PROTO_PATH} \
	--go-grpc_opt paths=source_relative \
	$(PWD)/../loms/api/loms/v1/loms.proto
	go mod tidy