FROM golang:1.12.0 as builder
WORKDIR /go/src/github.com/gerald1248/k8s-network-policy-viewer/
COPY * ./
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GO111MODULE on
RUN \
  go mod download && \
  go get && \
  go vet && \
  go test -v -cover && \
  go build -o k8s-network-policy-viewer .

FROM ubuntu:18.10
WORKDIR /app/
EXPOSE 8080
ENV NETWORK_POLICY_VIEWER_BLACKLIST default,kube,flux
RUN apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get -qq install curl graphviz
COPY --from=builder /go/src/github.com/gerald1248/k8s-network-policy-viewer/k8s-network-policy-viewer /usr/bin/
USER 1000
CMD ["k8s-network-policy-viewer", "-s=true"]
