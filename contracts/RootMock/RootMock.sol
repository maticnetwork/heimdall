pragma solidity ^0.4.24;
contract RootMock {
    event HeaderBlock(bytes32,uint256,uint256,bytes);
    function submitHeaderBlock (bytes32 root, uint256 start, uint256 end ,bytes sigs) public {
        emit HeaderBlock(root,start,end,sigs);
    }
}