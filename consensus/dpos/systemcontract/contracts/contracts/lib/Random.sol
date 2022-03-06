// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library Random {
    /**
     * @dev generate random seed
     */
    function makeSeed(uint256 i) public view returns (uint64) {
        return uint64(uint256(keccak256(abi.encodePacked(blockhash(block.number - 1)))) + i);
    }

    function luckyWheel() public returns (address[] memory) {}
}
