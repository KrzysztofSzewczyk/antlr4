/*
 * Copyright (c) 2012 The ANTLR Project Contributors. All rights reserved.
 * Use is of this file is governed by the BSD 3-clause license that
 * can be found in the LICENSE.txt file in the project root.
 */

package org.antlr.v4.gui;

import org.antlr.v4.runtime.tree.Tree;

public interface TreeTextProvider {
	String getText(Tree node);
}
