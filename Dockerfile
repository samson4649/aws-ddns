FROM golang:1-alpine as builder

COPY src /go/src/gitlab.monarch-ares.io/devops/aws-ddns

RUN apk add --update --no-cache git \
  && cd /go/src/gitlab.monarch-ares.io/devops/aws-ddns/cmd \
  && go get -v \
  && go build -o /awsddns -ldflags='-w -s' main.go



FROM alpine:latest
MAINTAINER Samuel Lock <samuel.lock@monarch-ares.com>

COPY --from=builder /awsddns /awsddns

ENTRYPOINT ["/awsddns"]
CMD ["--config","/etc/awsddns/awsddns.yml"]
