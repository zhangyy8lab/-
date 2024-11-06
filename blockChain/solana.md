

## Solana 节点

1. 验证者节点（Validator Nodes）

	• 负责处理交易，参与共识，维护网络的安全性和稳定性。

2. 归档节点（Archiver Nodes）

	• 归档节点存储整个区块链的历史数据，以便于后续的数据查询和分析。它们通常会保存完整的交易历史和状态数据。

3. RPC 节点（RPC Nodes）

	• RPC（Remote Procedure Call）节点提供接口供外部应用程序（如钱包、DApp）与 Solana 网络进行交互。RPC 节点通常会接收和处理用户请求，将请求转发给验证者节点进行处理。

4. 领导者节点（Leader Nodes）

	• 在每个区块时间段内，验证者节点将会被选为领导者，负责创建新区块并将其添加到链上。这些节点在执行交易时扮演关键角色。

5. 客户端节点（Client Nodes）

	• 客户端节点是指使用 Solana 网络的各种客户端应用，如钱包和去中心化应用（DApp），这些节点通常不直接参与共识，而是依赖 RPC 节点来与网络交互。

6. 监视节点（Monitoring Nodes）

	• 这些节点用于监控和记录网络的状态和性能，可以帮助开发者和操作人员检测网络的健康状况。



### Install

#### install rust

```bash 
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
rustc --version
```

#### install Solana Client

```bash
sh -c "$(curl -sSfL https://release.anza.xyz/stable/install)"
export PATH="$HOME/.local/share/solana/install/active_release/bin:$PATH"
solana --version
```

#### cargo build

```bash
# apt package
apt install libudev-dev

# 
cd ~ && git clone https://github.com/solana-labs/solana.git
cd solana && cargo build  --release

# build after path ~/solana/target/release/
```



### use solana

#### config

```bash
# get setting
solana config get 

# config info
Config File: /root/.config/solana/cli/config.yml
RPC URL: https://api.mainnet-beta.solana.com
WebSocket URL: wss://api.mainnet-beta.solana.com/ (computed)
Keypair Path: /root/.config/solana/id.json
Commitment: confirmed

# setting 
solana config set -um    # For mainnet-beta https://api.mainnet-beta.solana.com
solana config set -ud    # For devnet https://api.devnet.solana.com 
solana config set -ul    # For localhost http://localhost:8899
solana config set -ut    # For testnet https://api.testnet.solana.com

# setting local
solana config export-address-labels 47.76.174.179
solana config get
Config File: /root/.config/solana/cli/config.yml
RPC URL: http://localhost:8899
WebSocket URL: ws://localhost:8900/ (computed)
Keypair Path: /root/.config/solana/id.json
Commitment: confirmed

# 
8899 TCP - 通过 HTTP 进行 JSONRPC。使用 `--rpc-port RPC_PORT`` 进行更改
8900 TCP - 通过 Websockets 进行 JSONRPC。派生。用途RPC_PORT + 1
```



#### keypair

##### id.json

```bash
mkdir ~./config/solana

# solana-keygen not found & cp ~/solana/target/release/solana-keygen /usr/local/bin/ 
solana-keygen new -o ~/.config/solana/id.json 
....
Wrote new keypair to /root/.config/solana/id.json
=============================================================================
pubkey: 2qhnjHGdfZgLbbfaaVJSzsec8iJikFb3bVEqPj5pm93C
=============================================================================
Save this seed phrase and your BIP39 passphrase to recover your new keypair:
private before favorite future wood trade mad remove over source review short
=============================================================================

# id.json defualt 500000000 SOL 
```

##### validator.json 

```bash
solana-keygen new -o ~/.config/solana/validator.json 
....
Wrote new keypair to /root/.config/solana/validator-keypair.json
========================================================================
pubkey: 7TnUFn573MnKwN782BchuWB6ykesvyhNd9fzw7goEERG
========================================================================
Save this seed phrase and your BIP39 passphrase to recover your new keypair:
jazz cross panel lesson adapt obey swap dignity toss special butter clog
========================================================================
```

##### user1.json

```bash
solana-keygen new -o ~/.config/solana/user1.json 
...
Wrote new keypair to /root/.config/solana/user1.json
=========================================================================
pubkey: 9883VJ3WJfdraH5gJMxroXDgYvUBPJNdV5uiD1qBGWVr
=========================================================================
Save this seed phrase and your BIP39 passphrase to recover your new keypair:
member pitch supreme injury salon churn state urge hard bonus once tongue
=========================================================================
```





#### start validaor

```bash
solana-test-validator -C ~/.config/solana/cli/config.yml 
```

#### params 

```bash
--gossip-host：指定其他节点连接此节点时使用的外部地址（常用于公网 IP）。
--bind-address：指定节点监听的本地接口地址（常用于限制本地接口绑定）。
```





#### airDrop get SOL balance

```bash 
solana airdrop 2 <user1.json>

# get balance 
solana balance <user1.json>
2 SOL
```

#### 转账

```bash
# 转账
# solana transfer --from .config/solana/id.json <user1.json> 500  --allow-unfunded-recipient
# --allow-unfunded-recipient 参数允许系统跳过账户初始化检查，将 SOL 转入尚未资助的地址。 
solana transfer --from .config/solana/id.json  9883VJ3WJfdraH5gJMxroXDgYvUBPJNdV5uiD1qBGWVr 500

Signature: 5DsKvUKub79Wqtvo6cbRUC2AYvZi9XvQNwsUzaSf56gcQhtTiGgXa4XXcawCUAsFafoyRiNX5snQfPh2pFL9bafL


# 查询转账信息 
solana balance <user1.json> 
502 SOL

# 查看 confirm 状态
# solana confirm <transaction_signature> 
solana confirm 5DsKvUKub79Wqtvo6cbRUC2AYvZi9XvQNwsUzaSf56gcQhtTiGgXa4XXcawCUAsFafoyRiNX5snQfPh2pFL9bafL
Finalized

# 账户的交易列表
# solana transaction-history <user1.json> 
solana transaction-history  9883VJ3WJfdraH5gJMxroXDgYvUBPJNdV5uiD1qBGWVr
5DsKvUKub79Wqtvo6cbRUC2AYvZi9XvQNwsUzaSf56gcQhtTiGgXa4XXcawCUAsFafoyRiNX5snQfPh2pFL9bafL
2aFNvxcUBQQE28crHiApvUNme64oALo4NoQF9QXEXCfjZTsp3nreNWTiXxbfSGkUsie48sNrXhumjMZMGW6ZMYak
2 transactions found

```



#### 性能测试

```bash
# transfer
for i in {1..1000}; do
  solana transfer --from ~/.config9883VJ3WJfdraH5gJMxroXDgYvUBPJNdV5uiD1qBGWVr 1 &
done


# 
for i in {1..1000}; do
  solana airdrop 1 --url http://127.0.0.1:8899
done

```





### solana-explorer

```bash
git clone https://github.com/solana-labs/explorer.git
cd explorer 
# npm install --legacy-peer-deps # 忽略一些依赖冲突
pnpm i 

# add env
vim .env 
REACT_APP_RPC_URL=http://127.0.0.1:8899

# start
pnpm dev
```



#### nginx-proxy

```bash
```





## Solana总体架构

### 账号

> • 系统账户：用于用户资产管理，支付交易费用。
> • 程序账户：用于存储和执行智能合约代码，实现去中心化逻辑。
> • 数据账户：用于存储合约的状态数据，便于合约持久化信息。

#### 系统账户(SystemAccount)

记录该账户持有的 SOL 数量，可用于支付交易费用或进行转账。及交易的手续费，也用来创建其他账户或分配内存给新账户。



#### 程序账户(ProgramAccount)

用于存储和运行智能合约代码的账户,

程序账户本身不持有 SOL 余额，它的核心用途是部署和执行智能合约逻辑.一旦部署后，程序账户是不可修改的



#### 数据账户(DataAccount)

> 也称 PDA - Program Derived Address

在去中心化应用中，数据账户用来存储用户的余额、配置、NFT 资产元数据等状态信息，以便在合约交互中读写和修改。

数据账户通常使用 PDA（Program Derived Address）生成，这些地址没有私钥，也无法直接通过普通签名控制, 存储程序合约的状态数据，通常与程序账户关联, 程序账户可以设置特定的权限，仅允许相关的程序对数据账户进行操作，防止数据泄露或恶意篡改。



### 组件

#### Leader， PoH Generator

1. Leader 是被选举出来的PoH Generator
2. 它接收用户的交易， 并输出所有交易的PoH序列，该序列保障了Solana系统中的全球一致的顺序
3. 针对一批次的交易， 该Leader会针对交易的顺序运行的结果而产生的状态， 进行签名并发布
4. 该签名使用的Leader的私钥签的



#### 状态

1. solana 系统状态由Hash表来维护，而且该表是基于用户地址来索引的
2. 表中每个条目饮食了用户完整的地址， 以及计算要用到的信息
3. 下面是两个状态表的例子

- 交易表

![image-20241023164414475](/Users/zhangyy/Library/Application Support/typora-user-images/image-20241023164414475.png)

共占用了32个字节

- PoS质押表

  ![image-20241023164547545](/Users/zhangyy/Library/Application Support/typora-user-images/image-20241023164547545.png)

  共占用64个字节

  #### verifier 状态复制

  1. verifier节点用来复制Solana链的状态， 并确保该链状态的高可用性
  2. 节点要复制的内容这或目标(target)， 是由共识算法决定的； 共识算法中的Validators会基于链下定义好的准则，选择并通过投票来确定PoRep节点
  3. Solana网络配置了最少的PoS质押质押数额，并且一个复制身份(replicator identity),只能有一个质押账户

  #### Validators

  1. 这些节点是虚拟节点， 跑在Verifiers或者Leader所在的机器上， 或者独立的机器上
  2. 它们专门用来执行Solana配置的共识算法
  3. 当机器扮演Leader角色时， 这些Validators是不运行的； 也就是扮演Verifier的角色时才会运行

### 网络限制(Network Limits)

1. Leader 期望

- 接收所有的用户发来的请求

- 将用户请求的包以最高效率排好序

- 将排好序的请求， 编排进PoH时间流逝包(package)；

- 将编排好的时间流逝发布给下游的Verifiers

2. 因此在PoH Generator这里会产生网络瓶颈

> 会出现一个多对一， 一对多的情况 多进1， 1还要出多

![image-20241023170704388](/Users/zhangyy/Library/Application Support/typora-user-images/image-20241023170704388.png)

3. 要想提高PoH的网络吞吐能力，交易的内存访问方式很关键； 因此交易都按序排放， 这样可以让错误而降到最低， 而让预测先取的交易， 可以做到最大
4. 下面是不同数据包的存放格式

- 输入包格式

![image-20241023171904330](/Users/zhangyy/Library/Application Support/typora-user-images/image-20241023171904330.png)

> Last Valid Hash 交易算出后的hash
>
> Counter： 计数器
>
> Signature 1/2 是Leader进行的签名

共占用： 20+8+16+8+32+32+32=148 bytes

- 上面包格式里， 最小的payload是一个目标账号

  > 这个是包含一个转账号信息， 

![image-20241023172320802](/Users/zhangyy/Library/Application Support/typora-user-images/image-20241023172320802.png)

存在payload时的包占用空间为 176 bytes

5. 对于PoH 时间流逝包里面包含了以下内容：

- 当前Hash
- counter
- 在该时间流逝中所有新消息（Message）的hash
- 所有消息被处理后的状态签名

![image-20241023173131460](/Users/zhangyy/Library/Application Support/typora-user-images/image-20241023173131460.png)

此类包的最小占用存储空间： 132 bytes

这个时间流逝包每隔N个消息会， 广播一次

6. 对于1gbps的网络连接

- 最大的TPS = 1gbps / 176 bytes = 710k