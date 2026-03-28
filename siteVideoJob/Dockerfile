FROM alpine:latest
ARG config
COPY ./logs /data/siteVideoJob/logs
COPY ./config /data/siteVideoJob/config
COPY ./siteVideoJob /data/siteVideoJob/siteVideoJob
WORKDIR /data/siteVideoJob
VOLUME ["/data/siteVideoJob/logs"]

RUN apk update \
    && apk upgrade \
    && apk add --no-cache  \
    && apk add ca-certificates \
    && apk add update-ca-certificates 2>/dev/null || true \
    && apk add tzdata \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && mv ${config} ./config/active.ini

CMD ["/data/siteVideoJob/siteVideoJob", "--config=./config/active.ini"]