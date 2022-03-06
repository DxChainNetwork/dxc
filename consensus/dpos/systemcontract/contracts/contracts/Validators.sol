// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "./Base.sol";
import "./interfaces/ISystemRewards.sol";
import "./interfaces/IProposals.sol";
import "./interfaces/INodeVote.sol";

contract Validators is Base, Initializable {
    using EnumerableSet for EnumerableSet.AddressSet;
    using Address for address;

    enum ValidatorStatus {
        canceled,
        canceling,
        kickout,
        effictive
    }

    address[] public curEpochValidators;
    mapping(address => uint256) public curEpochValidatorsIdMap;

    EnumerableSet.AddressSet effictiveValidators;

    /// @notice canceled、canceling、kickout
    EnumerableSet.AddressSet invalidValidators;

    struct Validator {
        ValidatorStatus status;
        uint256 deposit;
        /// @notice based on 100
        uint8 rate;
        uint256 totalVotes;
        uint256 unstakeLockingEndBlock;
        uint256 rateSettLockingEndBlock;
    }

    mapping(address => Validator) public validators;

    mapping(address => EnumerableSet.AddressSet) validatorToVoters;

    /// @notice TVL
    uint256 public totalDeposit;

    mapping(address => uint256) validatorsDepositMap;

    /// @notice SystemRewards contract
    ISystemRewards public sysRewards;

    /// @notice Proposals contract
    IProposals public proposals;

    /// @notice NodeVote contract
    INodeVote public nodeVote;

    event LogAddValidator(address indexed _val, uint256 _deposit, uint256 _rate);
    event LogUpdateValidatorDeposit(address indexed _val, uint256 _deposit);
    event LogUpdateValidatorRate(address indexed _val, uint8 _rate);
    event LogKickoutValidator(address indexed _val);
    event LogUnstakeValidator(address indexed _val, uint256 _lockEnd);
    event LogRedeemValidator(address indexed _val);
    event LogTryElect(uint256 _validatorCount, uint256 _deposit);

    modifier onlyNotCurEpochValidator() {
        require(
            validators[msg.sender].status == ValidatorStatus.effictive && curEpochValidatorsIdMap[msg.sender] == 0,
            "Validators: illegal msg.sender"
        );
        _;
    }

    /**
     * @dev only Proposals contract address
     */
    modifier onlyProposalsC() {
        require(msg.sender == address(proposals), "Validators: not Proposals contract address");
        _;
    }

    /**
     * @dev only SystemRewards contract address
     */
    modifier onlySysRewardsC() {
        require(msg.sender == address(sysRewards), "Validators: not SystemRewards contract address");
        _;
    }

    /**
     * @dev only NodeVote contract address
     */
    modifier onlyNodeVoteC() {
        require(msg.sender == address(nodeVote), "Validators: not NodeVote contract address");
        _;
    }

    modifier onlyValidator() {
        require(
            effictiveValidators.contains(msg.sender) || invalidValidators.contains(msg.sender),
            "Validators: Not Validator"
        );
        _;
    }

    /**
     * @dev initialize
     */
    function initialize(
        address _proposal,
        address _sysReward,
        address _nodeVote,
        address _initVal,
        uint256 _initDeposit,
        uint8 _initRate
    ) external payable initializer {
        sysRewards = ISystemRewards(_sysReward);
        proposals = IProposals(_proposal);
        nodeVote = INodeVote(_nodeVote);

        require(!_initVal.isContract(), "Validators: validator address error");
        require(msg.value == _initDeposit && _initDeposit >= MIN_DEPOSIT, "Validators: deposit or value error");
        require(
            _initRate >= MIN_RATE && _initRate <= MAX_RATE,
            "Validators: Rate must greater than MIN_RATE and less than MAX_RATE"
        );

        Validator storage val = validators[_initVal];
        val.status = ValidatorStatus.effictive;
        val.deposit = _initDeposit;
        val.totalVotes += _initDeposit;
        val.rate = _initRate;

        effictiveValidators.add(_initVal);
        totalDeposit += _initDeposit;

        curEpochValidators.push(_initVal);
        curEpochValidatorsIdMap[_initVal] = curEpochValidators.length;
        sysRewards.updateValDeposit(_initVal, _initDeposit);
        sysRewards.updateEpochWhileElect(totalDeposit, curEpochValidators.length);

        emit LogAddValidator(_initVal, _initDeposit, _initRate);
    }

    function getCurEpochValidators() external view returns (address[] memory) {
        return curEpochValidators;
    }

    function isEffictiveValidator(address addr) external view returns (bool) {
        return validators[addr].status == ValidatorStatus.effictive;
    }

    function effictiveValsLength() public view returns (uint256) {
        return effictiveValidators.length();
    }

    function getEffictiveValidators() public view returns (address[] memory) {
        uint256 len = effictiveValidators.length();
        address[] memory vals = new address[](len);

        for (uint256 i = 0; i < len; i++) {
            vals[i] = effictiveValidators.at(i);
        }
        return vals;
    }

    function invalidValsLength() public view returns (uint256) {
        return invalidValidators.length();
    }

    function getInvalidValidators() public view returns (address[] memory) {
        uint256 len = invalidValidators.length();
        address[] memory vals = new address[](len);

        for (uint256 i = 0; i < len; i++) {
            vals[i] = invalidValidators.at(i);
        }
        return vals;
    }

    function updateValidatorDeposit(uint256 _deposit) external payable onlyNotCurEpochValidator {
        Validator storage val = validators[msg.sender];
        if (_deposit >= val.deposit) {
            require(msg.value >= _deposit - val.deposit, "Validators: illegal deposit");
            val.deposit = val.deposit + msg.value;
            val.totalVotes = val.totalVotes + msg.value;
            totalDeposit += (_deposit - val.deposit);
            payable(msg.sender).transfer(_deposit - val.deposit);
        } else {
            require(_deposit >= MIN_DEPOSIT, "Validators: illegal deposit");
            uint256 sub = val.deposit - _deposit;
            payable(msg.sender).transfer(sub);
            val.deposit = _deposit;
            val.totalVotes = val.totalVotes - sub;
            totalDeposit -= (val.deposit - _deposit);
        }

        _earnValidatorReward();

        emit LogUpdateValidatorDeposit(msg.sender, val.deposit);
    }

    function updateValidatorRate(uint8 _rate) external onlyNotCurEpochValidator {
        Validator storage val = validators[msg.sender];

        require(val.rateSettLockingEndBlock < block.number, "Validators: illegal rate set block");
        require(_rate >= MIN_RATE && val.rate <= MAX_RATE, "Validators: illegal Allocation ratio");

        if (val.rate > _rate) {
            require(
                (uint256(val.rate - _rate) * 10000) / uint256(val.rate) <= uint256(MAX_RATE_OF_CHANGE),
                "Validators: illegal rate of change"
            );
        } else {
            require(
                (uint256(_rate - val.rate) * 10000) / uint256(val.rate) <= uint256(MAX_RATE_OF_CHANGE),
                "Validators: illegal rate of change"
            );
        }

        _earnValidatorReward();

        val.rate = _rate;
        val.rateSettLockingEndBlock = block.number + RATE_SET_LOCK_BLOCKS;

        emit LogUpdateValidatorRate(msg.sender, _rate);
    }

    function addValidatorFromProposal(
        address _val,
        uint256 _deposit,
        uint8 _rate
    ) external payable onlyProposalsC {
        require(!_val.isContract(), "Validators: validator address error");
        require(msg.value == _deposit, "Validators: deposit not equal msg.value");

        Validator storage val = validators[_val];
        require(
            val.status == ValidatorStatus.canceled || val.status == ValidatorStatus.kickout,
            "Validators: validator status error"
        );

        val.status = ValidatorStatus.effictive;
        val.deposit = _deposit;
        val.totalVotes += _deposit;
        val.rate = _rate;

        effictiveValidators.add(_val);
        invalidValidators.remove(_val);
        totalDeposit += _deposit;

        emit LogAddValidator(_val, _deposit, _rate);
    }

    function kickoutValidator(address _val) external onlySysRewardsC {
        Validator storage val = validators[_val];
        require(val.status == ValidatorStatus.effictive, "Validators: validator status error");

        val.status = ValidatorStatus.kickout;
        uint256 index = curEpochValidatorsIdMap[_val] - 1;
        for (uint256 i = index; i < curEpochValidators.length - 1 - index; i++) {
            address nextAddr = curEpochValidators[index + 1];
            uint256 nextAddrId = curEpochValidatorsIdMap[nextAddr];
            curEpochValidatorsIdMap[nextAddr] = nextAddrId - 1;
            curEpochValidators[index] = curEpochValidators[index + 1];
        }
        curEpochValidators.pop();

        effictiveValidators.remove(_val);
        invalidValidators.add(_val);

        emit LogKickoutValidator(_val);
    }

    function unstake() external {
        Validator storage val = validators[msg.sender];
        require(
            (val.status == ValidatorStatus.effictive || val.status == ValidatorStatus.kickout) &&
                curEpochValidatorsIdMap[msg.sender] == 0,
            "Validators: illegal msg.sender"
        );

        _earnValidatorReward();

        val.status = ValidatorStatus.canceling;
        val.unstakeLockingEndBlock = block.number + VALIDATOR_UNSTAKE_LOCK_BLOCKS;

        effictiveValidators.remove(msg.sender);
        invalidValidators.add(msg.sender);

        emit LogUnstakeValidator(msg.sender, val.unstakeLockingEndBlock);
    }

    function redeem() external {
        Validator storage val = validators[msg.sender];
        require(val.unstakeLockingEndBlock < block.number, "Validators: illegal redeem block");
        require(
            val.status == ValidatorStatus.canceling && curEpochValidatorsIdMap[msg.sender] == 0,
            "Validators: illegal msg.sender"
        );

        _earnValidatorReward();

        val.status = ValidatorStatus.canceled;
        val.totalVotes -= val.deposit;
        val.deposit = 0;

        totalDeposit -= val.deposit;
        payable(msg.sender).transfer(val.deposit);

        emit LogRedeemValidator(msg.sender);
    }

    function voteValidator(
        address _voter,
        address _val,
        uint256 _votes
    ) external payable onlyNodeVoteC {
        validators[_val].totalVotes += _votes;
        totalDeposit += _votes;
        validatorToVoters[_val].add(_voter);
    }

    function cancelVoteValidator(
        address _voter,
        address _val,
        uint256 _votes
    ) external onlyNodeVoteC {
        validators[_val].totalVotes -= _votes;
        totalDeposit -= _votes;
        payable(msg.sender).transfer(_votes);
        validatorToVoters[_val].remove(_voter);
    }

    function tryElect() external onlySystem onlyMiner {
        uint256 nextEpochValCount = nextEpochValidatorCount();

        for (uint256 i = 0; i < curEpochValidators.length; i++) {
            delete curEpochValidatorsIdMap[curEpochValidators[i]];
        }
        delete curEpochValidators;

        uint256 total = makeValToDepositMap();

        emit LogTryElect(nextEpochValCount, total);

        for (uint256 i = 0; i < nextEpochValCount; i++) {
            if (total == 0) {
                break;
            }
            uint256 randDeposit = rand(total, i);
            for (uint256 j = 0; j < effictiveValidators.length(); j++) {
                address val = effictiveValidators.at(j);
                uint256 deposit = validatorsDepositMap[val];
                if (curEpochValidatorsIdMap[val] != 0) continue;
                if (randDeposit <= deposit) {
                    curEpochValidators.push(val);
                    curEpochValidatorsIdMap[val] = curEpochValidators.length;
                    total -= deposit;
                    sysRewards.updateValDeposit(val, deposit);
                    break;
                }
                randDeposit -= deposit;
            }
        }

        sysRewards.updateEpochWhileElect(totalDeposit, curEpochValidators.length);
    }

    function makeValToDepositMap() internal returns (uint256 total) {
        for (uint256 i = 0; i < effictiveValidators.length(); i++) {
            delete validatorsDepositMap[effictiveValidators.at(i)];
        }
        for (uint256 i = 0; i < effictiveValidators.length(); i++) {
            address val = effictiveValidators.at(i);
            validatorsDepositMap[val] = validators[val].totalVotes;
            total += validators[val].totalVotes;
        }
    }

    function rand(uint256 _length, uint256 _i) internal view returns (uint256) {
        uint256 random = uint256(keccak256(abi.encodePacked(blockhash(block.number - _i), _i)));
        return random % _length;
    }

    function recentFourteenEpochAvgValCount() internal view returns (uint256) {
        uint256 curEpoch = sysRewards.currentEpoch();
        if (curEpoch == 1) {
            return effictiveValidators.length();
        }
        uint256 sumValidatorCount = 0;
        uint256 avg = 15;
        if (curEpoch < avg) {
            avg = curEpoch;
        }
        for (uint256 i = 1; i < avg; i++) {
            (, , uint256 validatorCount) = sysRewards.epochs(curEpoch - i);
            sumValidatorCount += validatorCount;
        }
        return sumValidatorCount / (avg - 1);
    }

    function nextEpochValidatorCount() internal view returns (uint256) {
        uint256 avgCount = recentFourteenEpochAvgValCount();
        if (avgCount < MIN_LEVEL_VALIDATOR_COUNT) {
            return MaxValidatorSizeUnder60;
        }
        if (avgCount < MEDIUM_LEVEL_VALIDATOR_COUNT) {
            return MaxValidatorSizeFrom60To90;
        }
        if (avgCount < MAX_LEVEL_VALIDATOR_COUNT) {
            return MaxValidatorSizeFrom90To120;
        }
        // avgCount >= MAX_LEVEL_VALIDATOR_COUNT
        return MaxValidatorSizeOver120;
    }

    function _earnValidatorReward() internal {
        sysRewards.earnValRewardFromValidatorC(msg.sender);
    }
}
