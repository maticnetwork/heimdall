pragma solidity ^0.4.24;
contract RootMock {
    bool public headerBockSubmitted=false;
    event HeaderBlock(bytes32 root,uint256 start,uint256 end,bytes sigs);
    function submitHeaderBlock (bytes32 root, uint256 start, uint256 end ,bytes sigs) public returns(bool) {
        emit HeaderBlock(root,start,end,sigs);
        headerBockSubmitted=true;
        return headerBockSubmitted;
    }
}