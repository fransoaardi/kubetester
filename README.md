# kubetester

## intro
- kubernetes gRPC, HTTP loadbalancing 테스트를 위한 시스템 구성
- 마이크로서비스 인 액션 (모건 브루스, 파울로 페레이라) 9장 컨테이너와 스케줄러를 이용해 배포하기를 읽고 테스트 구현해봄  
- gRPC loadbalance 는 꼭 `Headless Service` 말고 `linkerd` 등 application layer load balancer, 혹은 service mesh 구성이 있겠지만, 일단은 구성해봄 
  - https://kubernetes.io/blog/2018/11/07/grpc-load-balancing-on-kubernetes-without-tears/

## troubleshoot 
- `Service` 에 `NodePort` 를 사용하지 않은 이유
  - `NodePort` 를 이용하면 GCP 를 이용할때 다른 node 에 port 로 접근하는 경우에 방화벽 문제가 생겨서 clusterIP 로 접근함
  - 제일 처음 진입점인 `serve-svc` 는 `LoadBalancer` 로 정의했고, GCP 의 LoadBalancer 와 연동되어 `ExternalIP` 를 부여받아 쉽게 테스트 하였음

- `service dns` 사용 
  - `<service-name>.<namespace>.svc.cluster.local` 형태로 dns 가 생성됨 
  - e.g. 
```shell script
gRPC: grpc-svc.default.svc.cluster.local
HTTP: http://http-svc.default.svc.cluster.local/hellohttp
```
  - gRPC port 설정을 헤맸는데, `hellogrpc` 는 8000 을 listen, `grpc-svc` 도 동일한 port(8000) 값을 주고, dns 도 8000 port 호출하여 해결함 
```yaml
ports:
  - protocol: TCP
    port: 8000
``` 
          
- gRPC round robin 설정
  - gRPC 는 HTTP2 기반으로 동작하고, HTTP1.1 과는 달리 한 connection 에 여러 request/response 를 재사용함 
  - 결국 connection layer 인 layer-3 에서 동작하는 LoadBalancer 가 아닌, application layer 인 layer-7 에서 동작하는 LoadBalancer 이용이 필요함 
  - 따라서 `Headless service` 를 정의(ClusterIP 값은 None)하고, 아래와같이 dns 세팅을 하고 `"round_robin"` 설정함
   
1) `grpc.Dial` 에 `dns:///` (/ 3개 맞음) 를 명시함
```go
conn, err := grpc.Dial("dns:///grpc-svc:8000",
        grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancerName("round_robin"))
```  
2) 또는, `resolver.SetDefaultScheme("dns")` 를 명시
```go
resolver.SetDefaultScheme("dns")
conn, err := grpc.Dial("grpc-svc:8000",
        grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancerName("round_robin"))
```

- gRPC compile, code generation 은 아래 script 로 대체함
```proto
$ cd hellogrpc 
$ protoc -I hellogrpc/ proto/hello.proto --go_out=plugins=grpc:hellogrpc
``` 

- 각 directory 에 `vendor` 는 왜만들었는지?
  - `docker build` 할때 `go build` 하며 dependency pkg 다운받는 시간이 오래걸려서 바꿔버림.
  - dockerfile 단순화를 위해 multi-stage build 였던 부분을 제거함  

## how-to run
### docker build
- 각 directory 내부의 `dockerbuild.sh` 실행 후 docker hub 에 올려서 사용 (개인 docker hub 이용)

### kubectl 을 이용한 scheduling
```shell script
$ kubectl apply -f kubeyamls
```

### yaml 파일들 다시 적용시
```shell script
$ ./restart.sh
``` 

## demo
> curl 34.64.214.65/hellogrpc 반복

> pod 가 많으면 더 잘보이겠지만 일단은 round-robin 되는것으로 추정됨  
```json
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-q2b4b","response":{"version":"v1-hellogrpc","hostname":"hellogrpc-6dd9d78b49-j66hq"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-q2b4b","response":{"version":"v1-hellogrpc","hostname":"hellogrpc-6dd9d78b49-j66hq"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-png98","response":{"version":"v1-hellogrpc","hostname":"hellogrpc-6dd9d78b49-z2zm9"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-q2b4b","response":{"version":"v1-hellogrpc","hostname":"hellogrpc-6dd9d78b49-j66hq"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-q2b4b","response":{"version":"v1-hellogrpc","hostname":"hellogrpc-6dd9d78b49-z2zm9"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-png98","response":{"version":"v1-hellogrpc","hostname":"hellogrpc-6dd9d78b49-z2zm9"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-q2b4b","response":{"version":"v1-hellogrpc","hostname":"hellogrpc-6dd9d78b49-j66hq"}}
```
> curl 34.64.214.65/hellohttp 반복

> pod 가 많으면 더 잘보이겠지만 일단은 round-robin 되는것으로 추정됨
```json
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-vtwq8","response":{"hostname":"hellohttp-7859b685bb-c7cr7","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-q2b4b","response":{"hostname":"hellohttp-7859b685bb-ggd7b","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-png98","response":{"hostname":"hellohttp-7859b685bb-c7cr7","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-png98","response":{"hostname":"hellohttp-7859b685bb-c7cr7","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-vtwq8","response":{"hostname":"hellohttp-7859b685bb-c7cr7","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-vtwq8","response":{"hostname":"hellohttp-7859b685bb-c7cr7","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-q2b4b","response":{"hostname":"hellohttp-7859b685bb-ggd7b","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-png98","response":{"hostname":"hellohttp-7859b685bb-c7cr7","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-q2b4b","response":{"hostname":"hellohttp-7859b685bb-ggd7b","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-vtwq8","response":{"hostname":"hellohttp-7859b685bb-c7cr7","name":"","version":"v1-hellohttp"}}
{"name":"","version":"v1-helloserve","hostname":"helloserve-84ff4df84d-vtwq8","response":{"hostname":"hellohttp-7859b685bb-c7cr7","name":"","version":"v1-hellohttp"}}
```
 
## flowchart
| | | | | | | | | |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
|`client`| -> | `serve-svc` | -> | `hello-serve` | -> | ( /hellohttp ) | -> | `hellohttp` |
| | | | | | -> | ( /hellogrpc ) | -> | `hellogrpc` |

![image](https://user-images.githubusercontent.com/34496756/76628501-29385b80-6580-11ea-947c-ec154e7db570.png)

## system architecture

### deployment
#### hellogrpc
- gRPC server, listen `8000/tcp`

#### hellohttp
- HTTP server, listen `8100/tcp`

#### helloserve
- API gateway, listen `8200/tcp`
- api

| path | description | response |  
| --- | --- | --- |
| /hellohttp | hellohttp 호출 | |
| /hellogrpc | hellogrpc 호출 | |


```shell script
$ kubectl get svc
NAME         TYPE           CLUSTER-IP    EXTERNAL-IP    PORT(S)        AGE
grpc-svc     ClusterIP      None          <none>         8000/TCP       5s
http-svc     ClusterIP      10.84.4.182   <none>         80/TCP         12m
kubernetes   ClusterIP      10.84.0.1     <none>         443/TCP        141m
serve-svc    LoadBalancer   10.84.1.232   34.64.214.65   80:31615/TCP   139m
```

```shell script
$ kubectl get deploy
NAME         READY   UP-TO-DATE   AVAILABLE   AGE
hellogrpc    3/3     3            3           12m
hellohttp    3/3     3            3           12m
helloserve   3/3     3            3           12m
```

### service
#### serve-svc
- type: `LoadBalancer`
- targetPort: 8200
- helloserve 의 loadbalancer 로 동작함 (round-robin) 
   
#### grpc-svc
- type: `ClusterIP (Headless Service)`
- port: 8000
- hellogrpc 의 loadbalancer 로 동작
- 다만 headless service 이기 때문에 helloserve 에서 client loadbalancer 구현해서 호출함

#### http-svc
- type: `ClusterIP`
- targetPort: 8100
- hellohttp 의 loadbalancer 로 동작함 (round-robin)


## reference
- https://medium.com/@ammar.daniel/grpc-client-side-load-balancing-in-go-cd2378b69242
- https://corgipan.tistory.com/9
- https://www.marwan.io/blog/grpc-dns-load-balancing
- https://blog.nobugware.com/post/2019/kubernetes_mesh_network_load_balancing_grpc_services/
