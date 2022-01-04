// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma abicoder v2;

import "../Validators.sol";

interface IValidators {
    function isEffictiveValidator(address addr) external view returns (bool);

    function getEffictiveValidators() external view returns (address[] memory);

    function getInvalidValidators() external view returns (address[] memory);

    function effictiveValsLength() external view returns (uint256);

    function invalidValsLength() external view returns (uint256);

    function validators(address _val) external view returns (Validators.Validator calldata);

    function kickoutValidator(address _val) external;

    function addValidatorFromProposal(
        address addr,
        uint256 deposit,
        uint8 rate
    ) external payable;

    function voteValidator(
        address _voter,
        address _val,
        uint256 _votes
    ) external payable;

    function cancelVoteValidator(
        address _voter,
        address _val,
        uint256 _votes
    ) external payable;
}
