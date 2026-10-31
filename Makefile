.PHONY: test lint build docker tf-validate helm-lint run-local clean

test:
	@echo "==> Running Go unit tests with coverage..."
	cd app && go test -v -race -coverprofile=../coverage.out ./src/...

lint:
	@echo "==> Running static analysis & vet..."
	cd app && go vet ./src/...

build:
	@echo "==> Compiling static binary..."
	cd app && CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/core-service ./src

docker:
	@echo "==> Building distroless container image..."
	docker build -t core-service:latest -f app/Dockerfile app/

run-local:
	@echo "==> Starting core service locally on :8080..."
	cd app && go run ./src/main.go

helm-lint:
	@echo "==> Linting Helm charts..."
	helm lint gitops/helm/core-service

tf-validate:
	@echo "==> Validating AWS Terraform..."
	cd terraform/environments/aws && terraform init -backend=false && terraform validate
	@echo "==> Validating GCP Terraform..."
	cd terraform/environments/gcp && terraform init -backend=false && terraform validate
	@echo "==> Validating Azure Terraform..."
	cd terraform/environments/azure && terraform init -backend=false && terraform validate

clean:
	rm -rf bin/ coverage.out
