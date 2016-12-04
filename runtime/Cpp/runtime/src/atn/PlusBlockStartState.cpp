/* Copyright (c) 2012 The ANTLR Project Authors. All rights reserved.
 * Use of this file is governed by the BSD 3-clause license that
 * can be found in the LICENSE.txt file in the project root.
 */

#include "atn/PlusBlockStartState.h"

using namespace antlr4::atn;

size_t PlusBlockStartState::getStateType() {
  return PLUS_BLOCK_START;
}
