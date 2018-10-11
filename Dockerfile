FROM alpine:3.8

RUN apk add --update --no-cache ca-certificates

ADD ./microkit-example /microkit-example

ENTRYPOINT ["/microkit-example"]

CMD [ "--server.listen.address=http://0.0.0.0:8000" ]
