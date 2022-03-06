// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma abicoder v2;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "./Base.sol";
import "./interfaces/IValidators.sol";

contract Proposals is Base, Initializable {
    using EnumerableSet for EnumerableSet.Bytes32Set;

    enum ProposalType {
        init,
        recover
    }

    enum ProposalStatus {
        pending,
        pass,
        cancel
    }

    struct ProposalInfo {
        bytes4 id;
        address proposer;
        ProposalType pType;
        uint256 deposit;
        uint8 rate;
        string details;
        uint256 initBlock;
        address guarantee;
        uint256 updateBlock;
        ProposalStatus status;
    }

    mapping(bytes4 => ProposalInfo) public proposalInfos;
    mapping(address => bytes4[]) public proposals;

    EnumerableSet.Bytes32Set proposalsBytes;

    /// @notice Validators contract
    IValidators public validators;

    event LogInitProposal(bytes32 indexed id, address indexed proposer, uint256 block, uint256 deposit, uint256 rate);
    event LogGuarantee(bytes32 indexed id, address indexed guarantee, uint256 block);
    event LogCancelProposal(bytes32 indexed id, address indexed proposer, uint256 block);
    event LogUpdateProposal(bytes32 indexed id, address indexed proposer, uint256 block, uint256 deposit, uint256 rate);

    modifier onlyEffictiveValidator() {
        require(
            validators.isEffictiveValidator(msg.sender) || validators.effictiveValsLength() == 0,
            "Proposals: msg sender must be validator"
        );
        _;
    }

    modifier onlyEffictiveProposal(bytes4 id) {
        require(block.number <= proposalInfos[id].initBlock + PROPOSAL_DURATION, "Proposals: Proposal has expired");
        _;
    }

    modifier checkValidatorLength() {
        require(
            validators.effictiveValsLength() < MAX_VALIDATORS_COUNT,
            "Proposals: length of the validator must be less than MAX_VALIDATORS_COUNT"
        );
        _;
    }

    /**
     * @dev initialize
     */
    function initialize(address _validator) public initializer {
        validators = IValidators(_validator);
    }

    /**
     * @dev initProposal
     */
    function initProposal(
        ProposalType pType,
        uint8 rate,
        string memory details
    ) external payable checkValidatorLength {
        // check msg.sender
        require(!validators.isEffictiveValidator(msg.sender), "Proposals: The msg.sender can not be validator");
        require(!Address.isContract(msg.sender), "Proposals: The msg.sender can not be contract address");
        // check rate 、deposit(msg.value)、deposit
        require(bytes(details).length <= MAX_PROPOSAL_DETAIL_LENGTH, "Proposals: Details is too long");
        require(msg.value >= MIN_DEPOSIT, "Proposals: Deposit must greater than MIN_DEPOSIT");
        require(
            rate >= MIN_RATE && rate <= MAX_RATE,
            "Proposals: Rate must greater than MIN_RATE and less than MAX_RATE"
        );

        bytes4[] memory lastIds = proposals[msg.sender];
        if (lastIds.length > 0) {
            bytes4 lastId = lastIds[lastIds.length - 1];
            require(
                proposalInfos[lastId].status != ProposalStatus.pending,
                "Proposals: The msg.sender's latest proposal is still in pending"
            );
        }
        bytes4 id = bytes4(keccak256(abi.encodePacked(msg.sender, msg.value, rate, details, block.number)));
        require(proposalInfos[id].initBlock == 0, "Proposals: Proposal already exists");
        // new ProposalInfo
        ProposalInfo memory proposal;
        proposal.deposit = msg.value;
        proposal.id = id;
        proposal.details = details;
        proposal.initBlock = block.number;
        proposal.proposer = msg.sender;
        proposal.pType = pType;
        proposal.status = ProposalStatus.pending;
        proposal.rate = rate;
        proposalInfos[id] = proposal;
        proposals[address(msg.sender)].push(id);
        proposalsBytes.add(id);
        emit LogInitProposal(id, msg.sender, block.number, msg.value, rate);
    }

    /**
     * @dev guarantee
     */
    function guarantee(bytes4 id) external onlyEffictiveValidator onlyEffictiveProposal(id) {
        require(proposalInfos[id].initBlock != 0, "Proposals: proposal not exist");
        // check status
        require(
            proposalInfos[id].status == ProposalStatus.pending,
            "Proposals: The status of proposal must be pending"
        );
        proposalInfos[id].updateBlock = block.number;
        proposalInfos[id].guarantee = msg.sender;
        validators.addValidatorFromProposal{value: proposalInfos[id].deposit}(
            proposalInfos[id].proposer,
            proposalInfos[id].deposit,
            proposalInfos[id].rate
        );
        proposalInfos[id].status = ProposalStatus.pass;
        emit LogGuarantee(id, msg.sender, block.number);
    }

    /**
     * @dev updateProposal
     */
    function updateProposal(
        bytes4 id,
        uint8 rate,
        uint256 deposit,
        string memory details
    ) external payable onlyEffictiveProposal(id) {
        require(proposalInfos[id].initBlock != 0, "Proposals: proposal not exist");
        // check status
        require(
            proposalInfos[id].status == ProposalStatus.pending,
            "Proposals: The status of proposal must be pending"
        );
        // check msg.value 、rate、details
        require(bytes(details).length <= MAX_PROPOSAL_DETAIL_LENGTH, "Proposals: details is too long");
        require(deposit >= MIN_DEPOSIT, "Proposals: deposit must greater than MIN_DEPOSIT");
        require(
            rate >= MIN_RATE && rate <= MAX_RATE,
            "Proposals: rate must greater than MIN_RATE and less than MAX_RATE"
        );

        uint256 lastDeposit = proposalInfos[id].deposit;
        if (lastDeposit > deposit) {
            address payable receiver = payable(address(msg.sender));
            receiver.transfer(lastDeposit - deposit);
        } else if (lastDeposit < deposit) {
            require(deposit - lastDeposit == msg.value, "Proposals: msg value not true");
        } else {
            if (msg.value != 0) {
                address payable receiver = payable(address(msg.sender));
                receiver.transfer(msg.value);
            }
        }
        proposalInfos[id].deposit = deposit;
        proposalInfos[id].rate = rate;
        proposalInfos[id].updateBlock = block.number;
        proposalInfos[id].details = details;
        emit LogUpdateProposal(id, msg.sender, block.number, deposit, rate);
    }

    /**
     * @dev cancelProposal
     */
    function cancelProposal(bytes4 id) external onlyEffictiveProposal(id) {
        require(proposalInfos[id].initBlock != 0, "Proposals: proposal not exist");
        // check status
        require(
            proposalInfos[id].status == ProposalStatus.pending,
            "Proposals: The status of proposal must be pending"
        );

        proposalInfos[id].updateBlock = block.number;
        address payable receiver = payable(address(msg.sender));
        receiver.transfer(proposalInfos[id].deposit);
        proposalInfos[id].status = ProposalStatus.cancel;
        emit LogCancelProposal(id, msg.sender, block.number);
    }

    function allProposals(uint256 page, uint256 size) public view returns (ProposalInfo[] memory) {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        uint256 start = (page - 1) * size;
        if (proposalsBytes.length() < start) {
            size = 0;
        } else {
            uint256 length = proposalsBytes.length() - start;
            if (length < size) {
                size = length;
            }
        }

        ProposalInfo[] memory proposalDir = new ProposalInfo[](size);
        for (uint256 i = 0; i < size; i++) {
            proposalDir[i] = proposalInfos[bytes4(proposalsBytes.at(i + start))];
        }
        return proposalDir;
    }

    function allProposalSets(uint256 page, uint256 size) public view returns (bytes4[] memory) {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        uint256 start = (page - 1) * size;
        if (proposalsBytes.length() < start) {
            size = 0;
        } else {
            uint256 length = proposalsBytes.length() - start;
            if (length < size) {
                size = length;
            }
        }
        bytes4[] memory proposalDir = new bytes4[](size);
        for (uint256 i = 0; i < size; i++) {
            proposalDir[i] = bytes4(proposalsBytes.at(i + start));
        }
        return proposalDir;
    }

    function addressProposals(
        address val,
        uint256 page,
        uint256 size
    ) public view returns (ProposalInfo[] memory) {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        bytes4[] memory addressProposalIds = proposals[val];
        uint256 start = (page - 1) * size;
        if (addressProposalIds.length < start) {
            size = 0;
        } else {
            uint256 length = addressProposalIds.length - start;
            if (length < size) {
                size = length;
            }
        }

        ProposalInfo[] memory proposalDir = new ProposalInfo[](size);
        for (uint256 i = 0; i < size; i++) {
            proposalDir[i] = proposalInfos[addressProposalIds[i + start]];
        }
        return proposalDir;
    }

    function addressProposalSets(
        address val,
        uint256 page,
        uint256 size
    ) public view returns (bytes4[] memory) {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        bytes4[] memory addressProposalIds = proposals[val];
        uint256 start = (page - 1) * size;
        if (addressProposalIds.length < start) {
            size = 0;
        } else {
            uint256 length = addressProposalIds.length - start;
            if (length < size) {
                size = length;
            }
        }

        bytes4[] memory proposalDir = new bytes4[](size);
        for (uint256 i = 0; i < size; i++) {
            proposalDir[i] = addressProposalIds[i + start];
        }
        return proposalDir;
    }

    function proposalCount() public view returns (uint256) {
        return proposalsBytes.length();
    }

    function addressProposalCount(address val) public view returns (uint256) {
        return proposals[val].length;
    }
}
