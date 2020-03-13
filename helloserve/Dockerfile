FROM golang:1.14-alpine
WORKDIR /helloserve
COPY . /helloserve
RUN ["go", "build"]

FROM alpine:latest
WORKDIR /helloserve
COPY --from=0 /helloserve/ /helloserve
EXPOSE 8100/tcp
CMD ["sh", "-c", "/helloserve/helloserve"]
