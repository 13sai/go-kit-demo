FROM alpine:3.12 
RUN mkdir -p /user
WORKDIR /user
COPY ./user /user
COPY ./docker.yaml /user
ENTRYPOINT ["./user", "-c", "/user/docker.yaml"]