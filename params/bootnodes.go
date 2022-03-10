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
	"enode://f664e6843fd5aab3a847d3a39e5f46fd86bfda22c888e423d83d5099fac1d43689e61a76cb23e792606ae5d7f07b05d0eb134803e4e8b994de5e11b0054e6974@18.141.209.219:32668",
	"enode://57c806e61d81aa678365e2568987a3b3803e9fa7e4dfa18c21929111fbbb41b3d9311237b072f462789b6a348a56ebf354e99ffc3d950ee7f4a73bbcb5b95f08@18.140.54.231:32668",
	"enode://181c289b212f53fa1d3cab3d8bc825b136bf6479fc641200b3a1c25ec7947adc45851a120eab63120744734fbfca816640b84dec056179dc8dd4b983d9ef3354@13.229.233.232:32668",
	"enode://798ee3a418198eb1689e6364940ced4c93c8490632660020d7bb46f666c74db4048ed00a20d11ba013daa1743eb9a3989f59cbc20730ff7d94425c4d412925bd@18.136.209.237:32668",
	"enode://ac197857bb25fa5068927f804503ea4312f2e01caa19f6f0d148c623ba98d4da1e737b89d054b58060014f76b74fa4d4b53c17dfebd577cdaa902c0bc18428b3@13.213.40.80:32668",
}

var V5Bootnodes = []string{}

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	return ""
}
