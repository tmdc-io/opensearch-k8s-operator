### CRD Installation (server-side apply for large CRDs)

CRD_DIR = charts/opensearch-operator/files

.PHONY: install-crds

install-crds:
	@echo "=== Installing CRDs with server-side apply ==="
	@for f in $(CRD_DIR)/*.yaml; do \
		echo "Applying $$f ..."; \
		kubectl apply --server-side --force-conflicts -f "$$f"; \
	done
	@echo "=== CRDs installed successfully ==="

### Helm push to ECR steps

CH_DIR = charts
CHARTS = opensearch-operator opensearch-cluster
VERSION = ${TAG}

.PHONY: push-all-charts push-oci-chart ecr-login ecr-logout

push-all-charts: ecr-login
	@for chart in $(CHARTS); do \
		echo ""; \
		echo "=== processing chart: $$chart ==="; \
		echo "=== package OCI chart: $$chart ==="; \
		helm3.16.4 package $(CH_DIR)/$$chart/ --version $(VERSION); \
		echo "=== create repository: $$chart ==="; \
		aws ecr describe-repositories --repository-names $$chart --no-cli-pager 2>/dev/null || \
			aws ecr create-repository --repository-name $$chart --region $(AWS_DEFAULT_REGION) --no-cli-pager; \
		echo "=== push OCI chart: $$chart ==="; \
		helm3.16.4 push $$chart-$(VERSION).tgz oci://$(ECR_HOST); \
		echo ""; \
	done
	@$(MAKE) ecr-logout

push-oci-chart:
	@test -n "$(DIR)" || (echo "ERROR: DIR is required. Usage: make push-oci-chart DIR=opensearch-operator"; exit 1)
	@$(MAKE) ecr-login
	@echo
	@echo "=== package OCI chart: $(DIR) ==="
	helm3.16.4 package $(CH_DIR)/$(DIR)/ --version $(VERSION)
	@echo
	@echo "=== create repository: $(DIR) ==="
	aws ecr describe-repositories --repository-names $(DIR) --no-cli-pager 2>/dev/null || \
		aws ecr create-repository --repository-name $(DIR) --region $(AWS_DEFAULT_REGION) --no-cli-pager
	@echo
	@echo "=== push OCI chart: $(DIR) ==="
	helm3.16.4 push $(DIR)-$(VERSION).tgz oci://$(ECR_HOST)
	@$(MAKE) ecr-logout

ecr-login:
	@echo "=== login to OCI registry ==="
	aws ecr get-login-password --region $(AWS_DEFAULT_REGION) | helm3.16.4 registry login $(ECR_HOST) --username AWS --password-stdin --debug

ecr-logout:
	@echo "=== logout of registry ==="
	helm3.16.4 registry logout $(ECR_HOST)

## TMDC
IMAGE_REPO ?= docker.io/tmdcio/opensearch-operator
PLATFORM ?= linux/amd64
GITHUB_TAGS ?= latest

.PHONY: build-tmdc-docker
build-tmdc-docker: ## Build TMDC Docker image for linux/amd64
	@echo "build docker container"
	@cp -f LICENSE opensearch-operator/LICENSE
	docker buildx build \
		--output type=docker \
		--platform $(PLATFORM) \
		-f opensearch-operator/Dockerfile \
		--tag $(IMAGE_REPO):$(GITHUB_TAGS) \
		opensearch-operator

.PHONY: push-tmdc-docker
push-tmdc-docker: ## Push TMDC Docker image to Docker Hub
	docker push $(IMAGE_REPO):$(GITHUB_TAGS)

CLUSTER_IMAGE_REPO ?= docker.io/tmdcio/opensearch-cluster

.PHONY: build-tmdc-cluster-docker
build-tmdc-cluster-docker: ## Build TMDC OpenSearch cluster Docker image for linux/amd64
	@echo "build cluster docker container"
	@cp -f LICENSE opensearch-cluster/LICENSE
	docker buildx build \
		--output type=docker \
		--platform $(PLATFORM) \
		-f opensearch-cluster/Dockerfile \
		--tag $(CLUSTER_IMAGE_REPO):$(GITHUB_TAGS) \
		opensearch-cluster

.PHONY: push-tmdc-cluster-docker
push-tmdc-cluster-docker: ## Push TMDC OpenSearch cluster Docker image to Docker Hub
	docker push $(CLUSTER_IMAGE_REPO):$(GITHUB_TAGS)