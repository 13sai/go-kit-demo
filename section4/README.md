
```
自定义network

docker network create kong-net

启动

docker-compose up

kong启动应该会exit的，因为KONG_DATABASE需要初始化

迁移
docker run --rm --network=kong-net -e "KONG_DATABASE=postgres" -e "KONG_PG_HOST=kong-database" -e "KONG_PG_PASSWORD=123456" kong:latest kong migrations bootstrap

docker ps -a

查看kong依旧exit了

docker logs 查看kong的失败日志，按照日志执行
比如：
2021/01/09 12:40:53 [error] 1#0: init_by_lua error: /usr/local/share/lua/5.1/kong/cmd/utils/migrations.lua:20: New migrations available; run 'kong migrations up' to proceed

那么执行一下 kong migrations up

docker run --rm --network=kong-net -e "KONG_DATABASE=postgres" -e "KONG_PG_HOST=kong-database" -e "KONG_PG_PASSWORD=123456" kong:latest kong migrations up


docker ps -a
可以看到kong也成功运行起来了

```

遇到的问题汇总：

```
add router 错误
konga.  3 schema violations (headers: unknown field; https_redirect_status_code: unknown field; path_handling: unknown field)x-some-header:foo,bar


升级kong到2.x
```

```
curl -i -X POST --url http://localhost:8001/services/ --data 'name=demo1' --data 'url=http://github.13sai.com/'


curl -i -X POST --url http://localhost:8001/services/demo1/routes --data 'hosts[]=sai0556.com'
```
![konga](https://imgconvert.csdnimg.cn/aHR0cHM6Ly9zZWdtZW50ZmF1bHQuY29tL2ltZy9iVmJ4RXhV?x-oss-process=image/format,png)
```
http://kong:8001
```

services
````
Protocol https
Host github.13sai.com
Port 443 
```

routes
```
Hosts sai.com
Paths /test
```

```
curl http://localhost:8000/test/2020/11/21/283/ --header 'Host: sai.com'
```