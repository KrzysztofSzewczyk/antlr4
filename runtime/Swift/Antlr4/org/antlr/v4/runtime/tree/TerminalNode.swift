/* Copyright (c) 2012 The ANTLR Project Contributors. All rights reserved.
 * Use of this file is governed by the BSD 3-clause license that
 * can be found in the LICENSE.txt file in the project root.
 */

public class TerminalNode: ParseTree {
    public func getSymbol() -> Token? {
        RuntimeException(" must overriden !")
        fatalError()

    }
}
