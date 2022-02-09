// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./RLPReader.sol";
import "./RLPWriter.sol";
import "./BytesSlice.sol";

contract RLPCodecTest {
    function testReadWriteUint(uint256 val) pure public {
        bytes memory rlp = RLPWriter.writeUint(val);
        uint256 decoded = RLPReader.readUint256(rlp);
        require(decoded == val);
    }

    function checkWriteUint(uint256 val, bytes memory encoded) pure internal {
        require(BytesSlice.equal(RLPWriter.writeUint(val), encoded));
        uint256 decoded = RLPReader.readUint256(encoded);
        require(decoded == val);
    }

    function checkBool(bool val, bytes memory encoded) pure internal {
        require(BytesSlice.equal(RLPWriter.writeBool(val), encoded));
        bool decoded = RLPReader.readBool(encoded);
        require(decoded == val);
    }

    function testWrite() pure public {
        checkWriteUint(0, hex"80");
        checkWriteUint(1, hex"01");
        checkWriteUint(0, hex"80");
        checkWriteUint(0, hex"80");
        checkWriteUint(1, hex"01");
        checkWriteUint(127, hex"7F");
        checkWriteUint(128, hex"8180");
        checkWriteUint(256, hex"820100");
        checkWriteUint(1024, hex"820400");
        checkWriteUint(0xFFFFFF, hex"83FFFFFF");
        checkWriteUint(0xFFFFFFFF, hex"84FFFFFFFF");
        checkWriteUint(0xFFFFFFFF, hex"84FFFFFFFF");
        checkWriteUint(0xFFFFFFFFFF, hex"85FFFFFFFFFF");
        checkWriteUint(0xFFFFFFFFFFFF, hex"86FFFFFFFFFFFF");
        checkWriteUint(0xFFFFFFFFFFFFFF, hex"87FFFFFFFFFFFFFF");
        checkWriteUint(0xFFFFFFFFFFFFFFFF, hex"88FFFFFFFFFFFFFFFF");

        checkBool(true, hex"01");
        checkBool(false, hex"80");

        checkWriteBytes("", hex"80");
        checkWriteBytes(hex"00", hex"00");
        checkWriteBytes(hex"01", hex"01");
        checkWriteBytes(hex"7f", hex"7f");
        checkWriteBytes(hex"80", hex"8180");
        checkWriteBytes(hex"ff", hex"81ff");
        checkWriteBytes(hex"010203", hex"83010203");
        checkWriteBytes("dog", hex"83646F67");
        checkWriteBytes(
            "Lorem ipsum dolor sit amet, consectetur adipisicing eli",
            hex"B74C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C69"
        );
        checkWriteBytes(
            "Lorem ipsum dolor sit amet, consectetur adipisicing elit",
            hex"B8384C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C6974"
        );
        checkWriteBytes(
            "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat",
            hex"B904004C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E73656374657475722061646970697363696E6720656C69742E20437572616269747572206D6175726973206D61676E612C20737573636970697420736564207665686963756C61206E6F6E2C20696163756C697320666175636962757320746F72746F722E2050726F696E20737573636970697420756C74726963696573206D616C6573756164612E204475697320746F72746F7220656C69742C2064696374756D2071756973207472697374697175652065752C20756C7472696365732061742072697375732E204D6F72626920612065737420696D70657264696574206D6920756C6C616D636F7270657220616C6971756574207375736369706974206E6563206C6F72656D2E2041656E65616E2071756973206C656F206D6F6C6C69732C2076756C70757461746520656C6974207661726975732C20636F6E73657175617420656E696D2E204E756C6C6120756C74726963657320747572706973206A7573746F2C20657420706F73756572652075726E6120636F6E7365637465747572206E65632E2050726F696E206E6F6E20636F6E76616C6C6973206D657475732E20446F6E65632074656D706F7220697073756D20696E206D617572697320636F6E67756520736F6C6C696369747564696E2E20566573746962756C756D20616E746520697073756D207072696D697320696E206661756369627573206F726369206C756374757320657420756C74726963657320706F737565726520637562696C69612043757261653B2053757370656E646973736520636F6E76616C6C69732073656D2076656C206D617373612066617563696275732C2065676574206C6163696E6961206C616375732074656D706F722E204E756C6C61207175697320756C747269636965732070757275732E2050726F696E20617563746F722072686F6E637573206E69626820636F6E64696D656E74756D206D6F6C6C69732E20416C697175616D20636F6E73657175617420656E696D206174206D65747573206C75637475732C206120656C656966656E6420707572757320656765737461732E20437572616269747572206174206E696268206D657475732E204E616D20626962656E64756D2C206E6571756520617420617563746F72207472697374697175652C206C6F72656D206C696265726F20616C697175657420617263752C206E6F6E20696E74657264756D2074656C6C7573206C65637475732073697420616D65742065726F732E20437261732072686F6E6375732C206D65747573206163206F726E617265206375727375732C20646F6C6F72206A7573746F20756C747269636573206D657475732C20617420756C6C616D636F7270657220766F6C7574706174"
        );
    }

    function testWriteList() pure public {
        bytes[] memory list = new bytes[](15);
        list[0] = "aaa";
        list[1] = "bbb";
        list[2] = "ccc";
        list[3] = "ddd";
        list[4] = "eee";
        list[5] = "fff";
        list[6] = "ggg";
        list[7] = "hhh";
        list[8] = "iii";
        list[9] = "jjj";
        list[10] = "kkk";
        list[11] = "lll";
        list[12] = "mmm";
        list[13] = "nnn";
        list[14] = "ooo";
        bytes[] memory encodedList = new bytes[](15);
        for (uint256 i = 0; i < list.length; i++) {
            encodedList[i] = RLPWriter.writeBytes(list[i]);
        }
        require(
            BytesSlice.equal(
                RLPWriter.writeList(encodedList),
                hex"F83C836161618362626283636363836464648365656583666666836767678368686883696969836A6A6A836B6B6B836C6C6C836D6D6D836E6E6E836F6F6F"
            )
        );
    }

    function checkWriteBytes(bytes memory val, bytes memory encoded) pure internal {
        require(BytesSlice.equal(RLPWriter.writeBytes(val), encoded));
        bytes memory decoded = RLPReader.readBytes(encoded);
        require(BytesSlice.equal(decoded, val));
    }

    //the normal rlp code is: hex"F83C836161618362626283636363836464648365656583666666836767678368686883696969836A6A6A836B6B6B836C6C6C836D6D6D836E6E6E836F6F6F"
    function testFailReadInconsistentRlp() public {
        RLPReader.readList(
            hex"F83C836161618362626283636363836464648365656583666666836767678368686883696969836A6A6A836B6B6B836C6C6C836D6D6D836E6E6E836F6F6F66"
        );
    }
}
