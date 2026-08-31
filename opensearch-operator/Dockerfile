# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.25.13 AS builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/   pkg/
COPY webhook/ webhook/
COPY opensearch-gateway/   opensearch-gateway/

# Build with Go's FIPS 140-3 module compiled in (Go 1.24+)
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOFIPS140=v1.0.0 go build -a -o manager main.go

FROM --platform=$BUILDPLATFORM docker.io/tmdcio/ubi:ubi10-minimal-fips-curated.v1

LABEL name="opensearch-operator" \
      vendor="OpenSearch Project" \
      summary="OpenSearch Kubernetes Operator" \
      description="Kubernetes operator for managing OpenSearch clusters and OpenSearch Dashboards" \
      io.k8s.display-name="OpenSearch Operator" \
      io.k8s.description="Kubernetes operator for managing OpenSearch clusters and OpenSearch Dashboards" \
      io.openshift.tags="opensearch,operator,database,search" \
      com.redhat.component="opensearch-operator" \
      maintainer="OpenSearch Project <support@opensearch.org>"
COPY LICENSE /licenses/LICENSE
WORKDIR /
COPY --from=builder /workspace/manager .
ENV GODEBUG=fips140=on
USER 65532:65532

ENTRYPOINT ["/manager"]
