# MNC_raft




./geth --datadir=raft --http --http.api admin,debug,web3,eth,txpool,miner,net,raft --http.addr 127.0.0.1  --port=10240  --raft console


./geth --datadir=solo2  --solo   --solo.main.addr=http://127.0.0.1:8545  --port=10241 console