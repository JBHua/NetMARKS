# Testing Using k6.io

1. Install k6.io
https://k6.io/docs/get-started/installation/

2. Write the script
https://k6.io/docs/using-k6/http-requests/

3. Run it!
```shell
sh test.sh <serviceName>
```

# Testing Setup
## Environment
- minikube v1.31.2, used for local k8s clusters
- Grafana k6 v0.47.0, used for load-testing
- Mac Studio, 20 cores, 64 GB RAM; macOS 14.1.1

## k8s
For each node (5 nodes total):
 - Max 4 CPUs
 - Max 16 GB RAM

For Each Service (4 + 15 = 19 total):
 - 4 replicas for base service (services that has no dependencies, like fish & water)
 - 2 replicas for derived service (like boat, coin and sword)
 - 1 replicas for most complicated service (coin, iron, tools and sword)
 - each process (a pod) have 2 cpu allocated

## Mean Application Response Time
10 Requests total, should be complete under 60s. 1 second of sleep in between requests.
Test for each service is conducted independently (no parallel testing)

For each request, only ask the application to produce one result. Response size is changed upon request
Take the `avg` of `http_req_duration` field as result

Ignore failed request, or put it another way, set up the test so that they don't exceed maximum capacity

## Inter-Node Interactions 
