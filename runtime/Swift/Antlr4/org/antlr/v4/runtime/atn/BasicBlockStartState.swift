/* Copyright (c) 2012 The ANTLR Project Contributors. All rights reserved.
 * Use of this file is governed by the BSD 3-clause license that
 * can be found in the LICENSE.txt file in the project root.
 */



/**
 *
 * @author Sam Harwell
 */

public final class BasicBlockStartState: BlockStartState {
    override
    public func getStateType() -> Int {
        return BlockStartState.BLOCK_START
    }
}
