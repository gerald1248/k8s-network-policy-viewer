package main

import (
	"testing"
)

func TestIsolated(t *testing.T) {
	inputYaml := []byte(`
apiVersion: v1
kind: List
items:
- kind: Namespace
  apiVersion: v1
  metadata:
    name: alice
    labels:
      app: alice
  spec:
    finalizers:
    - kubernetes
- kind: Namespace
  apiVersion: v1
  metadata:
    name: bob
    labels:
      app: bob
  spec:
    finalizers:
    - kubernetes
- apiVersion: v1
  kind: Pod
  metadata:
    name: alice
    namespace: alice
    labels:
      app: alice
  spec:
    containers:
    - name: httpd
      image: centos/httpd-24-centos7
      ports:
      - containerPort: 8080
  status:
    containerStatuses:
    - ready: true
- apiVersion: v1
  kind: Pod
  metadata:
    name: bob
    namespace: bob
    labels:
      app: bob
  spec:
    containers:
    - name: httpd
      image: centos/httpd-24-centos7
      ports:
      - containerPort: 8080
  status:
    containerStatuses:
    - ready: true
- apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: isolated
    namespace: alice
  spec:
    podSelector: {}
    policyTypes:
    - Ingress
    - Egress
    ingress: []
    egress: []
- apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: isolated
    namespace: bob
  spec:
    podSelector: {}
    policyTypes:
    - Ingress
    - Egress
    ingress: []
    egress: []
`)
	output := "dot"
	_, isolation, coverage, err := processBytes(inputYaml, &output)

	if err != nil {
		t.Errorf("Must accept input YAML")
		return
	}

	if isolation < 100 || coverage < 100 {
		t.Errorf("Must recognise full isolation network policy - got isolation=%d, coverage=%d", isolation, coverage)
	}
}

func TestNonIsolated(t *testing.T) {
	inputYaml := []byte(`
apiVersion: v1
kind: List
items:
- kind: Namespace
  apiVersion: v1
  metadata:
    name: alice
    labels:
      app: alice
  spec:
    finalizers:
    - kubernetes
- kind: Namespace
  apiVersion: v1
  metadata:
    name: bob
    labels:
      app: bob
  spec:
    finalizers:
    - kubernetes
- apiVersion: v1
  kind: Pod
  metadata:
    name: alice
    namespace: alice
    labels:
      app: alice
  spec:
    containers:
    - name: httpd
      image: centos/httpd-24-centos7
      ports:
      - containerPort: 8080
  status:
    containerStatuses:
    - ready: true
- apiVersion: v1
  kind: Pod
  metadata:
    name: bob
    namespace: bob
    labels:
      app: bob
  spec:
    containers:
    - name: httpd
      image: centos/httpd-24-centos7
      ports:
      - containerPort: 8080
  status:
    containerStatuses:
    - ready: true
`)
	output := "dot"
	_, isolation, coverage, err := processBytes(inputYaml, &output)

	if err != nil {
		t.Errorf("Must accept input YAML")
		return
	}

	if isolation > 0 || coverage > 0 {
		t.Errorf("Must recognise absence of network policy - got isolation=%d, coverage=%d", isolation, coverage)
	}
}

func TestAsymmetricalWhitelist(t *testing.T) {
	inputYaml := []byte(`
apiVersion: v1
kind: List
items:
- kind: Namespace
  apiVersion: v1
  metadata:
    name: alice
    labels:
      app: alice
  spec:
    finalizers:
    - kubernetes
- kind: Namespace
  apiVersion: v1
  metadata:
    name: bob
    labels:
      app: bob
  spec:
    finalizers:
    - kubernetes
- apiVersion: v1
  kind: Pod
  metadata:
    name: alice
    namespace: alice
    labels:
      app: alice
  spec:
    containers:
    - name: httpd
      image: centos/httpd-24-centos7
      ports:
      - containerPort: 8080
  status:
    containerStatuses:
    - ready: true
- apiVersion: v1
  kind: Pod
  metadata:
    name: bob
    namespace: bob
    labels:
      app: bob
  spec:
    containers:
    - name: httpd
      image: centos/httpd-24-centos7
      ports:
      - containerPort: 8080
  status:
    containerStatuses:
    - ready: true
- apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: isolated
    namespace: alice
  spec:
    podSelector: {}
    policyTypes:
    - Ingress
    - Egress
    ingress: []
    egress: []
- apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: isolated
    namespace: bob
  spec:
    podSelector: {}
    policyTypes:
    - Ingress
    - Egress
    ingress: []
    egress: []
- apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: whitelist
    namespace: bob
  spec:
    podSelector:
      matchLabels:
        app: bob
    policyTypes:
    - Ingress
    ingress:
    - from:
      - podSelector:
          matchLabels:
            app: alice
        namespaceSelector:
          matchLabels:
            app: alice
  `)
	output := "dot"
	_, isolation, coverage, err := processBytes(inputYaml, &output)

	if err != nil {
		t.Errorf("Must accept input YAML")
		return
	}

	if isolation != 50 || coverage != 100 {
		t.Errorf("Must recognise asymmetrical isolation - got isolation=%d, coverage=%d", isolation, coverage)
	}
}

func TestIngress(t *testing.T) {
	inputYaml := []byte(`
apiVersion: v1
kind: List
items:
- kind: Namespace
  apiVersion: v1
  metadata:
    name: alice
    labels:
      app: alice
  spec:
    finalizers:
    - kubernetes
- kind: Namespace
  apiVersion: v1
  metadata:
    name: ingress
    labels:
      app: ingress
  spec:
    finalizers:
    - kubernetes
- apiVersion: v1
  kind: Pod
  metadata:
    name: alice
    namespace: alice
    labels:
      app: alice
  spec:
    containers:
    - name: httpd
      image: centos/httpd-24-centos7
      ports:
      - containerPort: 8080
  status:
    containerStatuses:
    - ready: true
- apiVersion: v1
  kind: Pod
  metadata:
    name: ingress
    namespace: ingress
    labels:
      app: ingress
  spec:
    containers:
    - name: httpd
      image: centos/httpd-24-centos7
      ports:
      - containerPort: 8080
  status:
    containerStatuses:
    - ready: true
- apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: isolated
    namespace: alice
  spec:
    podSelector: {}
    policyTypes:
    - Ingress
    - Egress
    ingress: []
    egress: []
- apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: isolated
    namespace: alice
  spec:
    podSelector: {}
    policyTypes:
    - Ingress
    - Egress
    ingress: []
    egress: []
- apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: whitelist
    namespace: alice
  spec:
    podSelector:
      matchLabels:
        app: alice
    policyTypes:
    - Ingress
    ingress:
    - from:
      - podSelector:
          matchLabels:
            app: ingress
        namespaceSelector:
          matchLabels:
            app: ingress
  `)
	output := "dot"
	_, isolation, coverage, err := processBytes(inputYaml, &output)

	if err != nil {
		t.Errorf("Must accept input YAML")
		return
	}

	if isolation != 50 || coverage != 50 {
		t.Errorf("Must accept ingress connections - got isolation=%d, coverage=%d", isolation, coverage)
	}
}
