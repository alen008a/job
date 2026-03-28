FROM alpine:latest
ARG config
COPY ./logs /data/siteLetterJob/logs
COPY ./config /data/siteLetterJob/config
COPY ./siteLetterJob /data/siteLetterJob/siteLetterJob
WORKDIR /data/siteLetterJob
VOLUME ["/data/siteLetterJob/logs"]

RUN apk update \
    && apk upgrade \
    && apk add --no-cache  \
    && apk add ca-certificates \
    && apk add update-ca-certificates 2>/dev/null || true \
    && apk add tzdata \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && mv ${config} ./config/active.ini

CMD ["/data/siteLetterJob/siteLetterJob", "--config=./config/active.ini"]