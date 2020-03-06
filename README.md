# evrynet-tools 

`evrynet-tools` contains tools for development & testing evrynet network.

## Build accounts command line interface
* Go 1.12+
```shell script
$ make accounts
$ ./build/accounts -help
NAME:
   accounts - The accounts command line interface

USAGE:
   accounts [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   generate  generate a number accounts based on a seed
   deposit   Deposit EVR to the generated accounts
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

To generate accounts you can use this command  
`./build/accounts generate --num 5 --seed testnet`

To deposit to accounts you can use this command  
`./build/accounts deposit --num 10 --seed testnet --expectedbalance "1000000000000000000" --senderpk "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1" --rpcendpoint "http://0.0.0.0:22001"`

<details>
<summary>Output</summary> 

```json 
{
	"private_key": "bd35ed6ecf65de973d82d81692075e24dd1c432f780cee3ab34cef5a56e1d751",
	"public_key": "043e9039812f828d3086d1f5383be5d0125c7a40049c2ed9aa02affa13ce897548902773446822333551bb31b07344a5212e6cdb4f7ca6fe6a73b92914dfb5bcb1",
	"address": "0x879B0b268dbA7668678FeFe283a9995FB5f8cBeb"
}
{
	"private_key": "deb1ff1f17ece293c576d5a0c1202af4fee9280791c0baa1d2e4e8659847f646",
	"public_key": "04fb49ad4df6cbf272f03f40ab722b00be9db48af075a8e957674e7402aa6c4fe531f665747155d035debedf453b04167049b2a6c2b1b1b3ea2bb44aec3ceaebc1",
	"address": "0xF44B353c9d3bAcdd1B22898a4b14372bC85a40cB"
}
{
	"private_key": "bc9d6000f18f5963c810515ed5b90dc1c2f11ce9f4027e82b08b6725daff404b",
	"public_key": "04678ab7ac69e9ea5bf967119977e9175ca00c12b13c20d4a49da940ea7e7839db1be998d8120ac2bc85d3019ec2d03fdadc39a3da88e1e66728061fb4f6e6ad8a",
	"address": "0x65fE8cc4E7ce281Afb5dC0B875DaB983D57522BD"
}
{
	"private_key": "db676ee7ff9cff6ed067d18e8e754ff3be955a5bba695ccde7d5c24645681251",
	"public_key": "0413f6148b74b15c9d14a6c0851643e9da948027e2fc39971c669cbde506618da8503050cc283c3ab0191aad10328b97c91710b80a02db81c7b77583cccbad5517",
	"address": "0xAE2c412B2651d3aABce6F2F67Ab079f5B06a2ADd"
}
{
	"private_key": "8d8546977f0f85f0ffd1399a813793c7f4a1d80ec66b9f66f5c09c6c46be86d5",
	"public_key": "04d097709ee34bf0c857eedb6599de9e3d1b0aaee7b5b6332c3faee5115ddf677f5e919ca602966211c939cad329d6aa123269f4af84c4257cb78b4d1b551d27ba",
	"address": "0x844e6d9b98c88924a042514d218c415406cE1846"
}
```
</details>

## Build transactions command line interface  
```shell script
$ make tx_flood
$ ./build/tx_flood -h

NAME:
   tx_flood - The tx_flood command line interface

USAGE:
   tx_flood [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --num value             Number of accounts want to generate (default: 4)
   --seed value            Seed to generate private key account (default: "evrynet")
   --num-tx-per-acc value  Number of transactions want to use for an account (default: 1)
   --flood-mode value      Flood mode when send Tx: 0: Random, 1: Normal Tx, 2: Tx with SC (default: 0)
   --continous             Flood continously if set to true
   --sleep-duration value  Time to sleep after each batch of numAccount*numTxPerAcc flooding (default: 1s)
   --rpcendpoint value     RPC endpoint to send request (default: "http://0.0.0.0:22001")
   --help, -h              show help
   --version, -v           print the version
```  
To use tx flood you can use this command  
`./build/tx_flood --num 3 --num-tx-per-acc 2 --seed testnet --rpcendpoint "http://0.0.0.0:22001" --flood-mode 2`

## Build transactions metric command line interface  
```shell script
$ make tx_metric
$ ./build/tx_metric -h
NAME:
   tx_metric - The tx_metric command line interface

USAGE:
   tx_metric [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --rpcendpoint value  RPC endpoint to send request (default: "http://0.0.0.0:22001")
   --start-block value  Where blocknumber start at (default: 0)
   --duration value     Duration to calculate metric (default: 1m0s)
   --help, -h           show help
   --version, -v        print the version
```  
To use tx metric you can use this command  
`./build/tx_metric --rpcendpoint "http://0.0.0.0:22001" --start-block 1681 --duration 60s`

## Build block monitor command line interface  
```shell script
$ make blockmonitor
$ ./build/blockmonitor -h
NAME:
   blcMonitor - sends messages to telegram when node dont increase blocks

USAGE:
   blockmonitor [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   start    Alert to telegram when block is stuck
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

$ ./build/blockmonitor start -h
NAME:
   blockmonitor start - Alert to telegram when block is stuck

USAGE:
   blockmonitor start [command options] [arguments...]

DESCRIPTION:
   Alert to telegram when block is stuck

OPTIONS:
   --apiToken value     The API token (default: "935265043:AAG02pbFdpOhT31ALs1tzlLTl_NVnxYuTF4")
   --chatId value       The ID of group/chanel (default: -375817595)
   --rpcendpoint value  RPC endpoint to send request (default: "http://0.0.0.0:22001")

```  

## Build staking command line interface  
```shell script
$ make staking
$ ./build/staking -h
NAME:
   stakingClient - sends a vote/ unvote/ register/ resign for a candidate to a node

USAGE:
   staking [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   vote      Sends a vote for a candidate
   unVote    Sends a un-vote for a candidate
   resign    resign a candidate, only called by owner of that candidate
   register  register a new candidate
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

```  

To use staking tool you can use this command  

```
 ./build/staking register --stakingsc 0x2d5bd25efa0ab97aaca4e888c5fbcb4866904e46 --senderpk 85af6fd1be0b4314fc00e8da30091541fff1a6a7159032ba9639fea4449e86cc --candidate 0x71562b71999873DB5b286dF957af199Ec94617F7

 ./build/staking vote --stakingsc 0x2d5bd25efa0ab97aaca4e888c5fbcb4866904e46 --senderpk 85af6fd1be0b4314fc00e8da30091541fff1a6a7159032ba9639fea4449e86cc --candidate 0x71562b71999873DB5b286dF957af199Ec94617F7 --amount 10

 ./build/staking unvote --stakingsc 0x2d5bd25efa0ab97aaca4e888c5fbcb4866904e46 --senderpk 85af6fd1be0b4314fc00e8da30091541fff1a6a7159032ba9639fea4449e86cc --candidate 0x71562b71999873DB5b286dF957af199Ec94617F7 --amount 1

 ./build/staking resign --stakingsc 0x2d5bd25efa0ab97aaca4e888c5fbcb4866904e46 --senderpk 85af6fd1be0b4314fc00e8da30091541fff1a6a7159032ba9639fea4449e86cc --candidate 0x71562b71999873DB5b286dF957af199Ec94617F7 

 ```
