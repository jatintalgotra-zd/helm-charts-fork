# ChromaDB Helm Chart

The ChromaDB Helm chart provides an easy way to deploy ChromaDB, a high-performance embedding database designed for AI applications. This chart allows you to manage ChromaDB instances on Kubernetes with customizable resource allocation, persistence, and scaling options.

---

## Prerequisites

- Kubernetes 1.19+  
- Helm 3+

---

## Add Helm Repository

Before deploying the ChromaDB chart, add the Helm repository to your local setup:

```bash
helm repo add zopdev https://helm.zop.dev
helm repo update
```

For more details, refer to the [Helm Repository Documentation](https://helm.sh/docs/helm/helm_repo/).

---

## Install Helm Chart

To install the ChromaDB Helm chart, use the following command:

```bash
helm install [RELEASE_NAME] zopdev/chromadb
```

Replace `[RELEASE_NAME]` with your desired release name. For example:

```bash
helm install my-chromadb zopdev/chromadb
```

To customize configurations, provide a `values.yaml` file or override values via the command line.

See [Helm Install Documentation](https://helm.sh/docs/helm/helm_install/) for more information.

---

## Uninstall Helm Chart

To remove the ChromaDB deployment and all associated Kubernetes resources, use the following command:

```bash
helm uninstall [RELEASE_NAME]
```

For example:

```bash
helm uninstall my-chromadb
```

See [Helm Uninstall Documentation](https://helm.sh/docs/helm/helm_uninstall/) for additional details.

---

## Configuration

The ChromaDB Helm chart includes several configuration options to tailor the deployment to your needs. Below is a summary of the key configurations:

| **Input**               | **Type**  | **Description**                                                                                | **Default**           |
|--------------------------|-----------|------------------------------------------------------------------------------------------------|-----------------------|
| `image`                  | `string`  | Docker image and tag for the ChromaDB container.                                               | `ghcr.io/chroma-core/chroma:latest` |
| `resources.requests.memory` | `string` | Minimum memory resources required by the ChromaDB container.                                   | `"1Gi"`              |
| `resources.requests.cpu` | `string` | Minimum CPU resources required by the ChromaDB container.                                      | `"100m"`             |
| `resources.limits.memory` | `string` | Maximum memory resources the ChromaDB container can use.                                       | `"2Gi"`              |
| `resources.limits.cpu`   | `string`  | Maximum CPU resources the ChromaDB container can use.                                          | `"1"`                |
| `diskSize`               | `string`  | Size of the persistent volume for storing ChromaDB data.                                       | `"10Gi"`             |

You can override these values in a `values.yaml` file or via the command line during installation.

---

### Example `values.yaml` File

```yaml
version: 0.6.3

resources:
  requests:
    memory: "1Gi"
    cpu: "100m"
  limits:
    memory: "2Gi"
    cpu: "1000m"

diskSize: "10Gi"
```

To use this configuration, save it to a `values.yaml` file and apply it during installation:

```bash
helm install my-chromadb zopdev/chromadb -f values.yaml
```

---

## Features

- **Scalable Architecture:** Configure resources and scaling options to optimize performance for embedding-intensive workloads.
- **Persistent Storage:** Keep ChromaDB data intact across pod restarts with configurable persistent volumes.
- **Customizable Resource Allocation:** Tailor CPU and memory resources to match workload requirements.
- **Easy Deployment:** Simplified Helm chart for rapid deployment of ChromaDB in Kubernetes environments.

---

## Contributing

We welcome contributions to improve this Helm chart. Please refer to the [CONTRIBUTING.md](../../CONTRIBUTING.md) file for contribution guidelines.

---

## Code of Conduct

To maintain a healthy and collaborative community, please adhere to our [Code of Conduct](../../CODE_OF_CONDUCT.md).

---

## License

This project is licensed under the [LICENSE](../../LICENSE). Please review it for terms of use.

---

## Connection Config

- **CHROMA_SERVER_HOST** : Hostname or service name for the ChromaDB instance.
- **CHROMA_SERVER_HTTP_PORT** : HTTP port on which ChromaDB listens. Default is 8000 .
- **ANONYMIZED_TELEMETRY** : Enables/disables sending anonymized telemetry.False means disabled.
- **ALLOW_RESET** : Allows resetting ChromaDB data if set to true.
- **IS_PERSISTENT** : Determines if ChromaDB data should be persistent.
- **CHROMA_LOG_LEVEL** : Log level for ChromaDB .
- **CHROMA_DATA_DIR** : Directory path where ChromaDB stores its data.
- **CHROMA_IMPORT_SETTING** : Indicates whether import settings are enabled .

---

