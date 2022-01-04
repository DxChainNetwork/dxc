// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "./Base.sol";
import "./interfaces/IValidators.sol";
import "./interfaces/ISystemRewards.sol";

contract NodeVotes is Base, Initializable {
    using EnumerableSet for EnumerableSet.AddressSet;

    struct VoteInfo {
        address validator;
        uint256 votes;
        uint256 rewardDebt;
        uint256 updateRewardEpoch;
    }

    struct Voter {
        EnumerableSet.AddressSet validators;
        EnumerableSet.AddressSet redeemValidators;
        mapping(address => uint256) redeemBlock;
        mapping(address => VoteInfo) voteInfos;
    }

    struct RedeemVoterInfo {
        address validator;
        uint256 votes;
        uint256 redeemBlock;
    }

    mapping(address => Voter) voters;

    /// @notice Validators contract
    IValidators public validators;

    /// @notice SystemRewards contract
    ISystemRewards public sysRewards;

    event LogVote(address indexed voter, address indexed validator, uint256 votes);
    event LogCancelVote(address indexed voter, address indexed validator, uint256 votes);
    event LogRedeem(address indexed voter, uint256 votes);
    event LogEarn(address indexed voter, address indexed validator, uint256 reward);

    receive() external payable {}

    /**
     * @dev initialize
     */
    function initialize(address _validator, address _sysReward) public initializer {
        validators = IValidators(_validator);
        sysRewards = ISystemRewards(_sysReward);
    }

    function pendingVoteReward(address _val, address _voter) public view returns (uint256) {
        Voter storage voter = voters[_voter];
        if (!voter.validators.contains(_val)) return 0;
        (, uint256 accRewardPerVote) = sysRewards.pendingVoterReward(_val);

        if (voter.voteInfos[_val].updateRewardEpoch == currentEpoch()) {
            return 0;
        } else {
            return voter.voteInfos[_val].votes * accRewardPerVote - voter.voteInfos[_val].rewardDebt;
        }
    }

    function pendingRedeem() public view returns (uint256) {
        Voter storage voter = voters[msg.sender];
        if (voter.redeemValidators.length() == 0) {
            return 0;
        }
        uint256 allVotes = 0;

        for (uint256 i = 0; i < voter.redeemValidators.length(); i++) {
            address val = voter.redeemValidators.at(i);
            if (voter.redeemBlock[val] < block.number) {
                allVotes += voter.voteInfos[val].votes;
            }
        }
        return allVotes;
    }

    function redeemInfo(
        address _val,
        uint256 page,
        uint256 size
    ) public view returns (RedeemVoterInfo[] memory) {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        Voter storage voter = voters[_val];
        uint256 start = (page - 1) * size;
        if (voter.redeemValidators.length() < start) {
            size = 0;
        } else {
            uint256 length = voter.redeemValidators.length() - start;
            if (length < size) {
                size = length;
            }
        }
        RedeemVoterInfo[] memory redeemDir = new RedeemVoterInfo[](size);
        for (uint256 i = 0; i < size; i++) {
            address val = voter.redeemValidators.at(i + start);
            RedeemVoterInfo memory redeemVoter = RedeemVoterInfo({
                validator: val,
                votes: voter.voteInfos[val].votes,
                redeemBlock: voter.redeemBlock[val]
            });
            redeemDir[i] = redeemVoter;
        }
        return redeemDir;
    }

    function voteList(
        address _val,
        uint256 page,
        uint256 size
    ) public view returns (VoteInfo[] memory) {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        Voter storage voter = voters[_val];
        uint256 start = (page - 1) * size;
        if (voter.validators.length() < start) {
            size = 0;
        } else {
            uint256 length = voter.validators.length() - start;
            if (length < size) {
                size = length;
            }
        }
        VoteInfo[] memory voteDir = new VoteInfo[](size);
        for (uint256 i = 0; i < size; i++) {
            voteDir[i] = voter.voteInfos[(voter.validators.at(i + start))];
        }
        return voteDir;
    }

    function voteListLength(address _val) public view returns (uint256) {
        return voters[_val].validators.length();
    }

    function cancelVoteValidatorList(
        address _val,
        uint256 page,
        uint256 size
    ) public view returns (address[] memory) {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        Voter storage voter = voters[_val];
        uint256 start = (page - 1) * size;
        if (voter.redeemValidators.length() < start) {
            size = 0;
        } else {
            uint256 length = voter.redeemValidators.length() - start;
            if (length < size) {
                size = length;
            }
        }

        address[] memory cancelValidators = new address[](size);
        for (uint256 i = 0; i < size; i++) {
            cancelValidators[i] = voter.redeemValidators.at(i + start);
        }
        return cancelValidators;
    }

    function cancelVoteValidatorListLength(address _val) public view returns (uint256) {
        return voters[_val].redeemValidators.length();
    }

    function earn(address _val) public {
        sysRewards.updateVoterReward(_val);
        Voter storage voter = voters[msg.sender];
        uint256 curEpoch = currentEpoch();
        require(
            voter.voteInfos[_val].updateRewardEpoch != curEpoch,
            "NodeVotes: Have already received rewards this epoch"
        );
        uint256 reward = pendingVoteReward(_val, msg.sender);
        voter.voteInfos[_val].rewardDebt = reward;
        voter.voteInfos[_val].updateRewardEpoch = curEpoch;
        payable(msg.sender).transfer(reward);
        emit LogEarn(msg.sender, _val, reward);
    }

    function vote(address _val) external payable {
        require(!validators.isEffictiveValidator(msg.sender), "NodeVotes: The msg.sender can not be validator");
        require(!Address.isContract(msg.sender), "NodeVotes: The msg.sender can not be contract address");
        require(validators.isEffictiveValidator(_val), "NodeVotes: The val must be validator");
        require(msg.value > 0, "NodeVotes: Vote must greater than zero");

        earn(_val);

        Voter storage voter = voters[msg.sender];
        require(!voter.redeemValidators.contains(_val), "NodeVotes: The validator is cancel voting");

        if (voter.validators.contains(_val)) {
            VoteInfo storage _voteInfo = voter.voteInfos[_val];
            _voteInfo.votes += msg.value;
        } else {
            voter.validators.add(_val);
            VoteInfo memory _voteInfo;
            _voteInfo.votes = msg.value;
            _voteInfo.validator = _val;
            voter.voteInfos[_val] = _voteInfo;
        }

        validators.voteValidator{value: msg.value}(msg.sender, _val, msg.value);
        emit LogVote(msg.sender, _val, msg.value);
    }

    function cancelVote(address _val) external payable {
        require(!validators.isEffictiveValidator(msg.sender), "NodeVotes: The msg.sender can not be validator");

        earn(_val);

        Voter storage voter = voters[msg.sender];
        require(voter.validators.contains(_val), "NodeVotes: The msg.sender did not vote this validator");
        voter.validators.remove(_val);
        voter.redeemValidators.add(_val);
        voter.redeemBlock[_val] = block.number + VOTE_CANCEL_BLOCK;
        uint256 votes = voter.voteInfos[_val].votes;

        validators.cancelVoteValidator(msg.sender, _val, votes);
        sysRewards.updateCancelVote(_val, votes);
        emit LogCancelVote(msg.sender, _val, votes);
    }

    function redeem(address[] memory vals) external {
        Voter storage voter = voters[msg.sender];
        uint256 allVotes = 0;

        for (uint256 i = 0; i < vals.length; i++) {
            address val = vals[i];
            require(voter.redeemValidators.contains(val), "NodeVotes: can not redeem");
            if (voter.redeemBlock[val] < block.number) {
                allVotes += voter.voteInfos[val].votes;
                voter.redeemValidators.remove(val);
                delete voter.voteInfos[val];
                delete voter.redeemBlock[val];
            }
        }
        payable(msg.sender).transfer(allVotes);
        emit LogRedeem(msg.sender, allVotes);
    }
}
