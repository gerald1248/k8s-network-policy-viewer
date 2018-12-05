package main

import (
	"testing"
)

func TestProcessFileInvalidPath(t *testing.T) {
	invalidPath := "/non/existent/file.yaml"
	output := "dot"
	_, err := processFile(invalidPath, &output)

	if err == nil {
		t.Errorf("Must reject invalid path %s", invalidPath)
	}
}

func TestProcessBytes(t *testing.T) {
	//don't allow XML
	xmlBuffer := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="true"?><root/>`)
	output := "dot"
	_, _, _, err := processBytes(xmlBuffer, &output)

	if err == nil {
		t.Errorf("Must reject XML input")
	}
}

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
