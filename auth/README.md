# 认证



## 测试数据库

192.68.136.27
new_chain_db

root
123456


## 节点认证地址
```shell
curl -X POST  -H "Content-type: application/json" --data '{"key":13434954,"nodeId":"4"}' http://192.168.136.39:9094/newchain/rod-api/node/getIsStart
```


## 修改节点状态
```shell
curl -X POST  -H "Content-type: application/json" --data '{"key":13434954,"status":"4"}' http://192.168.136.39:9094/newchain/rod-api/node/updateStatus
```


## 测试命令
```shell
geth_build_init_run --auth.server=http://192.168.136.39:9094/newchain --auth.code=123456


 geth_run --auth.server=http://192.168.136.39:9094/newchain --auth.code=320389200717164634
```