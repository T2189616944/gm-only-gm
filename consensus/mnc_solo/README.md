# MNC_solo
单机排序

1, 排序节点手动指定。启动参数。
2. 非排序节点, 需要将块发送到排序节点
3. 排序节点丢失， 非排序节点不会自动升级为排序节点。出块停止。
4. 同一高度块， 选择先到的为准。
5. 使用排序的密钥签名， 签名验证通过就认可结果。
6. 签名结果放在extradata 字段

--verbosity=9 --debug

./geth --datadir=solo --http --http.api admin,debug,web3,eth,txpool,miner,net,solo --http.addr 127.0.0.1 --solo --solo.key=16e96c512590dcf3cff5cb0c742e913271b9f50961750ec07de4b21db3d492bc  --solo.main  --port=10240  console


./geth --datadir=solo2  --solo   --solo.main.addr=http://127.0.0.1:8545  --port=10241 console



