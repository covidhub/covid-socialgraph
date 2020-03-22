FROM golang:1.14.1-alpine3.11 AS build

RUN apk add --no-cache curl g++ openssl-dev openssl-libs-static
RUN curl -qL -o - \https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-alpine-3.9.3.tar.gz \
    | tar -C / -zxv --strip-components=1

WORKDIR /socialgraph
COPY ./go.mod ./go.sum ./
RUN GOPROXY=https://proxy.golang.org go mod download
COPY ./ ./
RUN go build -o /go/bin/socialgraph -tags seabolt_static /socialgraph/cmd/socialgraph

FROM alpine:latest

COPY --from=build /go/bin/socialgraph ./socialgraph

CMD "./socialgraph"
