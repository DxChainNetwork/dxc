// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// Test data
contract Base {
    uint256 public constant EPOCH_BLOCKS = 20;

    uint8 public constant MIN_RATE = 70;
    uint8 public constant MAX_RATE = 100;

    uint16 public constant MAX_RATE_OF_CHANGE = 300;
    // mainnet : 4e7 ether
    uint256 public constant MIN_DEPOSIT = 1 ether;
    uint256 public constant RATE_SET_LOCK_BLOCKS = 100;
    uint256 public constant VALIDATOR_UNSTAKE_LOCK_BLOCKS = 100;
    uint256 public constant MAX_VALIDATORS_COUNT = 210;
    //TODO: need test
    uint256 public constant MAX_PROPOSAL_DETAIL_LENGTH = 3000;
    // redeem block limit
    uint256 public constant VOTE_CANCEL_BLOCK = 100;

    uint256 public constant MAX_PUNISH_COUNT = 3;

    // total deposit
    uint256 public constant TOTAL_DEPOSIT_150 = 1e18 * 1e8 * 150;
    uint256 public constant TOTAL_DEPOSIT_200 = 1e18 * 1e8 * 200;
    uint256 public constant TOTAL_DEPOSIT_250 = 1e18 * 1e8 * 250;
    uint256 public constant TOTAL_DEPOSIT_300 = 1e18 * 1e8 * 300;
    uint256 public constant TOTAL_DEPOSIT_350 = 1e18 * 1e8 * 350;

    // block reward
    uint256 public constant REWARD_DEPOSIT_UNDER_150 = 1e15 * 95250;
    uint256 public constant REWARD_DEPOSIT_FROM_150_TO_200 = 1e15 * 128250;
    uint256 public constant REWARD_DEPOSIT_FROM_200_TO_250 = 1e15 * 157125;
    uint256 public constant REWARD_DEPOSIT_FROM_250_TO_300 = 1e15 * 180750;
    uint256 public constant REWARD_DEPOSIT_FROM_300_TO_350 = 1e15 * 199875;
    uint256 public constant REWARD_DEPOSIT_OVER_350 = 1e15 * 214125;

    // validator count
    uint256 public constant MaxValidatorSizeUnder60 = 2;
    // uint256 public constant MaxValidatorSizeUnder60 = 21;
    uint256 public constant MaxValidatorSizeFrom60To90 = 33;
    uint256 public constant MaxValidatorSizeFrom90To120 = 66;
    uint256 public constant MaxValidatorSizeOver120 = 99;
    uint256 public constant MIN_LEVEL_VALIDATOR_COUNT = 60;
    uint256 public constant MEDIUM_LEVEL_VALIDATOR_COUNT = 90;
    uint256 public constant MAX_LEVEL_VALIDATOR_COUNT = 120;

    uint256 public constant PROPOSAL_DURATION = 201600;

    // dead address
    address public constant BLACK_HOLE = 0x0000000000000000000000000000000000000000;

    modifier onlySystem() {
        require(tx.gasprice == 0, "Prohibit external calls");
        _;
    }

    modifier onlyMiner() {
        require(msg.sender == block.coinbase, "msg.sender error");
        _;
    }

    /**
     * @dev get current epoch
     */
    function currentEpoch() public view returns (uint256) {
        return block.number / EPOCH_BLOCKS;
    }
}
