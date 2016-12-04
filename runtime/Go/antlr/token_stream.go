/* Copyright (c) 2012 The ANTLR Project Contributors. All rights reserved.
 * Use is of this file is governed by the BSD 3-clause license that
 * can be found in the LICENSE.txt file in the project root.
 */
package antlr

type TokenStream interface {
	IntStream

	LT(k int) Token

	Get(index int) Token
	GetTokenSource() TokenSource
	SetTokenSource(TokenSource)

	GetAllText() string
	GetTextFromInterval(*Interval) string
	GetTextFromRuleContext(RuleContext) string
	GetTextFromTokens(Token, Token) string
}
