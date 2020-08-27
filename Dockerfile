FROM golang:1.14 as builder

WORKDIR /go/src/k8s-audit-webhook
COPY . .

RUN go get -d -v ./cmd/k8s-audit-webhook
RUN GOOS=linux CGO_ENABLED=0 go install -v -a -ldflags '-extldflags "-static"' ./cmd/k8s-audit-webhook

FROM busybox:1.32 as runtime

COPY --from=builder /go/bin/k8s-audit-webhook /k8s-audit-webhook
CMD ["/k8s-audit-webhook"]
