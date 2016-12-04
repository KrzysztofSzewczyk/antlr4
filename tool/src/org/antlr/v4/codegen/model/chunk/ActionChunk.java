/*
 * Copyright (c) 2012 The ANTLR Project Authors. All rights reserved.
 * Use of this file is governed by the BSD 3-clause license that
 * can be found in the LICENSE.txt file in the project root.
 */

package org.antlr.v4.codegen.model.chunk;

import org.antlr.v4.codegen.model.OutputModelObject;
import org.antlr.v4.codegen.model.decl.StructDecl;

/** */
public class ActionChunk extends OutputModelObject {
	/** Where is the ctx that defines attrs,labels etc... for this action? */
	public StructDecl ctx;

	public ActionChunk(StructDecl ctx) {
		this.ctx = ctx;
	}
}
