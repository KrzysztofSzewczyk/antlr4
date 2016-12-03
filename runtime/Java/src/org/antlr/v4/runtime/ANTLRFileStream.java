/*
 * Copyright (c) 2012 The ANTLR Project Contributors. All rights reserved.
 * Use is of this file is governed by the BSD 3-clause license that
 * can be found in the LICENSE.txt file in the project root.
 */
package org.antlr.v4.runtime;

import org.antlr.v4.runtime.misc.Utils;

import java.io.IOException;

/**
 * This is an {@link ANTLRInputStream} that is loaded from a file all at once
 * when you construct the object.
 */
public class ANTLRFileStream extends ANTLRInputStream {
	protected String fileName;

	public ANTLRFileStream(String fileName) throws IOException {
		this(fileName, null);
	}

	public ANTLRFileStream(String fileName, String encoding) throws IOException {
		this.fileName = fileName;
		load(fileName, encoding);
	}

	public void load(String fileName, String encoding)
		throws IOException
	{
		data = Utils.readFile(fileName, encoding);
		this.n = data.length;
	}

	@Override
	public String getSourceName() {
		return fileName;
	}
}
