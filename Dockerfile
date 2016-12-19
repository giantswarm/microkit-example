FROM alpine:3.4

RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

ADD ./microkit-example /microkit-example

ENTRYPOINT ["/microkit-example"]

CMD [ "--server.listen.address=http://0.0.0.0:8000" ]
