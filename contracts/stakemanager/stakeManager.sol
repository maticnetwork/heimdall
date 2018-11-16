pragma solidity ^0.4.24;

import { SafeMath } from "./SafeMath.sol";
import {RLP} from "./rlplib.sol";
import { RLPEncode } from "./rlpencode.sol";
import { BytesLib } from "./byteslib.sol";
import "./ecverify.sol";
contract StakeManager is ECVerify {
  using SafeMath for uint8,uint256;
  using RLP for bytes;
  using RLP for RLP.RLPItem;
  using RLP for RLP.Iterator;

  struct Validator {
    uint256 votingPower;
    address validator;
  }

  Validator[] public validators;

  function addValidator(address validator, uint256 votingPower) public {
    validators.push(Validator(votingPower, validator)); //use index instead
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



}