# Blockchain-based Voting System

## How-to

```shell
git clone github.com/carterbloop/goblockchain
cd goblockchain
go run main.go -data-dir tmp -server-port 8080 -entry-ip localhost -entry-port 8081
go run main.go -data-dir tmp2 -server-port 8081 -entry-ip localhost -entry-port 8080
go run main.go -data-dir tmp3 -server-port 8082 -entry-ip localhost -entry-port 8081
```
