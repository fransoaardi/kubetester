```
protoc -I hellogrpc/ proto/hello.proto --go_out=plugins=grpc:hellogrpc
```