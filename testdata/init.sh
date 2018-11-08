#!/bin/bash

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

cat << EOF >network-policy-isolated.yaml
apiVersion: v1
items:
- apiVersion: extensions/v1beta1
  kind: NetworkPolicy
  metadata:
    name: isolated
  spec:
    podSelector: {}
    policyTypes:
    - Ingress
    - Egress
kind: List
metadata:
  name: network-policy-isolated
EOF

cat << EOF >network-policy-ingress-isolated-whitelist.yaml
apiVersion: v1
items:
- apiVersion: extensions/v1beta1
  kind: NetworkPolicy
  metadata:
    name: ingress-isolated-whitelist
  spec:
    podSelector:
      matchLabels:
        app: httpd-bob
    policyTypes:
    - Ingress
    ingress:
    - {}
kind: List
metadata:
  name: network-policy-ingress-isolated-whitelist
EOF

cat << EOF >network-policy-ingress-isolated.yaml
apiVersion: v1
items:
- apiVersion: extensions/v1beta1
  kind: NetworkPolicy
  metadata:
    name: ingress-isolated
  spec:
    podSelector: {}
    policyTypes:
    - Ingress
kind: List
metadata:
  name: network-policy-ingress-isolated
EOF

cat << EOF >network-policy-egress-isolated.yaml
apiVersion: v1
items:
- apiVersion: extensions/v1beta1
  kind: NetworkPolicy
  metadata:
    name: egress-isolated
  spec:
    podSelector: {}
    policyTypes:
    - Egress
kind: List
metadata:
  name: network-policy-egress-isolated
EOF

for NAMESPACE in isolated global ingress-isolated egress-isolated ingress-isolated-whitelist; do
  kubectl create namespace ${NAMESPACE} 
  kubectl create -f deployment.yaml -n ${NAMESPACE}
done

kubectl create -f pod-bob.yaml -n ingress-isolated-whitelist
kubectl create -f pod-alice.yaml -n global

kubectl create -f network-policy-isolated.yaml -n isolated
kubectl create -f network-policy-ingress-isolated.yaml -n ingress-isolated
kubectl create -f network-policy-egress-isolated.yaml -n egress-isolated
kubectl create -f network-policy-ingress-isolated-whitelist.yaml -n ingress-isolated-whitelist

