FROM golang:1.14-alpine
WORKDIR /helloserve
COPY . /helloserve
RUN ["go", "build"]
EXPOSE 8100/tcp
CMD ["sh", "-c", "/helloserve/helloserve"]
