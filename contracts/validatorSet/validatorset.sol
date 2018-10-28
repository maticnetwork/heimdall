pragma solidity ^0.4.24;

import { SafeMath } from "./SafeMath.sol";
import {RLP} from "./rlplib.sol";
import { RLPEncode } from "./rlpencode.sol";
import { BytesLib } from "./byteslib.sol";
import "./ecverify.sol";
contract ValidatorSet is ECVerify {
  using SafeMath for uint256;
  using SafeMath for uint8;
//   using BytesLib for bytes32;
  int256 constant INT256_MIN = -int256((2**255)-1);
  uint256 constant UINT256_MAX = (2**256)-1;
  using RLP for bytes;
  using RLP for RLP.RLPItem;
  using RLP for RLP.Iterator;
  event NewProposer(address indexed user, bytes data);

  struct Validator {
    uint256 votingPower;
    address validator;
    string pubkey;
  }

  address public proposer;
  uint256 public totalVotingPower;
  uint256 public lowestPower;
  Validator[] public validators;

  function addValidator(address validator, uint256 votingPower, string _pubkey) public {
    validators.push(Validator(votingPower, validator,_pubkey)); //use index instead
  }
  function getPubkey(uint256 index) public view returns(string){
    //  return BytesLib.concat(abi.encodePacked(validators[index].pubkey1), abi.encodePacked(validators[index].pubkey2));
    return validators[index].pubkey;
  }
  function removeValidator(uint256 _index){
    delete validators[_index];
  }

    function getValidatorSet() public view returns (uint256[] ,address[]){
        uint256[] memory powers= new uint256[](validators.length);
        address[] memory validatorAddresses= new address[](validators.length);
        for (uint8 i = 0; i < validators.length; i++) {
            validatorAddresses[i]=validators[i].validator;
            powers[i]=validators[i].votingPower;
        }
        return (powers,validatorAddresses);
    }
    // Inputs: start,end,roothash,vote bytes,signatures of validators ,tx(extradata)
    // Constants: chainid,type,votetype
    // extradata => start,end ,proposer etc rlp encoded hash
    //
    //
    // todo : check proposer verify signatures  for validators
    bytes public chainID = "test-chain-E5igIA";
    bytes public roundType = "vote";
    bytes public voteType = "0x02";

    // @params start-> startBlock , end-> EndBlock , roothash-> merkelRoot
    function validate(bytes vote,bytes sigs,bytes extradata)public returns(address,address,uint) {

        RLP.RLPItem[] memory dataList = vote.toRLPItem().toList();
        require(keccak256(dataList[0].toData())==keccak256(chainID),"Chain ID not same");
        require(keccak256(dataList[1].toData())==keccak256(roundType),"Round Type Not Same ");

        // require(keccak256(dataList[5].toData())==keccak256(_voteType),"Vote Type Not Same");

        // validate extra data using getSha256(extradata)
        require(keccak256(dataList[6].toData())==keccak256(getSha256(extradata)));
        // check proposer
        // require(msg.sender==dataList[5].toAddress());
        // decode extra data and validate start,end etc
        RLP.RLPItem[] memory txDataList;
        txDataList=extradata.toRLPItem().toList()[0].toList();

        // extract end and assign to current child
        // require(txDataList[2].toUint() == end, "End Block Does Not Match");

        //slice sigs and do ecrecover from validator set
        for (uint64 i = 0; i < sigs.length; i += 65) {
          bytes memory sigElement = BytesLib.slice(sigs, i, 65);
          address signer = ECVerify.ecrecoveryFromData(vote,sigElement);
        }
        // signer address , proposer address , end block
        return (signer,dataList[5].toAddress(),txDataList[2].toUint());

    }
    function setChainId(string _chainID) public {
        chainID = bytes(_chainID);
    }
    function getSha256(bytes input) public returns (bytes20) {
        bytes32 hash = sha256(input);
        return bytes20(hash);

    }
}