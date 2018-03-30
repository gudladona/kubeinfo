# kubeinfo

This HTTP service authenticates to the Kubernetes API from an application running inside the Kubernetes cluster and returns metadata about the pods.

## Assumptions

* Audience is running a Unix/Linux OS

* Audience has basic knowledge of Go, Docker, Kubernetes, AWS

* All necessary software(s) are installed on audience OS. viz. kubectl, make, Go, Docker

* Respective Paths viz. GOPATH are properly set on audience OS.

* A Kubernetes cluster is running in AWS

* Kubernetes cluster can access public Docker Repo

* Audience following this guide has necessary access to the kube cluster and proper kube-context set.

* The security groups for AWS has ports 30000-32767 open for audience IP Address/Block.

## Installation

```bash
go get github.com/gudladona87/kubeinfo
```

## Installation on Kubernetes

1. `export DOCKER_USER = "gudladona87"`

2. `docker login --username gudladona87 --password-stdin`     -- Password provided in the email

3. Install make

4. `make deploy`

## Running tests

    > make test

Current code coverage available at ./coverage/coverage.html

## Testing the API End to End

For the sake of simplicity the `type` in .kube/service.yml is set to `NodePort`,
which allows each Node to proxy that port (the same port number on every Node) into the Service.

Once this service is installed on a Kubernetes cluster,

* Grep for the Service `kubeinfo`. Copy the port the service is exposed on.

    ```
    > kubectl get svc | grep kubeinfo
    kubeinfo    NodePort    100.64.219.181  <none> 8080:31046/TCP  27m
    ```
* Find the pod that the service is running on. Although the same port is exposed on all the nodes,
for simplicity I choose to use the node where the pod is running.

    ```
    > kubectl get po | grep kube
    kubeinfo-7c854d54d8-vp2tw   1/1 Running   0 30m
    ```

* Find the Node that the pod is running on
    ```
    > kubectl describe pod kubeinfo-7c854d54d8-vp2tw | grep 'Node:'
    Node:           ip-10-75-24-178.ec2.internal/10.75.24.178
    ```
* Initiate a GET request on the Node and the port at the endpoint `/pods`
    ```
    > curl http://10.75.24.178:31046/pods
    {"pod_count":314,"message":"OK"}
    ```

Note:

A better approach would be to configure an ingress for the service at a registered domain such as kubeinfo.example.com.
But, this requires registering a domain, which I think is out of scope.