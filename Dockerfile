FROM alpine:3.7

RUN apk add --no-cache tini

COPY alarmclock /
COPY docker/docker-entrypoint.sh /

EXPOSE 5050
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/alarmclock", "-f", "config.yml"]
