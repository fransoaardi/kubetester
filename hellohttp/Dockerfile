FROM golang:1.14-alpine
WORKDIR /hellohttp
COPY . /hellohttp
RUN ["go", "build"]
EXPOSE 8100/tcp
CMD ["sh", "-c", "/hellohttp/hellohttp"]
