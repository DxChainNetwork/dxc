// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Base {
    uint256 public constant BLOCK_SECONDS = 6;
    /// @notice min rate. base on 100
    uint8 public constant MIN_RATE = 70;
    /// @notice max rate. base on 100
    uint8 public constant MAX_RATE = 100;

    /// @notice 10 * 60 / BLOCK_SECONDS
    uint256 public constant EPOCH_BLOCKS = 14400;
    /// @notice min deposit for validator
    uint256 public constant MIN_DEPOSIT = 4e7 ether;
    uint256 public constant MAX_PUNISH_COUNT = 139;

    /// @notice use blocks as units in code: RATE_SET_LOCK_EPOCHS * EPOCH_BLOCKS
    uint256 public constant RATE_SET_LOCK_EPOCHS = 1;
    /// @notice use blocks as units in code: VALIDATOR_UNSTAKE_LOCK_EPOCHS * EPOCH_BLOCKS
    uint256 public constant VALIDATOR_UNSTAKE_LOCK_EPOCHS = 1;
    /// @notice use blocks as units in code: PROPOSAL_DURATION_EPOCHS * EPOCH_BLOCKS
    uint256 public constant PROPOSAL_DURATION_EPOCHS = 7;
    /// @notice use epoch as units in code: VALIDATOR_REWARD_LOCK_EPOCHS
    uint256 public constant VALIDATOR_REWARD_LOCK_EPOCHS = 7;
    /// @notice use epoch as units in code: VOTE_CANCEL_EPOCHS
    uint256 public constant VOTE_CANCEL_EPOCHS = 1;

    uint256 public constant MAX_VALIDATORS_COUNT = 210;
    uint256 public constant MAX_VALIDATOR_DETAIL_LENGTH = 1000;
    uint256 public constant MAX_VALIDATOR_NAME_LENGTH = 100;

    // total deposit
    uint256 public constant TOTAL_DEPOSIT_LV1 = 1e18 * 1e8 * 150;
    uint256 public constant TOTAL_DEPOSIT_LV2 = 1e18 * 1e8 * 200;
    uint256 public constant TOTAL_DEPOSIT_LV3 = 1e18 * 1e8 * 250;
    uint256 public constant TOTAL_DEPOSIT_LV4 = 1e18 * 1e8 * 300;
    uint256 public constant TOTAL_DEPOSIT_LV5 = 1e18 * 1e8 * 350;

    // block reward
    uint256 public constant REWARD_DEPOSIT_UNDER_LV1 = 1e15 * 95250;
    uint256 public constant REWARD_DEPOSIT_FROM_LV1_TO_LV2 = 1e15 * 128250;
    uint256 public constant REWARD_DEPOSIT_FROM_LV2_TO_LV3 = 1e15 * 157125;
    uint256 public constant REWARD_DEPOSIT_FROM_LV3_TO_LV4 = 1e15 * 180750;
    uint256 public constant REWARD_DEPOSIT_FROM_LV4_TO_LV5 = 1e15 * 199875;
    uint256 public constant REWARD_DEPOSIT_OVER_LV5 = 1e15 * 214125;

    // validator count
    uint256 public constant MAX_VALIDATOR_COUNT_LV1 = 21;
    uint256 public constant MAX_VALIDATOR_COUNT_LV2 = 33;
    uint256 public constant MAX_VALIDATOR_COUNT_LV3 = 66;
    uint256 public constant MAX_VALIDATOR_COUNT_LV4 = 99;
    uint256 public constant MIN_LEVEL_VALIDATOR_COUNT = 60;
    uint256 public constant MEDIUM_LEVEL_VALIDATOR_COUNT = 90;
    uint256 public constant MAX_LEVEL_VALIDATOR_COUNT = 120;

    // dead address
    address public constant BLACK_HOLE_ADDRESS = 0x0000000000000000000000000000000000000000;

    uint256 public constant SAFE_MULTIPLIER = 1e18;

    modifier onlySystem() {
        require(tx.gasprice == 0, "Prohibit external calls");
        _;
    }

    modifier onlyMiner() {
        require(msg.sender == block.coinbase, "msg.sender error");
        _;
    }

    /**
     * @dev return current epoch
     */
    function currentEpoch() public view returns (uint256) {
        return block.number / EPOCH_BLOCKS;
    }
}
