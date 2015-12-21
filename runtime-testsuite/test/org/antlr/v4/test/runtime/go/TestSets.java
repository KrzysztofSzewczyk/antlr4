/* This file is generated by TestGenerator, any edits will be overwritten by the next generation. */
package org.antlr.v4.test.runtime.go;

import org.junit.Ignore;
import org.junit.Test;

import static org.junit.Assert.*;

@SuppressWarnings("unused")
public class TestSets extends BaseTest {

	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testCharSetLiteral() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(85);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : (A {console.log($A.text);})+ ;\n");
		grammarBuilder.append("A : [AaBb] ;\n");
		grammarBuilder.append("WS : (' '|'\\n')+ -> skip ;");
		String grammar = grammarBuilder.toString();
		String input ="A a B b";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals(
			"A\n" +
			"a\n" +
			"B\n" +
			"b\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testComplementSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(51);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("parse : ~NEW_LINE;\n");
		grammarBuilder.append("NEW_LINE: '\\r'? '\\n';");
		String grammar = grammarBuilder.toString();
		String input ="a";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "parse", input, false);
		assertEquals("", found);

		assertEquals(
			"line 1:0 token recognition error at: 'a'\n" +
			"line 1:1 missing {} at '<EOF>'\n", this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testLexerOptionalSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(77);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log(this._input.getText());} ;\n");
		grammarBuilder.append("A : ('a'|'b')? 'c' ;");
		String grammar = grammarBuilder.toString();
		String input ="ac";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("ac\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testLexerPlusSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(77);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log(this._input.getText());} ;\n");
		grammarBuilder.append("A : ('a'|'b')+ 'c' ;");
		String grammar = grammarBuilder.toString();
		String input ="abaac";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("abaac\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testLexerStarSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(77);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log(this._input.getText());} ;\n");
		grammarBuilder.append("A : ('a'|'b')* 'c' ;");
		String grammar = grammarBuilder.toString();
		String input ="abaac";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("abaac\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testNotChar() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(53);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log($A.text);} ;\n");
		grammarBuilder.append("A : ~'b' ;");
		String grammar = grammarBuilder.toString();
		String input ="x";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("x\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testNotCharSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(59);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log($A.text);} ;\n");
		grammarBuilder.append("A : ~('b'|'c') ;");
		String grammar = grammarBuilder.toString();
		String input ="x";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("x\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testNotCharSetWithLabel() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(61);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log($A.text);} ;\n");
		grammarBuilder.append("A : h=~('b'|'c') ;");
		String grammar = grammarBuilder.toString();
		String input ="x";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("x\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testNotCharSetWithRuleRef3() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(125);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log($A.text);} ;\n");
		grammarBuilder.append("A : ('a'|B) ;  // this doesn't collapse to set but works\n");
		grammarBuilder.append("fragment\n");
		grammarBuilder.append("B : ~('a'|'c') ;");
		String grammar = grammarBuilder.toString();
		String input ="x";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("x\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testOptionalLexerSingleElement() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(71);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log(this._input.getText());} ;\n");
		grammarBuilder.append("A : 'b'? 'c' ;");
		String grammar = grammarBuilder.toString();
		String input ="bc";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("bc\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testOptionalSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(69);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : ('a'|'b')? 'c' {console.log(this._input.getText());} ;");
		String grammar = grammarBuilder.toString();
		String input ="ac";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("ac\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testOptionalSingleElement() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(71);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A? 'c' {console.log(this._input.getText());} ;\n");
		grammarBuilder.append("A : 'b' ;");
		String grammar = grammarBuilder.toString();
		String input ="bc";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("bc\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testParserNotSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(57);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : t=~('x'|'y') 'z' {console.log($t.text);} ;");
		String grammar = grammarBuilder.toString();
		String input ="zz";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("z\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testParserNotToken() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(63);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : ~'x' 'z' {console.log(this._input.getText());} ;");
		String grammar = grammarBuilder.toString();
		String input ="zz";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("zz\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testParserNotTokenWithLabel() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(51);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : t=~'x' 'z' {console.log($t.text);} ;");
		String grammar = grammarBuilder.toString();
		String input ="zz";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("z\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testParserSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(52);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : t=('x'|'y') {console.log($t.text);} ;");
		String grammar = grammarBuilder.toString();
		String input ="x";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("x\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testPlusLexerSingleElement() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(71);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log(this._input.getText());} ;\n");
		grammarBuilder.append("A : 'b'+ 'c' ;");
		String grammar = grammarBuilder.toString();
		String input ="bbbbc";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("bbbbc\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testPlusSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(69);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : ('a'|'b')+ 'c' {console.log(this._input.getText());} ;");
		String grammar = grammarBuilder.toString();
		String input ="abaac";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("abaac\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testRuleAsSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(76);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a @after {console.log(this._input.getText());} : 'a' | 'b' |'c' ;");
		String grammar = grammarBuilder.toString();
		String input ="b";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("b\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testSeqDoesNotBecomeSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(113);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : C {console.log(this._input.getText());} ;\n");
		grammarBuilder.append("fragment A : '1' | '2';\n");
		grammarBuilder.append("fragment B : '3' '4';\n");
		grammarBuilder.append("C : A | B;");
		String grammar = grammarBuilder.toString();
		String input ="34";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("34\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testStarLexerSingleElement_1() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(71);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log(this._input.getText());} ;\n");
		grammarBuilder.append("A : 'b'* 'c' ;");
		String grammar = grammarBuilder.toString();
		String input ="bbbbc";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("bbbbc\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testStarLexerSingleElement_2() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(71);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : A {console.log(this._input.getText());} ;\n");
		grammarBuilder.append("A : 'b'* 'c' ;");
		String grammar = grammarBuilder.toString();
		String input ="c";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("c\n", found);
		assertNull(this.stderrDuringParse);

	}
	/* This file and method are generated by TestGenerator, any edits will be overwritten by the next generation. */
	@Test
	public void testStarSet() throws Exception {
		mkdir(tmpdir);
		StringBuilder grammarBuilder = new StringBuilder(69);
		grammarBuilder.append("grammar T;\n");
		grammarBuilder.append("a : ('a'|'b')* 'c' {console.log(this._input.getText());} ;");
		String grammar = grammarBuilder.toString();
		String input ="abaac";
		String found = execParser("T.g4", grammar, "TParser", "TLexer",
		                          "TListener", "TVisitor",
		                          "a", input, false);
		assertEquals("abaac\n", found);
		assertNull(this.stderrDuringParse);

	}

}
