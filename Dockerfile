FROM alpine:latest
MAINTAINER Samuel Lock <samuel.lock@monarch-ares.com>

RUN apk add --no-cache --update --virtual .build-deps
RUN apk add --no-cache bash make curl openssh git nodejs yarn
RUN apk -Uuv add groff less python py-pip

RUN pip install awscli

RUN apk --purge -v del py-pip

RUN rm /var/cache/apk/*

COPY ./template.json /template.json
COPY ./update.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]



