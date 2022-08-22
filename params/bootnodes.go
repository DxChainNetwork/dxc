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
	"enode://02c88819bd547e22e982c6811dc4eba6807e19d4a6827af2a20d467ae671e183c92c1fe9b93d2516a08e421a03559ef71fa9d766d7dcdde0b6dfb6e967171a66@54.68.71.166:32668",
	"enode://600856e36cf714659be073b8fd4fc148fbf303e4303bae2e3e97927e7031f5254a91310a6f93f1facbc9775f6af023eb4285d94f61323e853548ab8e9478c9dd@52.10.116.84:32668",
	"enode://a7285394350ac2eb7c12ad37e5c0cb9450824066ed13f079dabc7f5383b6f16c57ab6c2639418eab3606fc983c350ffdde77279d777b50883957266634d98460@54.71.25.84:32668",
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
