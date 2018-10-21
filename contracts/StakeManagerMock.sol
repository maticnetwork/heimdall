pragma solidity ^0.4.24;

contract StakeManager{
    struct Validator{
        string pubkey;
        int power;
    }
    // mapping(uint256=>Validator) public validators;
    uint256  public lastValidatorIndex=0;
    Validator[] public validators;
    function addValidator(string memory _pubKey ,int _power) public {
        uint256 last = validators.push(Validator(_pubKey,_power));
        lastValidatorIndex=last;
    }

}