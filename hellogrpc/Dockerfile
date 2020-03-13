FROM golang:1.14-alpine
WORKDIR /hellogrpc
COPY . /hellogrpc
RUN ["go", "build"]

FROM alpine:latest
WORKDIR /hellogrpc
COPY --from=0 /hellogrpc/ /hellogrpc
EXPOSE 8000/tcp
CMD ["sh", "-c", "/hellogrpc/hellogrpc"]
