FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN apk add --no-cache make build-base
RUN make build

FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /app/k8s-controller .
EXPOSE 8080
ENTRYPOINT ["/k8s-controller"]