FROM alpine:3.15

WORKDIR /opt/mongo-keys

RUN apk add openssl

COPY run.sh run.sh
RUN chmod +x run.sh

ENTRYPOINT ["sh", "-c", "/opt/mongo-keys/run.sh"]
