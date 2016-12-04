/* Copyright (c) 2012 The ANTLR Project Authors. All rights reserved.
 * Use of this file is governed by the BSD 3-clause license that
 * can be found in the LICENSE.txt file in the project root.
 */

#pragma once

#include "atn/DecisionState.h"

namespace antlr4 {
namespace atn {

  ///  The start of a regular {@code (...)} block.
  class ANTLR4CPP_PUBLIC BlockStartState : public DecisionState {
  public:
    BlockEndState *endState = nullptr;
  };

} // namespace atn
} // namespace antlr4
