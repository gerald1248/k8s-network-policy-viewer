FROM golang:1.11.1 as builder
WORKDIR /go/src/github.com/gerald1248/k8s-network-policy-viewer/
COPY * ./
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GO111MODULE on
RUN \
  #go mod tidy && \
  go mod download && \
  go get && \
  go vet && \
  go test -v && \
  go build -o k8s-network-policy-viewer .

FROM ubuntu:18.10
WORKDIR /app/
EXPOSE 8080
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get -qq install curl
COPY --from=builder /go/src/github.com/gerald1248/k8s-network-policy-viewer/k8s-network-policy-viewer /usr/bin/
USER 1000
CMD ["k8s-network-policy-viewer", "-s=true"]  
