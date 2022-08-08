// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import "github.com/DxChainNetwork/dxc/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	//"enode://a554c1cffa56f82e234d9acfc580914b34b32d45556c70935993454c3f83e44a82a5d27f6e7c5401b51faf5b151764379414717631333192cdc16cb6a7191833@13.215.48.4:32668",
	//"enode://c7ec3457a14f3a28f95ec1fd83d3d2be594b33535f3e39adbfd00969b29a1f4763160d7358636f849823f33792063eb5cae482f008565d174083ddc583c5c012@3.1.195.42:32668",
	//enode://882083551ef1d78f20caa4628027d7f0363ef6bb59148e23986e88df170b3cd75a4ec1c852b7b6d26b60f75d46f75e295e2d0b909f817ab54cc90067b9f6813d@127.0.0.1:32668
	"enode://882083551ef1d78f20caa4628027d7f0363ef6bb59148e23986e88df170b3cd75a4ec1c852b7b6d26b60f75d46f75e295e2d0b909f817ab54cc90067b9f6813d@13.229.247.212:32668",
	"enode://30e60336db8dee8b1deec00b894f954fbe55d3ce8efc5261a5cfa1dde00d9685ea2b438c3b7ae3399b07c27c7674388e88d469d20038949b50317ddf98b77716@13.215.163.181:32668",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
var TestnetBootnodes = []string{
	//"enode://6fd8199df5e00b091ac4d711200fa371b010b1b20828270d1bc78fcaaed6e76179e3cb49de919c34d965e28bceac6a33451db5d15a8caa2518e39e1925fb0d09@18.136.102.78:32668",
	//"enode://85b23ffafd9300176a6e6a3c8c1aaea0823b5f3a26052d5b7738cf526c80725ba38f81993388d467646cc73c078dad3d9d060a53f41043603c95d2670895d2a3@54.255.159.74:32668",
	//"enode://06a93f23ff0ebe15befd632ce6d4e23e067fe083b0f2c6a145ccc04dfa3dbb17785d0b79dc33a267b12c368704ccbd8570a9ddc53382ea5c73fbd0c143113896@13.229.92.245:32668",
}

var V5Bootnodes = []string{}

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	return ""
}
