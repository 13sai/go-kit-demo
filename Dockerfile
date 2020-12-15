FROM alpine:3.12 
RUN mkdir -p /register
WORKDIR /register
COPY ./register /register
EXPOSE 12312
ENTRYPOINT ./register -consul.addr=$consulAddr -service.addr=$serviceAddr