# go-kit-demo
go kit demo

```
docker pull mysql:5.7
docker pull alpine:3.12
docker pull redis:5.0

---

docker run -itd --name redis-5.0 -p 6389:6379 redis:5.0

docker run  -itd --name mysql-for-user -p 3316:3306 -e MYSQL_ROOT_PASSWORD=111111 mysql-for-user 

docker build -t user-alpine .

docker run -itd --name user --network host user-alpine
```