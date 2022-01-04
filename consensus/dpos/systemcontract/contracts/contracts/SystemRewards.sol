// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma abicoder v2;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "./interfaces/IValidators.sol";
import "./Base.sol";
import "./Validators.sol";

contract SystemRewards is Base, Initializable {
    using EnumerableSet for EnumerableSet.UintSet;

    struct Reward {
        uint256 pengingValidatorReward;
        uint256 pengingDelegatorsReward;
        uint256 totalVotes;
        uint256 cancelVotes;
    }

    struct SysRewards {
        uint256[] epochs;
        mapping(uint256 => Reward) rewards;
        uint256 nextValRewardEpochIndex;
        uint256 nextDelegatorsRewardEpochIndex;
        uint256 accRewardPerVote;
    }

    mapping(address => SysRewards) public sysRewards;

    struct Epoch {
        uint256 blockReward;
        uint256 tvl;
        uint256 validatorCount;
    }

    /// @notice epoch => Epoch
    mapping(uint256 => Epoch) public epochs;

    IValidators public validatorC;

    INodeVote public nodeVoteC;

    /// @notice validator => epoch => punish limit
    mapping(address => mapping(uint256 => uint256)) public punishInfo;

    event LogEarnValidatorReward(address indexed _val, uint256 _amount);
    event LogDistributeBlockReward(address indexed _val, uint256 _valReward, uint256 _delegatorReward);
    event LogPunish(address indexed _val, bool _kickout, uint256 _reward);

    modifier onlyValidatorC() {
        require(msg.sender == address(validatorC), "SystemRewards: not Validator contract");
        _;
    }

    modifier onlyNodeVoteC() {
        require(msg.sender == address(nodeVoteC), "SystemRewards: not NodeVote contract");
        _;
    }

    /**
     * @dev initialize
     */
    function initialize(address _validator, address _node) public initializer {
        validatorC = IValidators(_validator);
        nodeVoteC = INodeVote(_node);
    }

    function getValRewardEpochs(address _val) public view returns (uint256[] memory) {
        return sysRewards[_val].epochs;
    }

    function getValRewardInfoByEpoch(address _val, uint256 _epoch) public view returns (Reward memory) {
        return sysRewards[_val].rewards[_epoch];
    }

    function pendingValReward(address _val) public view returns (uint256, uint256) {
        SysRewards storage sysReward = sysRewards[_val];
        if (sysReward.epochs.length == 0) return (0, 0);

        (uint256 avaliable, uint256 frozen, uint256 curEpoch) = (0, 0, currentEpoch());

        for (uint256 i = 0; i < sysReward.epochs.length - sysReward.nextValRewardEpochIndex; i++) {
            uint256 epoch = sysReward.epochs[sysReward.nextValRewardEpochIndex + i];
            if (epoch + 6 < curEpoch) {
                avaliable += sysReward.rewards[epoch].pengingValidatorReward;
            } else {
                frozen += sysReward.rewards[epoch].pengingValidatorReward;
            }
        }
        return (avaliable, frozen);
    }

    function pendingVoterReward(address _val) public view returns (uint256, uint256) {
        SysRewards storage sysReward = sysRewards[_val];
        if (sysReward.epochs.length == 0) return (0, 0);
        (uint256 sum, uint256 accRewardPerVote, uint256 curEpoch) = (0, 0, currentEpoch());

        for (uint256 i = sysReward.nextDelegatorsRewardEpochIndex; i < sysReward.epochs.length; i++) {
            uint256 epoch = sysReward.epochs[i];
            if (epoch == curEpoch || sysReward.rewards[epoch].totalVotes == 0) continue;

            uint256 curEpochToBurn = (sysReward.rewards[epoch].pengingDelegatorsReward *
                sysReward.rewards[epoch].cancelVotes) / sysReward.rewards[epoch].totalVotes;
            accRewardPerVote += sysReward.rewards[epoch].pengingDelegatorsReward / sysReward.rewards[epoch].totalVotes;
            sum += sysReward.rewards[epoch].pengingDelegatorsReward - curEpochToBurn;
        }

        return (sum, accRewardPerVote);
    }

    /**
     * @dev earn validator reward
     */
    function earnValReward() external {
        _earnValidatorReward(msg.sender);
    }

    /**
     * @dev earn validator reward from Validator contract
     */
    function earnValRewardFromValidatorC(address _val) external onlyValidatorC {
        _earnValidatorReward(_val);
    }

    /**
     * @dev distribute BlockReward while mining
     */
    function distributeBlockReward(uint256 _reward) external payable onlyMiner onlySystem {
        // get validator distribution rate
        Validators.Validator memory val = validatorC.validators(msg.sender);

        // calculate block reward
        uint256 delegatorsReward = (_reward * uint256(val.rate)) / uint256(MAX_RATE);
        uint256 valReward = _reward - delegatorsReward;

        // distribute block reward
        SysRewards storage sysReward = sysRewards[msg.sender];
        uint256 curEpoch = sysReward.epochs[sysReward.epochs.length - 1];
        sysReward.rewards[curEpoch].pengingValidatorReward += valReward;
        sysReward.rewards[curEpoch].pengingDelegatorsReward += delegatorsReward;

        emit LogDistributeBlockReward(msg.sender, valReward, delegatorsReward);
    }

    function punish(address _val) external onlyMiner onlySystem returns (bool) {
        uint256 curEpoch = currentEpoch();
        punishInfo[_val][curEpoch] += 1;
        if (punishInfo[_val][curEpoch] > MAX_PUNISH_COUNT) {
            // Deduct this epoch reward
            SysRewards storage sysReward = sysRewards[_val];
            uint256 curEpochValReward = sysReward.rewards[curEpoch].pengingValidatorReward;
            sysReward.rewards[curEpoch].pengingValidatorReward = 0;
            payable(BLACK_HOLE).transfer(curEpochValReward);
            // kickout validator
            validatorC.kickoutValidator(_val);
            emit LogPunish(_val, true, curEpochValReward);
            return true;
        }

        emit LogPunish(_val, false, 0);
        return false;
    }

    function updateVoterReward(address _val) external onlyNodeVoteC returns (uint256, uint256) {
        SysRewards storage sysReward = sysRewards[_val];
        if (sysReward.epochs.length == 0) return (0, 0);
        (uint256 sum, uint256 updateEpochIndex, uint256 toBurn, uint256 curEpoch) = (0, 0, 0, currentEpoch());

        for (uint256 i = sysReward.nextDelegatorsRewardEpochIndex; i < sysReward.epochs.length; i++) {
            uint256 epoch = sysReward.epochs[i];
            updateEpochIndex = i + 1;

            if (epoch == curEpoch || sysReward.rewards[epoch].totalVotes == 0) {
                toBurn += sysReward.rewards[epoch].pengingDelegatorsReward;
                continue;
            }

            uint256 curEpochToBurn = (sysReward.rewards[epoch].pengingDelegatorsReward *
                sysReward.rewards[epoch].cancelVotes) / sysReward.rewards[epoch].totalVotes;
            sysReward.accRewardPerVote +=
                sysReward.rewards[epoch].pengingDelegatorsReward /
                sysReward.rewards[epoch].totalVotes;
            toBurn += curEpochToBurn;
            sum += sysReward.rewards[epoch].pengingDelegatorsReward - curEpochToBurn;
        }

        sysReward.nextDelegatorsRewardEpochIndex = updateEpochIndex;

        payable(BLACK_HOLE).transfer(toBurn);
        payable(address(nodeVoteC)).transfer(sum);

        return (sum, sysReward.accRewardPerVote);
    }

    function updateCancelVote(address _val, uint256 _amount) external onlyNodeVoteC {
        SysRewards storage sysReward = sysRewards[_val];
        sysReward.rewards[currentEpoch()].cancelVotes += _amount;
    }

    function updateValDeposit(address _val, uint256 _deposit) external onlyValidatorC {
        SysRewards storage sysReward = sysRewards[_val];
        uint256 curEpoch = currentEpoch();
        sysReward.epochs.push(curEpoch);
        sysReward.rewards[curEpoch].totalVotes = _deposit;
    }

    function updateEpochWhileElect(uint256 _tvl, uint256 _valCount) external onlyValidatorC {
        uint256 curEpoch = currentEpoch();
        epochs[curEpoch].tvl = _tvl;
        epochs[curEpoch].validatorCount = _valCount;
        _calculateBlockReward();
    }

    function _earnValidatorReward(address _val) internal {
        SysRewards storage sysReward = sysRewards[_val];
        if (sysReward.epochs.length == 0) return;

        (uint256 avaliable, uint256 updateEpochIndex, uint256 curEpoch) = (0, 0, currentEpoch());
        for (uint256 i = sysReward.nextValRewardEpochIndex; i < sysReward.epochs.length; i++) {
            uint256 epoch = sysReward.epochs[i];
            if (epoch + 6 < curEpoch) {
                avaliable += sysReward.rewards[epoch].pengingValidatorReward;
                updateEpochIndex = i + 1;
            }
        }
        sysReward.nextValRewardEpochIndex = updateEpochIndex;
        payable(_val).transfer(avaliable);

        emit LogEarnValidatorReward(_val, avaliable);
    }

    function _calculateBlockReward() internal {
        (uint256 curEpoch, uint256 sum) = (currentEpoch(), 0);
        uint256 count = curEpoch >= 13 ? 14 : curEpoch + 1;
        for (uint256 i = curEpoch + 1 - count; i <= curEpoch; i++) {
            sum += epochs[i].tvl;
        }
        uint256 avg = sum / count;
        if (avg < TOTAL_DEPOSIT_150) {
            epochs[curEpoch].blockReward = REWARD_DEPOSIT_UNDER_150;
        } else if (TOTAL_DEPOSIT_150 <= avg && avg < TOTAL_DEPOSIT_200) {
            epochs[curEpoch].blockReward = REWARD_DEPOSIT_FROM_150_TO_200;
        } else if (TOTAL_DEPOSIT_200 <= avg && avg < TOTAL_DEPOSIT_250) {
            epochs[curEpoch].blockReward = REWARD_DEPOSIT_FROM_200_TO_250;
        } else if (TOTAL_DEPOSIT_250 <= avg && avg < TOTAL_DEPOSIT_300) {
            epochs[curEpoch].blockReward = REWARD_DEPOSIT_FROM_250_TO_300;
        } else if (TOTAL_DEPOSIT_300 <= avg && avg < TOTAL_DEPOSIT_350) {
            epochs[curEpoch].blockReward = REWARD_DEPOSIT_FROM_300_TO_350;
        } else {
            epochs[curEpoch].blockReward = REWARD_DEPOSIT_OVER_350;
        }
    }
}
