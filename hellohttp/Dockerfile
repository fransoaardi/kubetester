FROM golang:1.14-alpine
WORKDIR /hellohttp
COPY . /hellohttp
RUN ["go", "build"]

FROM alpine:latest
WORKDIR /hellohttp
COPY --from=0 /hellohttp/ /hellohttp
EXPOSE 8100/tcp
CMD ["sh", "-c", "/hellohttp/hellohttp"]
