# Bitnami Production Runtime for Kubernetes

The Bitnami Kubernetes Production Runtime (BKPR) is a collection of services that makes it easy to run production workloads in Kubernetes.

Think of BKPR as a curated collection of the services you would need to deploy on top of your Kubernetes cluster to enable logging, monitoring, certificate management, automatic discovery of Kubernetes resources via public DNS servers and other common infrastructure needs.

BKPR is available for [Google Kubernetes Engine (GKE)](https://cloud.google.com/kubernetes-engine) and [Azure Kubernetes Service (AKS)](https://azure.microsoft.com/en-in/services/kubernetes-service/) clusters.

## Quickstart

* [AKS Quickstart](docs/quickstart-aks.md)
* GKE Quickstart

## Frequently Asked Questions (FAQ)

Please read the [FAQ](docs/FAQ.md).

## Versioning

The versioning used in BKPR is described [here](docs/versioning.md).

## Components

### Monitoring
* [Prometheus](https://prometheus.io/): A monitoring system and time series database
* [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/): An alert manager and router
### Logging
* [Elasticsearch](https://www.elastic.co/products/elasticsearch): A distributed, RESTful search and analytics engine
* [Kibana](https://www.elastic.co/products/kibana): A visualization tool for Elasticsearch data
* [Fluentd](https://www.fluentd.org/): A data collector for unified logging layer
### DNS and TLS certificates
* [ExternalDNS](https://github.com/kubernetes-incubator/external-dns): A component to synchronize exposed Kubernetes Services and Ingresses with DNS providers
* [cert-manager](https://github.com/jetstack/cert-manager): A Kubernetes add-on to automate the management and issuance of TLS certificates from various sources
### Others
* [OAuth2 Proxy](https://github.com/bitnami/bitnami-docker-oauth2-proxy): A reverse proxy and static file server that provides authentication using Providers (Google, GitHub, and others) to validate accounts by email, domain or group
* [nginx-ingress](https://github.com/kubernetes/ingress-nginx): A Controller to satisfy requests for Ingress objects

## Release Compatibility

### Kubernetes Version Support

The following matrix shows which Kubernetes versions are supported in AKS and GKE in the current releases of BKPR:

| BKPR release | Kubernetes version | AKS support? | GKE support? |
|:------------:|:------------------:|:------------:|:------------:|
|      1.0     |         1.8        |      Yes     |      Yes     |
|      1.0     |         1.9        |      Yes     |      Yes     |

### Components Version Support

The following matrix shows which versions of each component are used and supported in the curerent releases of BKPR:

| BKPR release |   Component   |          Release |
|:------------:|:-------------:|-----------------:|
|      1.0     |   Prometheus  |            2.3.2 |
|      1.0     |     Kibana    |           5.6.12 |
|      1.0     | Elasticsearch |           5.6.12 |
|      1.0     |  cert-manager |            0.3.2 |
|      1.0     |  alertmanager |           0.15.2 |
|      1.0     |  ExternalDNS  |            0.5.4 |
|      1.0     | nginx-ingress |           0.19.0 |
|      1.0     |  oauth2_proxy | 0.20180625.74543 |
|      1.0     |    heapster   |            1.5.2 |
|      1.0     |    fluentd    |            1.2.2 |