pragma solidity ^0.4.24;


import { StakeManager } from "./validatorSet.sol";

import { SafeMath } from "./SafeMath.sol";
import {RLP} from "./rlplib.sol";
import { RLPEncode } from "./rlpencode.sol";
import { BytesLib } from "./byteslib.sol";

contract RootMock {
    using RLP for bytes;
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using SafeMath for uint256;

    constructor (address _stakeManager) public {
    setStakeManager(_stakeManager);
    // set current header block
    _currentHeaderBlock = CHILD_BLOCK_INTERVAL;
    }
    uint256 public constant CHILD_BLOCK_INTERVAL = 10000;
    StakeManager public stakeManager;

    bytes32 public constant chain = keccak256("heimdall-40J7bp");
      // round type
    bytes32 public constant roundType = keccak256("vote");
      // vote type
    byte public constant voteType = 0x02;

    struct HeaderBlock {
     bytes32 root;
     uint256 start;
     uint256 end;
     uint256 createdAt;
     address proposer;
   }
    uint256 private _currentHeaderBlock;

    mapping(uint256 => HeaderBlock) public headerBlocks;

    function submitHeaderBlock(bytes vote, bytes sigs, bytes extradata) external {

        RLP.RLPItem[] memory dataList = vote.toRLPItem().toList();
        require(keccak256(dataList[0].toData()) == chain, "Chain ID not same");
        require(keccak256(dataList[1].toData()) == roundType, "Round type not same ");
        require(dataList[4].toByte() == voteType, "Vote type not same");

        // check proposer
        require(msg.sender == dataList[5].toAddress());

        // validate extra data using getSha256(extradata)
        require(keccak256(dataList[6].toData()) == keccak256(bytes20(sha256(extradata))));

        // extract end and assign to current child
        dataList = extradata.toRLPItem().toList()[0].toList();
        uint256 start = currentChildBlock();
        uint256 end = dataList[2].toUint();
        bytes32 root = dataList[3].toBytes32();

        if (start > 0) {
          start = start.add(1);
        }

        // Start on mainchain and matic chain must be same
        require(start == dataList[1].toUint());

        // Make sure we are adding blocks
        require(end > start);

        // Add the header root
        HeaderBlock memory headerBlock = HeaderBlock({
          root: root,
          start: start,
          end: end,
          createdAt: block.timestamp,
          proposer: msg.sender
        });

        headerBlocks[_currentHeaderBlock] = headerBlock;
        _currentHeaderBlock = _currentHeaderBlock.add(CHILD_BLOCK_INTERVAL);

  }
  function setStakeManager(address _stakeManager) public  {
    require(_stakeManager != address(0));
    stakeManager = StakeManager(_stakeManager);
  }

    function headerBlock(uint256 _headerNumber) public view returns (
    bytes32 _root,
    uint256 _start,
    uint256 _end,
    uint256 _createdAt
    ) {
    HeaderBlock memory _headerBlock = headerBlocks[_headerNumber];

    _root = _headerBlock.root;
    _start = _headerBlock.start;
    _end = _headerBlock.end;
    _createdAt = _headerBlock.createdAt;
  }
  function currentChildBlock() public view returns(uint256) {
    if (_currentHeaderBlock != CHILD_BLOCK_INTERVAL) {
      return headerBlocks[_currentHeaderBlock.sub(CHILD_BLOCK_INTERVAL)].end;
    }

    return 0;
  }

}