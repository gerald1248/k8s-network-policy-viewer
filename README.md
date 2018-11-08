k8s-network-policy-viewer
=========================

![Docker Automated](https://img.shields.io/docker/automated/gerald1248/k8s-network-policy-viewer.svg)
![Docker Build](https://img.shields.io/docker/build/gerald1248/k8s-network-policy-viewer.svg)

The network policy viewer visualizes the pod network. Many pieces are either unfinished or missing, but basic isolation rules can be represented in JSON, YAML or dot (Graphviz):

<img src="testdata/testdata.svg" alt="Sample visualization"/>

In this example, the names of the namespaces match their respective network policies, the exception being the `global` namespace which has none.

The policies `isolated`, `egress-isolated`, `ingress-isolated` each apply to the namespace as a whole.

`ingress-isolated-whitelist` whitelists `httpd-bob`, which is why `httpd-bob` can be reached from both `global` pods (including of course `httpd-alice`) and the `ingress-isolated` namespace.

Build
-----
The build steps are the following:
```
$ go mod download
$ go get
$ go vet
$ go test -v
$ go build -o k8s-network-policy-viewer .
```

Testdata
--------
To build the sample data, run:
```
$ make -C testdata init
$ make -C testdata create
```

Custom inputs
-------------
The application is intended for in-cluster use (the Helm chart with appropriate cluster role is in preparation), but you can use the application today by piping or supplying the output of `kubectl get po,cluster-policy --all-namespaces -o json`. The application accepts JSON and YAML, but you may wish to work with JSON so you can filter the input with `jq`.

Fuller coverage of pod selection, Helm chart, API and so on are all still in progress.
