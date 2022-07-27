// SPDX-License-Identifier: GPL-v3
pragma solidity ^0.8.0;

import "../test-helper/TestBase.sol";

contract TestChainStorageContainer is TestBase, ChainStorageContainer {
    address sender = address(7777);
    ChainStorageContainer chainStorageContainer;

    function setUp() public {
        _initialize();
        vm.startPrank(sender);
        ChainStorageContainer chainStorageContainerLogic = new ChainStorageContainer();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(chainStorageContainerLogic),
            address(proxyAdmin),
            abi.encodeWithSelector(
                chainStorageContainerLogic.initialize.selector,
                AddressName.L1_CROSS_LAYER_WITNESS,
                address(addressManager)
            )
        );
        chainStorageContainer = ChainStorageContainer(address(proxy));
        vm.stopPrank();
    }

    function testChainSize() public view {
        require(chainStorageContainer.chainSize() == 0);
    }

    function testAppend() public {
        vm.startPrank(address(l1CrossLayerWitness));
        bytes32 ele = 0x0;
        uint64 chainSize = chainStorageContainer.append(ele);
        require(chainSize == 1);
        require(chainStorageContainer.chainSize() == 1);
    }

    function testAppendWithCallerIsNotOwner() public {
        vm.startPrank(address(0x88));
        bytes32 ele = 0x0;
        vm.expectRevert("ChainStorageContainer: Function can only be called by the owner.");
        chainStorageContainer.append(ele);
    }

    function testResize() public {
        vm.startPrank(address(l1CrossLayerWitness));
        bytes32 ele = 0x0;
        chainStorageContainer.append(ele);
        chainStorageContainer.append(ele);
        uint64 chainSize = chainStorageContainer.append(ele);
        require(chainSize == 3);
        require(chainStorageContainer.chainSize() == 3);
        chainStorageContainer.resize(uint64(2));
        require(chainStorageContainer.chainSize() == 2);

        bytes32 ele1 = bytes32(uint256(5));
        chainStorageContainer.append(ele1);
        ele = chainStorageContainer.get(2);
        require(ele == ele1);
    }

    function testResizeWithFailed() public {
        vm.startPrank(address(l1CrossLayerWitness));
        bytes32 ele = 0x0;
        chainStorageContainer.append(ele);
        chainStorageContainer.append(ele);
        uint64 chainSize = chainStorageContainer.append(ele);
        require(chainSize == 3);
        require(chainStorageContainer.chainSize() == 3);
        vm.expectRevert("can't resize beyond chain length");
        chainStorageContainer.resize(uint64(4));
    }

    function testResizeWithCallerIsNotOwner() public {
        vm.startPrank(address(0x88));
        vm.expectRevert("ChainStorageContainer: Function can only be called by the owner.");
        chainStorageContainer.resize(0);
    }

    function testGet() public {
        vm.startPrank(address(l1CrossLayerWitness));
        bytes32 ele = 0x0;
        chainStorageContainer.append(ele);
        chainStorageContainer.append(ele);
        uint64 chainSize = chainStorageContainer.append(ele);
        require(chainSize == 3);
        require(chainStorageContainer.chainSize() == 3);
        require(chainStorageContainer.get(0) == ele);
        require(chainStorageContainer.get(1) == ele);
        require(chainStorageContainer.get(2) == ele);
        vm.expectRevert("beyond chain size");
        chainStorageContainer.get(3);
    }
}
