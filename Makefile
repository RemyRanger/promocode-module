# Docker Compose Up
docker_compose_up:
	docker compose up -d

# Docker Compose Up
docker_compose_down:
	docker compose down

# Install all Go dependencies
install:
	go mod download
	go install tool

# Generate OAS bundled document in /dist
generate_api_interfaces:
	oapi-codegen -version
	go generate -x $$(go list ./... | grep -v '/ports\|/adapters|/core' | tr '\n' ' ')

# Build binary
build: generate_api_interfaces
	go build -o cmd/promocode-module cmd/main.go

run: generate_api_interfaces
	go run cmd/main.go

update_deps:
	go get -u ./...
	go mod tidy

test_go: generate_api_interfaces
	ginkgo -v -r -cover -coverpkg=APIs/internal/... -coverprofile=coverage.cov ./...

coverage_html:
	go tool cover -html=coverage.cov