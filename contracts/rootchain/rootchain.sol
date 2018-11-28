pragma solidity ^0.4.24;

import { SafeMath } from "./SafeMath.sol";

import { Merkle } from "./Merkle.sol";
import { RLP } from "./RLP.sol";
import { RLPEncode } from "./RLPEncode.sol";

import { IRootChain } from "./IRootChain.sol";
import { StakeManager } from "./StakeManager.sol";


contract IManager {
  // chain identifier
  bytes32 public chain = keccak256("heimdall-A7lVlP");
  // round type
  bytes32 public constant roundType = keccak256("vote");
  // vote type
  byte public constant voteType = 0x02;
  // network id
  bytes public constant networkId = "\x0d";
  // child block interval between checkpoint
  uint256 public constant CHILD_BLOCK_INTERVAL = 10000;
}

contract RootChain is IRootChain, IManager {
  using SafeMath for uint256;
  using Merkle for bytes32;
  using RLP for bytes;
  using RLP for RLP.RLPItem;
  using RLP for RLP.Iterator;

  // child chain contract
  address public childChainContract;

  // list of header blocks (address => header block object)
  mapping(uint256 => HeaderBlock) public headerBlocks;

  // current header block number
  uint256 private _currentHeaderBlock;


  // stake interface
  StakeManager public stakeManager;

  //
  // Constructor
  //

  constructor (address _stakeManager) public {
    setStakeManager(_stakeManager);
    _currentHeaderBlock = CHILD_BLOCK_INTERVAL;
  }

  //
  // Events
  //

  event NewHeaderBlock(
    address indexed proposer,
    uint256 indexed number,
    uint256 start,
    uint256 end,
    bytes32 root
  );

  //
  // External functions
  //

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

    // Make sure enough validators sign off on the proposed header root
    require(stakeManager.checkSignatures(keccak256(vote), sigs));

    // Add the header root
    HeaderBlock memory headerBlock = HeaderBlock({
      root: root,
      start: start,
      end: end,
      createdAt: block.timestamp,
      proposer: msg.sender
    });
    headerBlocks[_currentHeaderBlock] = headerBlock;

    // emit new header block
    emit NewHeaderBlock(
      msg.sender,
      _currentHeaderBlock,
      headerBlock.start,
      headerBlock.end,
      root
    );

    // update current header block
    _currentHeaderBlock = _currentHeaderBlock.add(CHILD_BLOCK_INTERVAL);

    // finalize commit
    stakeManager.finalizeCommit();

    // TODO add rewards
  }

  function setChain(string c) public {
    chain = keccak256(c);
  }

  function currentChildBlock() public view returns(uint256) {
    if (_currentHeaderBlock != CHILD_BLOCK_INTERVAL) {
      return headerBlocks[_currentHeaderBlock.sub(CHILD_BLOCK_INTERVAL)].end;
    }

    return 0;
  }

  function currentHeaderBlock() public view returns (uint256) {
    return _currentHeaderBlock;
  }

  function setFakeHeaderBlock(uint256 start, uint256 end) public {
      // Add the header root
    HeaderBlock memory headerBlock = HeaderBlock({
      root: keccak256(abi.encodePacked(start, end)),
      start: start,
      end: end,
      createdAt: block.timestamp,
      proposer: msg.sender
    });
    headerBlocks[_currentHeaderBlock] = headerBlock;

    // update current header block
    _currentHeaderBlock = _currentHeaderBlock.add(CHILD_BLOCK_INTERVAL);
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

  // get flat deposit block
  function depositBlock(uint256 _depositCount)
      public
      view
      returns
    (
      uint256 _header,
      address _owner,
      address _token,
      uint256 _amount,
      uint256 _createdAt
    ) {}

  // set stake manager
  function setStakeManager(address _stakeManager) public {
    require(_stakeManager != address(0));
    stakeManager = StakeManager(_stakeManager);
  }

  // finalize commit
  function finalizeCommit(uint256) public {}

  // slash stakers if fraud is detected
  function slash() public {
    // TODO pass block/proposer
  }

  function transferAmount(
    address _token,
    address _user,
    uint256 _amount,
    bool isWeth
  ) public returns(bool) {
      return true;
  }
}