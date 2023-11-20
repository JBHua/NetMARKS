# NetMARKS Reproduction

## Project Structure

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

### 
Envoy sidecar


## How to Install Scheduler-plugins && Run your own scheduler
See:
https://github.com/kubernetes-sigs/scheduler-plugins/blob/master/doc/install.md

How to login minikube control plane:
https://minikube.sigs.k8s.io/docs/commands/ssh/

## TODO:
https://kubernetes.io/docs/concepts/cluster-administration/system-metrics/
