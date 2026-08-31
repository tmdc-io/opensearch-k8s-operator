![build](https://github.com/opensearch-project/opensearch-k8s-operator/actions/workflows/docker-build.yaml/badge.svg) ![test](https://github.com/opensearch-project/opensearch-k8s-operator/actions/workflows/testing.yaml/badge.svg) ![release](https://img.shields.io/github/v/release/opensearch-project/opensearch-k8s-operator) [![Golang Lint](https://github.com/opensearch-project/opensearch-k8s-operator/actions/workflows/linting.yaml/badge.svg)](https://github.com/opensearch-project/opensearch-k8s-operator/actions/workflows/linting.yaml) [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/opensearch-operator)](https://artifacthub.io/packages/search?repo=opensearch-operator)

# OpenSearch Kubernetes Operator

The Kubernetes OpenSearch Operator is used for automating the deployment, provisioning, management, and orchestration of OpenSearch clusters and OpenSearch dashboards.

> **API Group Migration Notice:** The operator is migrating from `opensearch.opster.io` to `opensearch.org` API group. Both are currently supported, but `opensearch.opster.io` is deprecated. Please see the Migration Guide for details.

## Getting started

The Operator can be easily installed using Helm on any CNCF-certified Kubernetes cluster:

1. Add the helm repo: `helm repo add opensearch-operator https://opensearch-project.github.io/opensearch-k8s-operator/`
2. Install the Operator: `helm install opensearch-operator opensearch-operator/opensearch-operator`

After installation, you can deploy your first OpenSearch cluster by creating an `OpenSearchCluster` custom resource. Please refer to the [User Guide](./docs/userguide/main.md) for detailed instructions and configuration options.

## Video Tutorial

Learn how to install and use the operator with the [OpenSearch Kubernetes Operator Tutorial Series on YouTube](https://pulse.support/kb/running-opensearch-on-kubernetes-video-tutorial-series).

[![Watch the video](https://github.com/user-attachments/assets/3e8881b4-4b93-4322-86e2-f46baa01cad0)](https://pulse.support/kb/running-opensearch-on-kubernetes-video-tutorial-series)


## Features

- Deploy and manage OpenSearch clusters with multiple node pools
- Deploy and configure OpenSearch Dashboards
- Configure all node roles (cluster_manager, data, ingest, coordinating, etc.)
- Scale cluster resources manually, per node pool
- Rolling version upgrades with quorum-safe restarts
- Online volume expansion for disk scaling
- Certificate management with TLS hot reloading
- Multi-namespace support for managing clusters across organizational boundaries
- Plugin installation during bootstrap phase

## Compatibility

The opensearch k8s operator aims to be compatible to all supported opensearch versions. Please check the table below for details:


| Operator Version                                            | Min Supported Opensearch Version | Max Supported Opensearch Version | Comment                                    |
| ----------------------------------------------------------- |----------------------------------| -------------------------------- | ------------------------------------------ |
| 3.0.0                                                       | 2.19.2                           | latest 3.x                       | Supports the latest OpenSearch 3.x version. |
| 2.8.0                                                       | 2.19.2                           | latest 3.x                       | Supports the latest OpenSearch 3.x version. |
| 2.7.0<br>2.6.1<br>2.6.0<br>2.5.1<br>2.5.0                    | 1.3.x                            | 2.19.2                       | Supports the OpenSearch 2.19.2 version. |


This table only lists versions that have been explicitly tested with the operator, the operator will not prevent you from using other versions. Newer minor versions (2.x) not listed here generally also work but you should proceed with caution and test it out in a non-production environment first.

## Development

If you want to develop the operator, please see the separate [developer docs](./docs/developing.md).

## TMDC Docker images (GitHub Actions)

Operator and cluster images are built and pushed to Docker Hub when a matching git tag is pushed. The image tag is the git tag with the prefix removed.

| Git tag pattern | Workflow | Image |
|-----------------|----------|--------|
| `docker-operator-*` | `.github/workflows/TMDC-DOCKER-image-push.yaml` | `docker.io/tmdcio/opensearch-operator:<version>` |
| `docker-cluster-*` | `.github/workflows/TMDC-DOCKER-cluster-image-push.yaml` | `docker.io/tmdcio/opensearch-cluster:<version>` |

Helm chart tags are separate. See [HELM-CHART-PUSH.md](./HELM-CHART-PUSH.md) for `operator-*` and `cluster-*`.

### Examples

| Git tag | Result |
|---------|--------|
| `docker-operator-3.0.1-d4` | `docker.io/tmdcio/opensearch-operator:3.0.1-d4` |
| `docker-cluster-3.3.1-d1` | `docker.io/tmdcio/opensearch-cluster:3.3.1-d1` |
| `operator-3.0.3-d1` | Helm chart only (operator) |
| `cluster-3.3.1-d1` | Helm chart only (cluster) |

```bash
# Operator image
git tag docker-operator-3.0.1-d4
git push origin docker-operator-3.0.1-d4

# Cluster image
git tag docker-cluster-3.3.1-d1
git push origin docker-cluster-3.3.1-d1
```

Local builds:

```bash
make build-tmdc-docker GITHUB_TAGS=3.0.1-d4
make push-tmdc-docker GITHUB_TAGS=3.0.1-d4

make build-tmdc-cluster-docker GITHUB_TAGS=3.3.1-d1
make push-tmdc-cluster-docker GITHUB_TAGS=3.3.1-d1
```

## Contributions

We welcome contributions! See how you can get involved by reading [CONTRIBUTING.md](./CONTRIBUTING.md).
