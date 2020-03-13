FROM golang:1.14-alpine
WORKDIR /hellogrpc
COPY . /hellogrpc
RUN ["go", "build"]
EXPOSE 8000/tcp
CMD ["sh", "-c", "/hellogrpc/hellogrpc"]
