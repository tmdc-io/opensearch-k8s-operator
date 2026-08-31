ARG OPENSEARCH_VERSION=2.19.3
FROM docker.io/opensearchproject/opensearch:${OPENSEARCH_VERSION}

USER 0
RUN mkdir -p /licenses
COPY LICENSE /licenses/LICENSE
USER 1000
