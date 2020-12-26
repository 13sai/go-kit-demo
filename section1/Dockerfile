FROM alpine:3.12 
WORKDIR /
COPY ./user /
COPY ./config.yaml /
ENTRYPOINT ["./user", "-c", "./config.yaml"]