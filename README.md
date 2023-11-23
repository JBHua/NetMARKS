# NetMARKS Reproduction

## Project Structure

## Setup:
### Dev Environment
1. Install Dependencies
    - Container manager like Docker, or Podman
    - minikube: https://minikube.sigs.k8s.io/docs/start/
    - k8s version: v1.27.4
2. Configure Minikube
    - See: https://minikube.sigs.k8s.io/docs/commands/config/
3. Start Local k8s cluster using minikube
   ```shell
   minikube start --nodes <nodeCount>
   ```

### minikube
When initializing minikube nodes, only some core dependencies will be installed (dashboard and dns). For the experiment,
more setup is needed:

1. Istio Service Mesh
    - https://istio.io/latest/docs/setup/getting-started/ 
    - Remember to enable Envoy sidecar injection!
2. Install Prometheus & Grafana (optional) using Istio
   - https://istio.io/latest/docs/ops/integrations/prometheus/
   - https://istio.io/latest/docs/ops/integrations/grafana
4. Enable Scheduler
    - Custom: See following section


### Services
Use `build-images.sh` under `<rootDir>/services/` to build mock services and push the images to minikube.
Use `deployment.yaml` under `<rootDir>/infra` to create deployment. Use `services.yaml` under `<rootDir>/infra` to create services

## Deployment 
### Intro
We are doing dev locally, which means we don't push our test service's image to cloud. Thus, we need to load image:
    See: https://minikube.sigs.k8s.io/docs/handbook/pushing/#7-loading-directly-to-in-cluster-container-runtime
```shell
minikube image load <imageName>
```

### Deploy a Pod

```shell
kubectl apply -f infra/pods.yaml
```

### Expose/Forward 
```shell
kubectl port-forward svc/prometheus 9090:9090 -n istio-system
```

### Get All Pods
```shell
kubectl get po -A
```

### 
Envoy sidecar


## How to Install Scheduler-plugins && Run your own scheduler
See:
https://github.com/kubernetes-sigs/scheduler-plugins/blob/master/doc/install.md

How to login minikube control plane:
https://minikube.sigs.k8s.io/docs/commands/ssh/

## TODO:
https://kubernetes.io/docs/concepts/cluster-administration/system-metrics/
