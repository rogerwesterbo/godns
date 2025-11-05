# GoDNS Web UI Helm Chart

This Helm chart deploys the GoDNS Web UI - a React-based web interface for managing DNS zones and records.

## Features

- 🔐 **Secure by default**: Runs as non-root user with read-only root filesystem
- 🐳 **Minimal footprint**: Based on Google Distroless for maximum security
- 📊 **Auto-scaling**: Optional HPA support
- 🔒 **Network policies**: Optional network isolation
- 🚀 **High availability**: Pod disruption budgets and anti-affinity

## Prerequisites

- Kubernetes 1.20+
- Helm 3.0+
- (Optional) Ingress controller for external access
- (Optional) cert-manager for TLS certificates

## Installation

### Add the Helm repository

```bash
helm repo add godns oci://ghcr.io/rogerwesterbo/helm
```

### Install the chart

```bash
helm install godnsweb oci://ghcr.io/rogerwesterbo/helm/godnsweb \
  --namespace godns \
  --create-namespace
```

### Install with custom values

```bash
helm install godnsweb oci://ghcr.io/rogerwesterbo/helm/godnsweb \
  --namespace godns \
  --create-namespace \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=godns.example.com \
  --set env[0].name=VITE_API_BASE_URL \
  --set env[0].value=https://api.godns.example.com
```

## Configuration

The following table lists the configurable parameters of the GoDNS Web chart and their default values.

| Parameter                   | Description           | Default                           |
| --------------------------- | --------------------- | --------------------------------- |
| `replicaCount`              | Number of replicas    | `2`                               |
| `image.repository`          | Image repository      | `ghcr.io/rogerwesterbo/godns-web` |
| `image.tag`                 | Image tag             | `""` (uses chart appVersion)      |
| `image.pullPolicy`          | Image pull policy     | `IfNotPresent`                    |
| `service.type`              | Service type          | `ClusterIP`                       |
| `service.port`              | Service port          | `80`                              |
| `ingress.enabled`           | Enable ingress        | `false`                           |
| `ingress.className`         | Ingress class name    | `nginx`                           |
| `ingress.hosts`             | Ingress hosts         | `[godns-web.example.com]`         |
| `resources.limits.cpu`      | CPU limit             | `200m`                            |
| `resources.limits.memory`   | Memory limit          | `256Mi`                           |
| `resources.requests.cpu`    | CPU request           | `100m`                            |
| `resources.requests.memory` | Memory request        | `128Mi`                           |
| `autoscaling.enabled`       | Enable HPA            | `false`                           |
| `autoscaling.minReplicas`   | Minimum replicas      | `2`                               |
| `autoscaling.maxReplicas`   | Maximum replicas      | `10`                              |
| `env`                       | Environment variables | See values.yaml                   |

## Environment Variables

The web application can be configured with the following environment variables:

- `VITE_KEYCLOAK_URL`: Keycloak server URL
- `VITE_KEYCLOAK_REALM`: Keycloak realm name
- `VITE_KEYCLOAK_CLIENT_ID`: Keycloak client ID
- `VITE_REDIRECT_URI`: OAuth2 redirect URI
- `VITE_POST_LOGOUT_REDIRECT_URI`: Post-logout redirect URI
- `VITE_API_BASE_URL`: GoDNS API base URL

## Examples

### Basic installation

```bash
helm install godnsweb oci://ghcr.io/rogerwesterbo/helm/godnsweb
```

### With ingress enabled

```bash
cat <<EOF | helm install godnsweb oci://ghcr.io/rogerwesterbo/helm/godnsweb -f -
ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: godns.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: godnsweb-tls
      hosts:
        - godns.example.com

env:
  - name: VITE_KEYCLOAK_URL
    value: "https://keycloak.example.com"
  - name: VITE_API_BASE_URL
    value: "https://api.godns.example.com"
EOF
```

### With autoscaling

```bash
helm install godnsweb oci://ghcr.io/rogerwesterbo/helm/godnsweb \
  --set autoscaling.enabled=true \
  --set autoscaling.minReplicas=3 \
  --set autoscaling.maxReplicas=20
```

### With network policies

```bash
helm install godnsweb oci://ghcr.io/rogerwesterbo/helm/godnsweb \
  --set networkPolicy.enabled=true
```

## Upgrading

```bash
helm upgrade godnsweb oci://ghcr.io/rogerwesterbo/helm/godnsweb \
  --namespace godns
```

## Uninstalling

```bash
helm uninstall godnsweb --namespace godns
```

## Security

This chart follows security best practices:

- Runs as non-root user (UID 65532)
- Uses read-only root filesystem
- Drops all Linux capabilities
- Uses Google Distroless base image (no shell, no package manager)
- Implements security context constraints
- Supports network policies for pod isolation
- Implements pod disruption budgets for high availability

## Support

For issues and questions, please visit:

- GitHub Issues: https://github.com/rogerwesterbo/godns/issues
- Documentation: https://github.com/rogerwesterbo/godns/tree/main/docs

## License

See the [LICENSE](https://github.com/rogerwesterbo/godns/blob/main/LICENSE) file for details.
