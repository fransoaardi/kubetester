kubectl delete deploy --all
kubectl delete svc/grpc-svc svc/http-svc

kubectl apply -f kubeyamls
