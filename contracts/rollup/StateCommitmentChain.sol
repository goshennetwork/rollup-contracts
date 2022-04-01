// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/* Library Imports */
import { OVMCodec } from "../libraries/OVMCodec.sol";
import { AddressResolver } from "../libraries/AddressResolver.sol";
import { MerkleTree } from "../libraries/MerkleTree.sol";

/* Interface Imports */
import { IStateCommitChain } from "../interfaces/IStateCommitChain.sol";
import { ICanonicalTransactionChain } from "../interfaces/ICanonicalTransactionChain.sol";
import { IStakingManager } from "../interfaces/IStakingManager.sol";
import { IChainStorageContainer } from "../interfaces/IChainStorageContainer.sol";

/**
 * @title StateCommitmentChain
 * @dev The State Commitment Chain (SCC) contract contains a list of proposed state roots which
 * Proposers assert to be a result of each transaction in the Canonical Transaction Chain (CTC).
 * Elements here have a 1:1 correspondence with transactions in the CTC, and should be the unique
 * state root calculated off-chain by applying the canonical transactions one by one.
 *
 */
contract StateCommitmentChain is IStateCommitChain, AddressResolver {
    /*************
     * Constants *
     *************/

    uint256 public FRAUD_PROOF_WINDOW;
    uint256 public SEQUENCER_PUBLISH_WINDOW;

    /***************
     * Constructor *
     ***************/

    /**
     * @param _libAddressManager Address of the Address Manager.
     */
    constructor(
        address _libAddressManager,
        uint256 _fraudProofWindow,
        uint256 _sequencerPublishWindow
    ) AddressResolver(_libAddressManager) {
        FRAUD_PROOF_WINDOW = _fraudProofWindow;
        SEQUENCER_PUBLISH_WINDOW = _sequencerPublishWindow;
    }

    /********************
     * Public Functions *
     ********************/

    /**
     * Accesses the batch storage container.
     * @return Reference to the batch storage container.
     */
    function batches() public view returns (IChainStorageContainer) {
        return IChainStorageContainer(resolve("ChainStorageContainer-SCC-batches"));
    }

    /**
     * @inheritdoc IStateCommitChain
     */
    function getCurrentBlockHeight() public view returns (uint256 _totalElements) {
        (uint40 totalElements, ) = _getBatchExtraData();
        return uint256(totalElements);
    }

    /**
     * @inheritdoc IStateCommitChain
     */
    function getTotalBatches() public view returns (uint256 _totalBatches) {
        return batches().length();
    }

    /**
     * @inheritdoc IStateCommitChain
     */
    // slither-disable-next-line external-function
    function getLastSequencerTimestamp() public view returns (uint256 _lastSequencerTimestamp) {
        (, uint40 lastSequencerTimestamp) = _getBatchExtraData();
        return uint256(lastSequencerTimestamp);
    }

    /**
     * @inheritdoc IStateCommitChain
     */
    // slither-disable-next-line external-function
    function appendStateBatch(bytes32[] memory _batch, uint256 _shouldStartAtElement) public {
        // Fail fast in to make sure our batch roots aren't accidentally made fraudulent by the
        // publication of batches by some other user.
        require(
            _shouldStartAtElement == getCurrentBlockHeight(),
            "Actual batch start index does not match expected start index."
        );

        // Proposers must have previously staked at the BondManager
        require(IStakingManager(resolve("StakingManager")).isStaking(msg.sender), "Proposer should be staking");

        require(_batch.length > 0, "Cannot submit an empty state batch.");

        require(
            getCurrentBlockHeight() + _batch.length <=
                ICanonicalTransactionChain(resolve("CanonicalTransactionChain")).getTotalElements(),
            "Number of state roots cannot exceed the number of canonical transactions."
        );

        // Pass the block's timestamp and the publisher of the data
        // to be used in the fraud proofs
        _appendBatch(_batch, abi.encode(block.timestamp, msg.sender));
    }

    /**
     * @inheritdoc IStateCommitChain
     */
    // slither-disable-next-line external-function
    function deleteStateBatch(OVMCodec.ChainBatchHeader memory _batchHeader) public {
        require(
            msg.sender == resolve("OVM_FraudVerifier"),
            "State batches can only be deleted by the OVM_FraudVerifier."
        );

        require(_isValidBatchHeader(_batchHeader), "Invalid batch header.");

        require(
            insideFraudProofWindow(_batchHeader),
            "State batches can only be deleted within the fraud proof window."
        );

        _deleteBatch(_batchHeader);
    }

    function verifyBatchHeader(OVMCodec.ChainBatchHeader memory _batchHeader) public view returns (bool) {
        return _isValidBatchHeader(_batchHeader);
    }

    /**
     * @inheritdoc IStateCommitChain
     */
    // slither-disable-next-line external-function
    function verifyStateCommitment(
        OVMCodec.BlockInfo memory _element,
        OVMCodec.ChainBatchHeader memory _batchHeader,
        OVMCodec.ChainInclusionProof memory _proof
    ) public view returns (bool) {
        require(_isValidBatchHeader(_batchHeader), "Invalid batch header.");

        require(
            MerkleTree.verify(
                _batchHeader.batchRoot,
                OVMCodec.hashBlockInfo(_element),
                _proof.index,
                _proof.siblings,
                _batchHeader.batchSize
            ),
            "Invalid inclusion proof."
        );

        return true;
    }

    /**
     * @inheritdoc IStateCommitChain
     */
    function insideFraudProofWindow(OVMCodec.ChainBatchHeader memory _batchHeader) public view returns (bool _inside) {
        (uint256 timestamp, ) = abi.decode(_batchHeader.extraData, (uint256, address));

        require(timestamp != 0, "Batch header timestamp cannot be zero");
        return (timestamp + FRAUD_PROOF_WINDOW) > block.timestamp;
    }

    /**********************
     * Internal Functions *
     **********************/

    /**
     * Parses the batch context from the extra data.
     * @return Total number of elements submitted.
     * @return Timestamp of the last batch submitted by the sequencer.
     */
    function _getBatchExtraData() internal view returns (uint40, uint40) {
        bytes27 extraData = batches().getGlobalMetadata();

        // solhint-disable max-line-length
        uint40 totalElements;
        uint40 lastSequencerTimestamp;
        assembly {
            extraData := shr(40, extraData)
            totalElements := and(extraData, 0x000000000000000000000000000000000000000000000000000000FFFFFFFFFF)
            lastSequencerTimestamp := shr(
                40,
                and(extraData, 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0000000000)
            )
        }
        // solhint-enable max-line-length

        return (totalElements, lastSequencerTimestamp);
    }

    /**
     * Encodes the batch context for the extra data.
     * @param _totalElements Total number of elements submitted.
     * @param _lastSequencerTimestamp Timestamp of the last batch submitted by the sequencer.
     * @return Encoded batch context.
     */
    function _makeBatchExtraData(uint40 _totalElements, uint40 _lastSequencerTimestamp)
        internal
        pure
        returns (bytes27)
    {
        bytes27 extraData;
        assembly {
            extraData := _totalElements
            extraData := or(extraData, shl(40, _lastSequencerTimestamp))
            extraData := shl(40, extraData)
        }

        return extraData;
    }

    /**
     * Appends a batch to the chain.
     * @param _batch Elements within the batch.
     * @param _extraData Any extra data to append to the batch.
     */
    function _appendBatch(bytes32[] memory _batch, bytes memory _extraData) internal {
        address sequencer = resolve("OVM_Proposer");
        (uint40 totalElements, uint40 lastSequencerTimestamp) = _getBatchExtraData();

        if (msg.sender == sequencer) {
            lastSequencerTimestamp = uint40(block.timestamp);
        } else {
            // We keep track of the last batch submitted by the sequencer so there's a window in
            // which only the sequencer can publish state roots. A window like this just reduces
            // the chance of "system breaking" state roots being published while we're still in
            // testing mode. This window should be removed or significantly reduced in the future.
            require(
                lastSequencerTimestamp + SEQUENCER_PUBLISH_WINDOW < block.timestamp,
                "Cannot publish state roots within the sequencer publication window."
            );
        }

        // For efficiency reasons getMerkleRoot modifies the `_batch` argument in place
        // while calculating the root hash therefore any arguments passed to it must not
        // be used again afterwards
        OVMCodec.ChainBatchHeader memory batchHeader = OVMCodec.ChainBatchHeader({
            batchIndex: getTotalBatches(),
            batchRoot: MerkleTree.getMerkleRoot(_batch),
            batchSize: _batch.length,
            prevTotalElements: totalElements,
            extraData: _extraData
        });

        emit StateBatchAppended(
            batchHeader.batchIndex,
            batchHeader.batchRoot,
            batchHeader.batchSize,
            batchHeader.prevTotalElements,
            batchHeader.extraData
        );

        batches().push(
            OVMCodec.hashBatchHeader(batchHeader),
            _makeBatchExtraData(uint40(batchHeader.prevTotalElements + batchHeader.batchSize), lastSequencerTimestamp)
        );
    }

    /**
     * Removes a batch and all subsequent batches from the chain.
     * @param _batchHeader Header of the batch to remove.
     */
    function _deleteBatch(OVMCodec.ChainBatchHeader memory _batchHeader) internal {
        require(_batchHeader.batchIndex < batches().length(), "Invalid batch index.");

        require(_isValidBatchHeader(_batchHeader), "Invalid batch header.");

        // slither-disable-next-line reentrancy-events
        batches().deleteElementsAfterInclusive(
            _batchHeader.batchIndex,
            _makeBatchExtraData(uint40(_batchHeader.prevTotalElements), 0)
        );

        // slither-disable-next-line reentrancy-events
        emit StateBatchDeleted(_batchHeader.batchIndex, _batchHeader.batchRoot);
    }

    /**
     * Checks that a batch header matches the stored hash for the given index.
     * @param _batchHeader Batch header to validate.
     * @return Whether or not the header matches the stored one.
     */
    function _isValidBatchHeader(OVMCodec.ChainBatchHeader memory _batchHeader) internal view returns (bool) {
        return OVMCodec.hashBatchHeader(_batchHeader) == batches().get(_batchHeader.batchIndex);
    }
}
