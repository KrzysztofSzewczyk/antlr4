package org.antlr.v4.test.runtime.javascript.node;

import org.antlr.v4.test.runtime.BaseRuntimeTest;
import org.antlr.v4.test.runtime.RuntimeTestDescriptor;
import org.antlr.v4.test.runtime.descriptors.PerformanceDescriptors;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;

@RunWith(Parameterized.class)
public class TestPerformance extends BaseRuntimeTest {
	public TestPerformance(RuntimeTestDescriptor descriptor) {
		super(descriptor,new BaseNodeTest());
	}

	@Parameterized.Parameters(name="{0}")
	public static RuntimeTestDescriptor[] getAllTestDescriptors() {
		return BaseRuntimeTest.getRuntimeTestDescriptors(PerformanceDescriptors.class, "Node");
	}
}
