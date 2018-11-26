#!/bin/bash

# namespaces
cat << EOF >namespace-isolated.yaml
kind: Namespace
apiVersion: v1
metadata:
  name: isolated
  labels:
    app: isolated
spec:
  finalizers:
  - kubernetes
EOF

cat << EOF >namespace-global.yaml
kind: Namespace
apiVersion: v1
metadata:
  name: global
  labels:
    app: global
spec:
  finalizers:
  - kubernetes
EOF

cat << EOF >namespace-ingress-isolated.yaml
kind: Namespace
apiVersion: v1
metadata:
  name: ingress-isolated
  labels:
    app: ingress-isolated
spec:
  finalizers:
  - kubernetes
EOF

cat << EOF >namespace-egress-isolated.yaml
kind: Namespace
apiVersion: v1
metadata:
  name: egress-isolated
  labels:
    app: egress-isolated
spec:
  finalizers:
  - kubernetes
EOF

cat << EOF >namespace-ingress-isolated-whitelist.yaml
kind: Namespace
apiVersion: v1
metadata:
  name: ingress-isolated-whitelist
  labels:
    app: ingress-isolated-whitelist
spec:
  finalizers:
  - kubernetes
EOF

cat << EOF >namespace-eve.yaml
kind: Namespace
apiVersion: v1
metadata:
  name: eve
  labels:
    app: eve
spec:
  finalizers:
  - kubernetes
EOF

# deployments 
cat << EOF >deployment.yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: httpd
  labels:
    app: httpd
spec:
  replicas: 1
  selector:
    matchLabels:
      app: httpd
  template:
    metadata:
      labels:
        app: httpd
    spec:
      containers:
      - name: httpd
        image: centos/httpd-24-centos7
EOF

cat << EOF >pod-alice.yaml
apiVersion: v1
kind: Pod
metadata:
  name: httpd-alice
  labels:
    app: httpd-alice
spec:
  containers:
  - name: httpd
    image: centos/httpd-24-centos7
    ports:
    - containerPort: 8080
EOF

cat << EOF >pod-bob.yaml
apiVersion: v1
kind: Pod
metadata:
  name: httpd-bob
  labels:
    app: httpd-bob
spec:
  containers:
  - name: httpd
    image: centos/httpd-24-centos7
    ports:
    - containerPort: 8080
EOF

cat << EOF >pod-eve.yaml
apiVersion: v1
kind: Pod
metadata:
  name: httpd-eve
  labels:
    app: httpd-eve
spec:
  containers:
  - name: httpd
    image: centos/httpd-24-centos7
    ports:
    - containerPort: 8080
EOF

cat << EOF >network-policy-isolated.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: isolated
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
  ingress: []
  egress: []
EOF

cat << EOF >network-policy-ingress-whitelist-pod.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: ingress-whitelist-pod
spec:
  podSelector:
    matchLabels:
      app: httpd-bob
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: httpd-alice
EOF

cat << EOF >network-policy-ingress-whitelist-namespace.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: ingress-whitelist-namespace
spec:
  podSelector:
    matchLabels:
      app: httpd-bob
  policyTypes:
  - Ingress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          app: eve
EOF

cat << EOF >network-policy-ingress-isolated.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: ingress-isolated
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  ingress: []
EOF

cat << EOF >network-policy-egress-isolated.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: egress-isolated
spec:
  podSelector: {}
  policyTypes:
  - Egress
  egress: []
EOF

for NAMESPACE in isolated global ingress-isolated egress-isolated ingress-isolated-whitelist eve; do
  kubectl create -f namespace-${NAMESPACE}.yaml
  kubectl create -f deployment.yaml -n ${NAMESPACE}
done

kubectl create -f pod-bob.yaml -n ingress-isolated-whitelist
kubectl create -f pod-alice.yaml -n ingress-isolated-whitelist

kubectl create -f network-policy-isolated.yaml -n eve
kubectl create -f pod-eve.yaml -n eve

kubectl create -f network-policy-isolated.yaml -n isolated
kubectl create -f network-policy-ingress-isolated.yaml -n ingress-isolated
kubectl create -f network-policy-egress-isolated.yaml -n egress-isolated
kubectl create -f network-policy-ingress-isolated.yaml -n ingress-isolated-whitelist
kubectl create -f network-policy-ingress-whitelist-pod.yaml -n ingress-isolated-whitelist
kubectl create -f network-policy-ingress-whitelist-namespace.yaml -n ingress-isolated-whitelist
