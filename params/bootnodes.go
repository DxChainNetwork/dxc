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
	"enode://a3c8d5d31d9f3dc0bb29540c9b6ce78cb50e7ddc1f1010fed5d9e58703c9d9b9ed7261863f14722f1107e45ef3b4f34d5bef8912283ee14c70317919876dce34@13.213.60.136:32668",
	"enode://72b1af6a9d863f82e9a0fb3b106e7cbaf4387513e8d6268c1e46010df41494af2b75940f388b79a3a26e8063f2a70771a2fe4b079c24bfcb3b56e7c39de9807a@18.136.202.227:32668",
	"enode://2ef0dec9583731ebec2f6cfa62f91f22e8d13643501f4048e897f2ad2466a69b6a8f25178b247986f0d49c53db17c43f21be5b9586065e46bcc5349eb149b8ed@54.255.243.52:32668",
}

var V5Bootnodes = []string{}

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	return ""
}
