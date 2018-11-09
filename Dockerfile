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

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app/
USER 1000
COPY --from=builder /go/src/github.com/gerald1248/k8s-network-policy-viewer/k8s-network-policy-viewer /usr/bin/
CMD ["k8s-network-policy-viewer", "-s=true"]  
