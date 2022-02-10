// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./MerkleTrie.sol";
import "./console.sol";


#define TEST_CASE( testnum, testreg, correctval, code... ) \
test_ ## testnum: \
code; \
li  x7, MASK_XLEN(correctval); \
li  TESTNUM, testnum; \
bne testreg, x7, fail;

