FROM golang:1.18-alpine as build
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

RUN apk add --no-cache make git
COPY . .
RUN make build

FROM alpine:3.14.2 as certs
RUN apk update && apk add ca-certificates

FROM gcr.io/distroless/base
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=build /build/bin/slsa-provenance /usr/local/bin/slsa-provenance
# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["slsa-provenance"]
