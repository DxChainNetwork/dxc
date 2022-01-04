// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface ISystemRewards {
    function earnValRewardFromValidatorC(address _val) external;

    function epochs(uint256 _epoch)
        external
        view
        returns (
            uint256,
            uint256,
            uint256
        );

    function currentEpoch() external view returns (uint256);

    function pendingVoterReward(address _val) external view returns (uint256, uint256);

    function updateVoterReward(address _val) external returns (uint256, uint256);

    function updateCancelVote(address _val, uint256 _amount) external;

    function updateValDeposit(address _val, uint256 _deposit) external;

    function updateEpochWhileElect(uint256 _tvl, uint256 _valCount) external;
}
