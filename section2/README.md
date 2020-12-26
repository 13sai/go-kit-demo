# 服务发现与注册

> 结合consul和go-kit，并运用k8s部署

## 构建

```
docker build -t  register .

kubectl apply -f consul-server-service.yaml
kubectl apply -f consul-server.yaml
kubectl apply -f consul-server-http.yaml
kubectl apply -f consul-client.yaml


kubectl apply -f register.yaml
```

![consul](./imgs/k8s-go-kit13.jpg)
![k8s](./imgs/k8s-dashbord13.jpg)

## 概念

> 服务注册：服务示例在启动的时候将自身信息注册到服务注册于发现中心，并在运行时通过心跳等方式向服务注册与发现中心上报自身服务状态

> 服务发现：服务示例根据服务名向服务注册和发现中心请求其他服务示例信息，用于进行接下来的远程调用