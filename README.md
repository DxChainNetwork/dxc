# DxChain 3.0

## The Ecosystem Powered by DxChain 3.0

### Smart Contract Platform

While compatible with all smart contract functions on Ethereum (completely migrated to the DxChain Mainnet), we try to provide worldwide users with easy-to-use and templated development solutions to help popularize smart contracts and provide more diversified dapps.

### Cross Chain

Under the premise of ensuring the absolute security of asset circulation, cross-chain circulation of assets with low handling fees, low latency and high concurrency can be completed through DxChain 3.0 on-chain transfer without any centralized platform. To achieve the circulation and exchange of tokens by connecting mainstream projects and building a global asset interactive network.

### DeFi

We are committed to building a more convenient, friendly, and safer decentralized financial service platform for worldwide users to transfer crypto assets easily and transparently. In the meantime, it will also provide our ecosystem partners with appropriate and feasible financial products and services, including but not limited to mainstream DeFi applications such as DEX, loans, and liquidity mining to meet the diverse financial needs.

### NFT
Currently, the Non-Fungible Token (NFT) can be used in crypto-collectibles, games and other applications. DxChain will help creators, developers, and collectors to perform NFT minting, development, and trading more conveniently and stably, providing the industry with integrative solutions.

### Metaverse
As a bridge between the real world and the virtual world, Metaverse will further affect social, entertainment, finance and other aspects in the future. In this regard, DxChain will join this market in advance, conduct a forward-looking exploration regarding cross-chain identification, social entertainment, integration of crypto and tangible assets, GameFi, etc.

## DxChain 3.0 Architecture

### Fully Compatible with EVM

DxChain 3.0 is fully compatible with Ethereum Virtual Machine (EVM), supports the compilation and execution of smart contracts, and supports various versions of Solidity. The Ethereum Virtual Machine is a Turing-complete state machine, an engine used to execute transactions or contract code. EVM provides a secure operating environment for each contract with an independent runtime stack, which contains a maximum of 1024 elements, and each element is 256bit. Moreover, EVM supports cyclic operation instructions and the contract supports complex logic functions, enabling any complex conceived programs to run smoothly.

### DPoS Consensus

The DxChain mainnet currently implements the DPoS (Delegated Proof of Stake) algorithm, which is considered an improved version of the PoS (Proof of Stake) algorithm and has the characteristics of democratization, low costs, low latency, and high concurrency. DxChain 3.0 improves the DPoS consensus algorithm by storing the staking and voting data in the built-in contract and keeps the data storage structure consistent with the state trie, allowing more efficient elections and block synchronization. 

* TPS: 500+
* Average block interval: 3s
* Validators: <=99
* Delegators: <=210

The system contracts code here: `consensus/dpos/systemcontract/contracts`

<img src="https://dxchain.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2Ft7Yq0wZVG1pJ3PiTzoaE%2Fuploads%2F2cFTPWdoV8zaNjS7uqCN%2F%E8%8B%B1%E6%96%878.png?alt=media&token=a4b0aaf6-77a1-4734-93f1-c89605cb8b17" style="width:75%;" />


### Cross Chain
DxChain Bridge uses AWS Nitro Enclave to build a fast, safe, and low-cost cross-chain bridge between Ethereum and DxChain. The DxBridge will be composed of Nitro Enclave and a list of trusted nodes (called Warden). Nitro Enclave will be used to build an isolated execution environment to prevent any centralized interference and reduce the complex steps without sacrificing system security.

DxBridge mainly consists of two parts:
- Nitro Enclave: AWS Trusted Execution Environment solution. By creating an isolated environment, users can use and process private keys with high security while preventing users and applications on the parent instance to views or obtaining those information.
- A set of Wardens: third-party searchers and verifiers of transactions. Warden is mainly responsible for retrieving DxChain and Ethereum and submitting legal transactions that need to be processed to Nitro Enclave. First, Warden will look for transactions that have been successfully sent to the Ethereum wallet or transactions that have been retrieved from DxChain. There must be enough funds to pay for related expenses, including the gas fee and cross-chain fees required, otherwise, the transaction will be rejected and Warden will not retrieve these transactions. Nitro Enclave requires a certain number of Wardens to submit the same transaction at the same time, then the bridge will send the corresponding transaction on another chain and submit legal transactions by providing a private key segment. 

<img src="https://dxchain.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2Ft7Yq0wZVG1pJ3PiTzoaE%2Fuploads%2FUAProUdP9pqJI5Ft3hil%2F%E8%8B%B1%E6%96%871.png?alt=media&token=52f10aa6-a1a6-49f4-bc11-46449a4229ea" style="width:75%;" />


Nitro Enclave can directly connect with Warden to obtain on-chain events and send transactions. The private keys of all addresses in the transaction are derived from the master private key generated during initialization that no other party can obtain. The master private key uses the Shamir Secret Sharing algorithm to distribute the private key segments to Warden, and uses TLS communication to verify the identity during the process. Nitro Enclave will ask Warden for private key segments via TLS connection to retrieve the master private key, and distribute new private key segments to Warden again after restart. In addition, transactions confirmed to have been processed by the bridge will be backed up locally.

## Tokenomic

The native token on the DxChain is DX, an exchange unit between users, provides economic incentives for on-chain activities and is also consumed as handling fees. The total supply of DX is 100 billion, of which the ecosystem pool is 45 billion. The undistributed rewards for DxChain 2.0 will continue to be distributed in DxChain 3.0.

<img src="https://dxchain.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2Ft7Yq0wZVG1pJ3PiTzoaE%2Fuploads%2FeMid4WM6kk77VdHcz817%2F%E8%8B%B1%E6%96%872.png?alt=media&token=ebb15a4b-b9aa-478b-b54d-40e0ebb00503" style="width:75%;" />

Users can participate in DxChain by voting DX to obtain block rewards or become a Validator to maintain the security and stability of the DxChain network. In DxChain 3.0, there is also a new token **burn mechanism**. If a Validator does not meet the required number of block generation in a certain epoch, all block rewards that have been generated will be burned and the Validator status will be canceled.

## Governance

DxChain mainly divides business, technology, and community aspects in community governance. It is expected that major decisions and policies will be governed by community voting. The community initiates proposals, evaluates the results, and the committee monitors the execution process so that the DxChain development team and the community can realize the co-governance. 

The DxChain team always firmly believes that a fair, reasonable, and transparent governance mechanism with multi-party participation can better improve the community's quality. DxChain will adhere to the following points:
* Improve the community incentive mechanism
    We will continue to run DPoS mining, keep the incentive mechanism updated and innovated, provide sufficient incentives to attract more users and achieve positive feedback, aiming to guarantee the benefits of all users.

* Community Co-governance
    Co-governance will be the core idea of DxChain 3.0 community governance. In order to increase the sense of ownership, community users will be advocated and guided to participate in community discussions, proposals, and voting in the design of the community governance mechanism.

* Foundation Assistance
    As the initiator of DxChain ecosystem, DX Foundation will play the role of mobilizing community participation, accelerating community merging, integrating community resources, solving community problems, and promoting community governance.

### Election

The validator is responsible for block generation and verification on the chain, and is an integral part of DxChain. In order to realize co-governance, the rules were firstly applied in node elections. If a user wants to participate in the DPoS and apply for a validator, the following conditions must be met:
- Possess the technical and hardware requirements to maintain a blockchain node
- Stake at least 40 million DX
- The proportion allocated to voters is between 70% and 100%
A proposal will be initiated on DxChain once the above conditions are met and all existing delegators can vote, the application will not get approved until at least one validator has voted. 
### DAO
DAO (Decentralized Autonomous Organization) involves on-chain governance. Project or community users can initiate proposals, such as adjustments to consensus algorithm and economic models, deciding whether to approve new delegators, greatly expressing our idea of co-governance and reducing the risks brought by centralized governance.
### Foundation

Currently, 5% of the block reward will be allocated to the DxChain Foundation to support future development and operation, community governance, external developer contribution rewards, ecosystem construction funds, etc. (including but not limited to Dapp Development and marketing). Community governance is a process in which delegators, community members and foundations supervise each other and work together. The DxChain team will continue to improve the further governance plan in DxChain 3.0 to make governance and supervision more transparent.

## Follow Us

Website：https://dxchain.com/

Telegram (EN): https://t.me/dxchain

Telegram (CN): https://t.me/DxChainGroup_CN

Discord: https://discord.gg/XbPwmErhDX

DX Explorer: https://www.dxscan.io/

DxFarm: https://bsc.dxchain.com/home

DxChain Wiki: https://dxchain.gitbook.io/dxchain/
