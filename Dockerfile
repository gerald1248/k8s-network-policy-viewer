FROM golang:1.11 as builder
WORKDIR /go/src/github.com/gerald1248/k8s-network-policy-viewer/
COPY * ./
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on \
  go get && \
  go module save && \
  go vet -v && \
  go test -v && \
  go build -a -installsuffix cgo -o k8s-network-policy-viewer .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app/
USER 1000
COPY --from=builder /go/src/github.com/gerald1248/k8s-network-policy-viewer/k8s-network-policy-viewer .
CMD ["./k8s-network-policy-viewer", "-h"]  
