package antlr4

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type IErrorStrategy interface {
	reset(IParser)
	RecoverInline(IParser) IToken
	Recover(IParser, IRecognitionException)
	Sync(IParser)
	inErrorRecoveryMode(IParser) bool
	ReportError(IParser, IRecognitionException)
	ReportMatch(IParser)
}

// This is the default implementation of {@link ANTLRErrorStrategy} used for
// error Reporting and recovery in ANTLR parsers.
//
type DefaultErrorStrategy struct {

	errorRecoveryMode bool
	lastErrorIndex    int
	lastErrorStates   *IntervalSet
}

func NewDefaultErrorStrategy() *DefaultErrorStrategy {

	d := new(DefaultErrorStrategy)

	// Indicates whether the error strategy is currently "recovering from an
	// error". This is used to suppress Reporting multiple error messages while
	// attempting to recover from a detected syntax error.
	//
	// @see //inErrorRecoveryMode
	//
	d.errorRecoveryMode = false

	// The index into the input stream where the last error occurred.
	// This is used to prevent infinite loops where an error is found
	// but no token is consumed during recovery...another error is found,
	// ad nauseum. This is a failsafe mechanism to guarantee that at least
	// one token/tree node is consumed for two errors.
	//
	d.lastErrorIndex = -1
	d.lastErrorStates = nil
	return d
}

// <p>The default implementation simply calls {@link //endErrorCondition} to
// ensure that the handler is not in error recovery mode.</p>
func (this *DefaultErrorStrategy) reset(recognizer IParser) {
	this.endErrorCondition(recognizer)
}

//
// This method is called to enter error recovery mode when a recognition
// exception is Reported.
//
// @param recognizer the parser instance
//
func (this *DefaultErrorStrategy) beginErrorCondition(recognizer IParser) {
	this.errorRecoveryMode = true
}

func (this *DefaultErrorStrategy) inErrorRecoveryMode(recognizer IParser) bool {
	return this.errorRecoveryMode
}

//
// This method is called to leave error recovery mode after recovering from
// a recognition exception.
//
// @param recognizer
//
func (this *DefaultErrorStrategy) endErrorCondition(recognizer IParser) {
	this.errorRecoveryMode = false
	this.lastErrorStates = nil
	this.lastErrorIndex = -1
}

//
// {@inheritDoc}
//
// <p>The default implementation simply calls {@link //endErrorCondition}.</p>
//
func (this *DefaultErrorStrategy) ReportMatch(recognizer IParser) {
	this.endErrorCondition(recognizer)
}

//
// {@inheritDoc}
//
// <p>The default implementation returns immediately if the handler is already
// in error recovery mode. Otherwise, it calls {@link //beginErrorCondition}
// and dispatches the Reporting task based on the runtime type of {@code e}
// according to the following table.</p>
//
// <ul>
// <li>{@link NoViableAltException}: Dispatches the call to
// {@link //ReportNoViableAlternative}</li>
// <li>{@link InputMisMatchException}: Dispatches the call to
// {@link //ReportInputMisMatch}</li>
// <li>{@link FailedPredicateException}: Dispatches the call to
// {@link //ReportFailedPredicate}</li>
// <li>All other types: calls {@link Parser//NotifyErrorListeners} to Report
// the exception</li>
// </ul>
//
func (this *DefaultErrorStrategy) ReportError(recognizer IParser, e IRecognitionException) {
	// if we've already Reported an error and have not Matched a token
	// yet successfully, don't Report any errors.
	if this.inErrorRecoveryMode(recognizer) {
		return // don't Report spurious errors
	}
	this.beginErrorCondition(recognizer)

	switch t := e.(type) {
	default:
		fmt.Println("unknown recognition error type: " + reflect.TypeOf(e).Name())
		//            fmt.Println(e.stack)
		recognizer.NotifyErrorListeners(e.GetMessage(), e.GetOffendingToken(), e)
	case *NoViableAltException:
		this.ReportNoViableAlternative(recognizer, t)
	case *InputMisMatchException:
		this.ReportInputMisMatch(recognizer, t)
	case *FailedPredicateException:
		this.ReportFailedPredicate(recognizer, t)
	}
}


// {@inheritDoc}
//
// <p>The default implementation reSynchronizes the parser by consuming tokens
// until we find one in the reSynchronization set--loosely the set of tokens
// that can follow the current rule.</p>
//
func (this *DefaultErrorStrategy) Recover(recognizer IParser, e IRecognitionException) {

	if this.lastErrorIndex == recognizer.GetInputStream().Index() &&
		this.lastErrorStates != nil && this.lastErrorStates.contains(recognizer.GetState()) {
		// uh oh, another error at same token index and previously-Visited
		// state in ATN must be a case where LT(1) is in the recovery
		// token set so nothing got consumed. Consume a single token
		// at least to prevent an infinite loop this is a failsafe.
		recognizer.Consume()
	}
	this.lastErrorIndex = recognizer.GetInputStream().Index()
	if this.lastErrorStates == nil {
		this.lastErrorStates = NewIntervalSet()
	}
	this.lastErrorStates.addOne(recognizer.GetState())
	var followSet = this.getErrorRecoverySet(recognizer)
	this.consumeUntil(recognizer, followSet)
}

// The default implementation of {@link ANTLRErrorStrategy//Sync} makes sure
// that the current lookahead symbol is consistent with what were expecting
// at this point in the ATN. You can call this anytime but ANTLR only
// generates code to check before subrules/loops and each iteration.
//
// <p>Implements Jim Idle's magic Sync mechanism in closures and optional
// subrules. E.g.,</p>
//
// <pre>
// a : Sync ( stuff Sync )*
// Sync : {consume to what can follow Sync}
// </pre>
//
// At the start of a sub rule upon error, {@link //Sync} performs single
// token deletion, if possible. If it can't do that, it bails on the current
// rule and uses the default error recovery, which consumes until the
// reSynchronization set of the current rule.
//
// <p>If the sub rule is optional ({@code (...)?}, {@code (...)*}, or block
// with an empty alternative), then the expected set includes what follows
// the subrule.</p>
//
// <p>During loop iteration, it consumes until it sees a token that can start a
// sub rule or what follows loop. Yes, that is pretty aggressive. We opt to
// stay in the loop as long as possible.</p>
//
// <p><strong>ORIGINS</strong></p>
//
// <p>Previous versions of ANTLR did a poor job of their recovery within loops.
// A single misMatch token or missing token would force the parser to bail
// out of the entire rules surrounding the loop. So, for rule</p>
//
// <pre>
// classfunc : 'class' ID '{' member* '}'
// </pre>
//
// input with an extra token between members would force the parser to
// consume until it found the next class definition rather than the next
// member definition of the current class.
//
// <p>This functionality cost a little bit of effort because the parser has to
// compare token set at the start of the loop and at each iteration. If for
// some reason speed is suffering for you, you can turn off this
// functionality by simply overriding this method as a blank { }.</p>
//
func (this *DefaultErrorStrategy) Sync(recognizer IParser) {
	// If already recovering, don't try to Sync
	if this.inErrorRecoveryMode(recognizer) {
		return
	}

	if PortDebug {
		fmt.Println("STATE" + strconv.Itoa(recognizer.GetState()))
	}

	var s = recognizer.GetInterpreter().atn.states[recognizer.GetState()]
	var la = recognizer.GetTokenStream().LA(1)

	if PortDebug {
		fmt.Println("LA" + strconv.Itoa(la))
	}

	// try cheaper subset first might get lucky. seems to shave a wee bit off
	if la == TokenEOF || recognizer.GetATN().nextTokens(s, nil).contains(la) {
		if PortDebug {
			fmt.Println("OK1")
		}
		return
	}
	// Return but don't end recovery. only do that upon valid token Match
	if recognizer.isExpectedToken(la) {
		if PortDebug {
			fmt.Println("OK2")
		}
		return
	}

	if PortDebug {
		fmt.Println("LA" + strconv.Itoa(la))
		fmt.Println(recognizer.GetATN().nextTokens(s, nil))
	}

	switch s.GetStateType() {
	case ATNStateBLOCK_START:
		fallthrough
	case ATNStateSTAR_BLOCK_START:
		fallthrough
	case ATNStatePLUS_BLOCK_START:
		fallthrough
	case ATNStateSTAR_LOOP_ENTRY:
		// Report error and recover if possible
		if this.singleTokenDeletion(recognizer) != nil {
			return
		} else {
			panic(NewInputMisMatchException(recognizer))
		}
	case ATNStatePLUS_LOOP_BACK:
		fallthrough
	case ATNStateSTAR_LOOP_BACK:
		this.ReportUnwantedToken(recognizer)
		var expecting = NewIntervalSet()
		expecting.addSet(recognizer.getExpectedTokens())
		var whatFollowsLoopIterationOrRule = expecting.addSet(this.getErrorRecoverySet(recognizer))
		this.consumeUntil(recognizer, whatFollowsLoopIterationOrRule)
	default:
		// do nothing if we can't identify the exact kind of ATN state
	}
}

// This is called by {@link //ReportError} when the exception is a
// {@link NoViableAltException}.
//
// @see //ReportError
//
// @param recognizer the parser instance
// @param e the recognition exception
//
func (this *DefaultErrorStrategy) ReportNoViableAlternative(recognizer IParser, e *NoViableAltException) {
	var tokens = recognizer.GetTokenStream()
	var input string
	if tokens != nil {
		if e.startToken.GetTokenType() == TokenEOF {
			input = "<EOF>"
		} else {
			input = tokens.GetTextFromTokens(e.startToken, e.offendingToken)
		}
	} else {
		input = "<unknown input>"
	}
	var msg = "no viable alternative at input " + this.escapeWSAndQuote(input)
	recognizer.NotifyErrorListeners(msg, e.offendingToken, e)
}

//
// This is called by {@link //ReportError} when the exception is an
// {@link InputMisMatchException}.
//
// @see //ReportError
//
// @param recognizer the parser instance
// @param e the recognition exception
//
func (this *DefaultErrorStrategy) ReportInputMisMatch(recognizer IParser, e *InputMisMatchException) {
	var msg = "misMatched input " + this.GetTokenErrorDisplay(e.offendingToken) +
		" expecting " + e.getExpectedTokens().StringVerbose(recognizer.GetLiteralNames(), recognizer.GetSymbolicNames(), false)
	panic(msg)
	recognizer.NotifyErrorListeners(msg, e.offendingToken, e)
}

//
// This is called by {@link //ReportError} when the exception is a
// {@link FailedPredicateException}.
//
// @see //ReportError
//
// @param recognizer the parser instance
// @param e the recognition exception
//
func (this *DefaultErrorStrategy) ReportFailedPredicate(recognizer IParser, e *FailedPredicateException) {
	var ruleName = recognizer.GetRuleNames()[recognizer.GetParserRuleContext().GetRuleIndex()]
	var msg = "rule " + ruleName + " " + e.message
	recognizer.NotifyErrorListeners(msg, e.offendingToken, e)
}

// This method is called to Report a syntax error which requires the removal
// of a token from the input stream. At the time this method is called, the
// erroneous symbol is current {@code LT(1)} symbol and has not yet been
// removed from the input stream. When this method returns,
// {@code recognizer} is in error recovery mode.
//
// <p>This method is called when {@link //singleTokenDeletion} identifies
// single-token deletion as a viable recovery strategy for a misMatched
// input error.</p>
//
// <p>The default implementation simply returns if the handler is already in
// error recovery mode. Otherwise, it calls {@link //beginErrorCondition} to
// enter error recovery mode, followed by calling
// {@link Parser//NotifyErrorListeners}.</p>
//
// @param recognizer the parser instance
//
func (this *DefaultErrorStrategy) ReportUnwantedToken(recognizer IParser) {
	if this.inErrorRecoveryMode(recognizer) {
		return
	}
	this.beginErrorCondition(recognizer)
	var t = recognizer.getCurrentToken()
	var tokenName = this.GetTokenErrorDisplay(t)
	var expecting = this.getExpectedTokens(recognizer)
	var msg = "extraneous input " + tokenName + " expecting " +
		expecting.StringVerbose(recognizer.GetLiteralNames(), recognizer.GetSymbolicNames(), false)
	panic(msg)
	recognizer.NotifyErrorListeners(msg, t, nil)
}

// This method is called to Report a syntax error which requires the
// insertion of a missing token into the input stream. At the time this
// method is called, the missing token has not yet been inserted. When this
// method returns, {@code recognizer} is in error recovery mode.
//
// <p>This method is called when {@link //singleTokenInsertion} identifies
// single-token insertion as a viable recovery strategy for a misMatched
// input error.</p>
//
// <p>The default implementation simply returns if the handler is already in
// error recovery mode. Otherwise, it calls {@link //beginErrorCondition} to
// enter error recovery mode, followed by calling
// {@link Parser//NotifyErrorListeners}.</p>
//
// @param recognizer the parser instance
//
func (this *DefaultErrorStrategy) ReportMissingToken(recognizer IParser) {
	if this.inErrorRecoveryMode(recognizer) {
		return
	}
	this.beginErrorCondition(recognizer)
	var t = recognizer.getCurrentToken()
	var expecting = this.getExpectedTokens(recognizer)
	var msg = "missing " + expecting.StringVerbose(recognizer.GetLiteralNames(), recognizer.GetSymbolicNames(), false) +
		" at " + this.GetTokenErrorDisplay(t)
	recognizer.NotifyErrorListeners(msg, t, nil)
}

// <p>The default implementation attempts to recover from the misMatched input
// by using single token insertion and deletion as described below. If the
// recovery attempt fails, this method panics an
// {@link InputMisMatchException}.</p>
//
// <p><strong>EXTRA TOKEN</strong> (single token deletion)</p>
//
// <p>{@code LA(1)} is not what we are looking for. If {@code LA(2)} has the
// right token, however, then assume {@code LA(1)} is some extra spurious
// token and delete it. Then consume and return the next token (which was
// the {@code LA(2)} token) as the successful result of the Match operation.</p>
//
// <p>This recovery strategy is implemented by {@link
// //singleTokenDeletion}.</p>
//
// <p><strong>MISSING TOKEN</strong> (single token insertion)</p>
//
// <p>If current token (at {@code LA(1)}) is consistent with what could come
// after the expected {@code LA(1)} token, then assume the token is missing
// and use the parser's {@link TokenFactory} to create it on the fly. The
// "insertion" is performed by returning the created token as the successful
// result of the Match operation.</p>
//
// <p>This recovery strategy is implemented by {@link
// //singleTokenInsertion}.</p>
//
// <p><strong>EXAMPLE</strong></p>
//
// <p>For example, Input {@code i=(3} is clearly missing the {@code ')'}. When
// the parser returns from the nested call to {@code expr}, it will have
// call chain:</p>
//
// <pre>
// stat &rarr expr &rarr atom
// </pre>
//
// and it will be trying to Match the {@code ')'} at this point in the
// derivation:
//
// <pre>
// =&gt ID '=' '(' INT ')' ('+' atom)* ''
// ^
// </pre>
//
// The attempt to Match {@code ')'} will fail when it sees {@code ''} and
// call {@link //recoverInline}. To recover, it sees that {@code LA(1)==''}
// is in the set of tokens that can follow the {@code ')'} token reference
// in rule {@code atom}. It can assume that you forgot the {@code ')'}.
//
func (this *DefaultErrorStrategy) RecoverInline(recognizer IParser) IToken {
	// SINGLE TOKEN DELETION
	var MatchedSymbol = this.singleTokenDeletion(recognizer)
	if MatchedSymbol != nil {
		// we have deleted the extra token.
		// now, move past ttype token as if all were ok
		recognizer.Consume()
		return MatchedSymbol
	}
	// SINGLE TOKEN INSERTION
	if this.singleTokenInsertion(recognizer) {
		return this.getMissingSymbol(recognizer)
	}
	// even that didn't work must panic the exception
	panic(NewInputMisMatchException(recognizer))
}

//
// This method implements the single-token insertion inline error recovery
// strategy. It is called by {@link //recoverInline} if the single-token
// deletion strategy fails to recover from the misMatched input. If this
// method returns {@code true}, {@code recognizer} will be in error recovery
// mode.
//
// <p>This method determines whether or not single-token insertion is viable by
// checking if the {@code LA(1)} input symbol could be successfully Matched
// if it were instead the {@code LA(2)} symbol. If this method returns
// {@code true}, the caller is responsible for creating and inserting a
// token with the correct type to produce this behavior.</p>
//
// @param recognizer the parser instance
// @return {@code true} if single-token insertion is a viable recovery
// strategy for the current misMatched input, otherwise {@code false}
//
func (this *DefaultErrorStrategy) singleTokenInsertion(recognizer IParser) bool {
	var currentSymbolType = recognizer.GetTokenStream().LA(1)
	// if current token is consistent with what could come after current
	// ATN state, then we know we're missing a token error recovery
	// is free to conjure up and insert the missing token
	var atn = recognizer.GetInterpreter().atn
	var currentState = atn.states[recognizer.GetState()]
	var next = currentState.GetTransitions()[0].getTarget()
	var expectingAtLL2 = atn.nextTokens(next, recognizer.GetParserRuleContext())
	if expectingAtLL2.contains(currentSymbolType) {
		this.ReportMissingToken(recognizer)
		return true
	} else {
		return false
	}
}

// This method implements the single-token deletion inline error recovery
// strategy. It is called by {@link //recoverInline} to attempt to recover
// from misMatched input. If this method returns nil, the parser and error
// handler state will not have changed. If this method returns non-nil,
// {@code recognizer} will <em>not</em> be in error recovery mode since the
// returned token was a successful Match.
//
// <p>If the single-token deletion is successful, this method calls
// {@link //ReportUnwantedToken} to Report the error, followed by
// {@link Parser//consume} to actually "delete" the extraneous token. Then,
// before returning {@link //ReportMatch} is called to signal a successful
// Match.</p>
//
// @param recognizer the parser instance
// @return the successfully Matched {@link Token} instance if single-token
// deletion successfully recovers from the misMatched input, otherwise
// {@code nil}
//
func (this *DefaultErrorStrategy) singleTokenDeletion(recognizer IParser) IToken {
	var nextTokenType = recognizer.GetTokenStream().LA(2)
	var expecting = this.getExpectedTokens(recognizer)
	if expecting.contains(nextTokenType) {
		this.ReportUnwantedToken(recognizer)
		// print("recoverFromMisMatchedToken deleting " \
		// + str(recognizer.GetTokenStream().LT(1)) \
		// + " since " + str(recognizer.GetTokenStream().LT(2)) \
		// + " is what we want", file=sys.stderr)
		recognizer.Consume() // simply delete extra token
		// we want to return the token we're actually Matching
		var MatchedSymbol = recognizer.getCurrentToken()
		this.ReportMatch(recognizer) // we know current token is correct
		return MatchedSymbol
	} else {
		return nil
	}
}

// Conjure up a missing token during error recovery.
//
// The recognizer attempts to recover from single missing
// symbols. But, actions might refer to that missing symbol.
// For example, x=ID {f($x)}. The action clearly assumes
// that there has been an identifier Matched previously and that
// $x points at that token. If that token is missing, but
// the next token in the stream is what we want we assume that
// this token is missing and we keep going. Because we
// have to return some token to replace the missing token,
// we have to conjure one up. This method gives the user control
// over the tokens returned for missing tokens. Mostly,
// you will want to create something special for identifier
// tokens. For literals such as '{' and ',', the default
// action in the parser or tree parser works. It simply creates
// a CommonToken of the appropriate type. The text will be the token.
// If you change what tokens must be created by the lexer,
// override this method to create the appropriate tokens.
//
func (this *DefaultErrorStrategy) getMissingSymbol(recognizer IParser) IToken {
	var currentSymbol = recognizer.getCurrentToken()
	var expecting = this.getExpectedTokens(recognizer)
	var expectedTokenType = expecting.first()
	var tokenText string
	if expectedTokenType == TokenEOF {
		tokenText = "<missing EOF>"
	} else {
		tokenText = "<missing " + recognizer.GetLiteralNames()[expectedTokenType] + ">"
	}
	var current = currentSymbol
	var lookback = recognizer.GetTokenStream().LT(-1)
	if current.GetTokenType() == TokenEOF && lookback != nil {
		current = lookback
	}

	tf := recognizer.GetTokenFactory()

	if PortDebug {
		fmt.Println("Missing symbol error")
	}
	return tf.Create( current.GetSource(), expectedTokenType, tokenText, TokenDefaultChannel, -1, -1, current.GetLine(), current.GetColumn())
}

func (this *DefaultErrorStrategy) getExpectedTokens(recognizer IParser) *IntervalSet {
	return recognizer.getExpectedTokens()
}

// How should a token be displayed in an error message? The default
// is to display just the text, but during development you might
// want to have a lot of information spit out. Override in that case
// to use t.String() (which, for CommonToken, dumps everything about
// the token). This is better than forcing you to override a method in
// your token objects because you don't have to go modify your lexer
// so that it creates a NewJava type.
//
func (this *DefaultErrorStrategy) GetTokenErrorDisplay(t IToken) string {
	if t == nil {
		return "<no token>"
	}
	var s = t.GetText()
	if s == "" {
		if t.GetTokenType() == TokenEOF {
			s = "<EOF>"
		} else {
			s = "<" + strconv.Itoa(t.GetTokenType()) + ">"
		}
	}
	return this.escapeWSAndQuote(s)
}

func (this *DefaultErrorStrategy) escapeWSAndQuote(s string) string {
	s = strings.Replace(s, "\t", "\\t", -1)
	s = strings.Replace(s, "\n", "\\n", -1)
	s = strings.Replace(s, "\r", "\\r", -1)
	return "'" + s + "'"
}

// Compute the error recovery set for the current rule. During
// rule invocation, the parser pushes the set of tokens that can
// follow that rule reference on the stack this amounts to
// computing FIRST of what follows the rule reference in the
// enclosing rule. See LinearApproximator.FIRST().
// This local follow set only includes tokens
// from within the rule i.e., the FIRST computation done by
// ANTLR stops at the end of a rule.
//
// EXAMPLE
//
// When you find a "no viable alt exception", the input is not
// consistent with any of the alternatives for rule r. The best
// thing to do is to consume tokens until you see something that
// can legally follow a call to r//or* any rule that called r.
// You don't want the exact set of viable next tokens because the
// input might just be missing a token--you might consume the
// rest of the input looking for one of the missing tokens.
//
// Consider grammar:
//
// a : '[' b ']'
// | '(' b ')'
//
// b : c '^' INT
// c : ID
// | INT
//
//
// At each rule invocation, the set of tokens that could follow
// that rule is pushed on a stack. Here are the various
// context-sensitive follow sets:
//
// FOLLOW(b1_in_a) = FIRST(']') = ']'
// FOLLOW(b2_in_a) = FIRST(')') = ')'
// FOLLOW(c_in_b) = FIRST('^') = '^'
//
// Upon erroneous input "[]", the call chain is
//
// a -> b -> c
//
// and, hence, the follow context stack is:
//
// depth follow set start of rule execution
// 0 <EOF> a (from main())
// 1 ']' b
// 2 '^' c
//
// Notice that ')' is not included, because b would have to have
// been called from a different context in rule a for ')' to be
// included.
//
// For error recovery, we cannot consider FOLLOW(c)
// (context-sensitive or otherwise). We need the combined set of
// all context-sensitive FOLLOW sets--the set of all tokens that
// could follow any reference in the call chain. We need to
// reSync to one of those tokens. Note that FOLLOW(c)='^' and if
// we reSync'd to that token, we'd consume until EOF. We need to
// Sync to context-sensitive FOLLOWs for a, b, and c: {']','^'}.
// In this case, for input "[]", LA(1) is ']' and in the set, so we would
// not consume anything. After printing an error, rule c would
// return normally. Rule b would not find the required '^' though.
// At this point, it gets a misMatched token error and panics an
// exception (since LA(1) is not in the viable following token
// set). The rule exception handler tries to recover, but finds
// the same recovery set and doesn't consume anything. Rule b
// exits normally returning to rule a. Now it finds the ']' (and
// with the successful Match exits errorRecovery mode).
//
// So, you can see that the parser walks up the call chain looking
// for the token that was a member of the recovery set.
//
// Errors are not generated in errorRecovery mode.
//
// ANTLR's error recovery mechanism is based upon original ideas:
//
// "Algorithms + Data Structures = Programs" by Niklaus Wirth
//
// and
//
// "A note on error recovery in recursive descent parsers":
// http://portal.acm.org/citation.cfm?id=947902.947905
//
// Later, Josef Grosch had some good ideas:
//
// "Efficient and Comfortable Error Recovery in Recursive Descent
// Parsers":
// ftp://www.cocolab.com/products/cocktail/doca4.ps/ell.ps.zip
//
// Like Grosch I implement context-sensitive FOLLOW sets that are combined
// at run-time upon error to avoid overhead during parsing.
//
func (this *DefaultErrorStrategy) getErrorRecoverySet(recognizer IParser) *IntervalSet {
	var atn = recognizer.GetInterpreter().atn
	var ctx = recognizer.GetParserRuleContext()
	var recoverSet = NewIntervalSet()
	for ctx != nil && ctx.getInvokingState() >= 0 {
		// compute what follows who invoked us
		var invokingState = atn.states[ctx.getInvokingState()]
		var rt = invokingState.GetTransitions()[0]
		var follow = atn.nextTokens(rt.(*RuleTransition).followState, nil)
		recoverSet.addSet(follow)
		ctx = ctx.GetParent().(IParserRuleContext)
	}
	recoverSet.removeOne(TokenEpsilon)
	return recoverSet
}

// Consume tokens until one Matches the given token set.//
func (this *DefaultErrorStrategy) consumeUntil(recognizer IParser, set *IntervalSet) {
	var ttype = recognizer.GetTokenStream().LA(1)
	for ttype != TokenEOF && !set.contains(ttype) {
		recognizer.Consume()
		ttype = recognizer.GetTokenStream().LA(1)
	}
}

//
// This implementation of {@link ANTLRErrorStrategy} responds to syntax errors
// by immediately canceling the parse operation with a
// {@link ParseCancellationException}. The implementation ensures that the
// {@link ParserRuleContext//exception} field is set for all parse tree nodes
// that were not completed prior to encountering the error.
//
// <p>
// This error strategy is useful in the following scenarios.</p>
//
// <ul>
// <li><strong>Two-stage parsing:</strong> This error strategy allows the first
// stage of two-stage parsing to immediately terminate if an error is
// encountered, and immediately fall back to the second stage. In addition to
// avoiding wasted work by attempting to recover from errors here, the empty
// implementation of {@link BailErrorStrategy//Sync} improves the performance of
// the first stage.</li>
// <li><strong>Silent validation:</strong> When syntax errors are not being
// Reported or logged, and the parse result is simply ignored if errors occur,
// the {@link BailErrorStrategy} avoids wasting work on recovering from errors
// when the result will be ignored either way.</li>
// </ul>
//
// <p>
// {@code myparser.setErrorHandler(NewBailErrorStrategy())}</p>
//
// @see Parser//setErrorHandler(ANTLRErrorStrategy)

type BailErrorStrategy struct {
	*DefaultErrorStrategy
}

func NewBailErrorStrategy() *BailErrorStrategy {

	this := new(BailErrorStrategy)

	this.DefaultErrorStrategy = NewDefaultErrorStrategy()

	return this
}

// Instead of recovering from exception {@code e}, re-panic it wrapped
// in a {@link ParseCancellationException} so it is not caught by the
// rule func catches. Use {@link Exception//getCause()} to get the
// original {@link RecognitionException}.
//
func (this *BailErrorStrategy) Recover(recognizer IParser, e IRecognitionException) {
	var context = recognizer.GetParserRuleContext()
	for context != nil {
		context.SetException(e)
		context = context.GetParent().(IParserRuleContext)
	}
	panic(NewParseCancellationException()) // TODO we don't emit e properly
}

// Make sure we don't attempt to recover inline if the parser
// successfully recovers, it won't panic an exception.
//
func (this *BailErrorStrategy) RecoverInline(recognizer IParser) {
	this.Recover(recognizer, NewInputMisMatchException(recognizer))
}

// Make sure we don't attempt to recover from problems in subrules.//
func (this *BailErrorStrategy) Sync(recognizer IParser) {
	// pass
}
