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
	"enode://f6e715cdacbaf28538aa93dbbd16c230bcdd0d1b055d56b6f155c50bc3b5f190ac48262a36e0e3fd62ec5c8ba10d994f445f33a501024d0266c55469dc7ee811@18.138.236.38:32668",
	"enode://ae9851d51ed679525ffed1d43c72b5f44dad81ffd47765786fbbf69af5d69c4bdba2fca0afc5599776edebc77dc62f2ee1c8b24664c0b0c8a58e897cad2260ed@13.214.153.30:32668",
	"enode://453a18b5c4fbd42d7cdade693304c5138c646af2684e44ba7e378e5d163913133557ac6e38bb0cff432036f01a057b14738325ca48271ea9797ffe9b9829fd62@13.212.159.171:32668",
	//"enode://a95fe8e7f8ed9fe98a90c7ee6cac677e05ebd5ecb82cb7c58ee1eb009aab06e6e429fa73e5ab3dc88d55362d2677203a5cd5275c210b1963e693362f20632b67@13.250.231.229:32668",
	//"enode://48532a4ad45b272ee3e5060ac20da50adea3d4c6a72fd9f137f204f755d9bfb3914921c1f1b5d4731a9283c368a8b73b0287be7a64a7d79f755b866a66c6e3f6@13.212.164.112:32668",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
var TestnetBootnodes = []string{
	"enode://6fd8199df5e00b091ac4d711200fa371b010b1b20828270d1bc78fcaaed6e76179e3cb49de919c34d965e28bceac6a33451db5d15a8caa2518e39e1925fb0d09@18.136.102.78:32668",
	"enode://85b23ffafd9300176a6e6a3c8c1aaea0823b5f3a26052d5b7738cf526c80725ba38f81993388d467646cc73c078dad3d9d060a53f41043603c95d2670895d2a3@54.255.159.74:32668",
	"enode://06a93f23ff0ebe15befd632ce6d4e23e067fe083b0f2c6a145ccc04dfa3dbb17785d0b79dc33a267b12c368704ccbd8570a9ddc53382ea5c73fbd0c143113896@13.229.92.245:32668",
}

var V5Bootnodes = []string{}

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	return ""
}
