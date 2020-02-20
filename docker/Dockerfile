FROM golang:1.13-alpine AS builder
ARG SVC
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/mainflux/agent
COPY . .
RUN apk update \
    && apk add make\
    && make $SVC \
    && mv build/mainflux-$SVC /exe

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
