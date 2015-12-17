package atn

//var LL1Analyzer = require('./../LL1Analyzer').LL1Analyzer
//var IntervalSet = require('./../IntervalSet').IntervalSet

type ATN struct {
    grammarType
    maxTokenType
    states
    decisionToState
    ruleToStartState
    ruleToStopState
    modeNameToStartState
    ruleToTokenType
    lexerActions
    modeToStartState
}

func NewATN(grammarType , maxTokenType) *ATN {

    atn := new(ATN)

    // Used for runtime deserialization of ATNs from strings///
    // The type of the ATN.
    atn.grammarType = grammarType
    // The maximum value for any symbol recognized by a transition in the ATN.
    atn.maxTokenType = maxTokenType
    atn.states = []
    // Each subrule/rule is a decision point and we must track them so we
    //  can go back later and build DFA predictors for them.  This includes
    //  all the rules, subrules, optional blocks, ()+, ()* etc...
    atn.decisionToState = []
    // Maps from rule index to starting state number.
    atn.ruleToStartState = []
    // Maps from rule index to stop state number.
    atn.ruleToStopState = nil
    atn.modeNameToStartState = {}
    // For lexer ATNs, atn.maps the rule index to the resulting token type.
    // For parser ATNs, atn.maps the rule index to the generated bypass token
    // type if the
    // {@link ATNDeserializationOptions//isGenerateRuleBypassTransitions}
    // deserialization option was specified otherwise, atn.is {@code nil}.
    atn.ruleToTokenType = nil
    // For lexer ATNs, atn.is an array of {@link LexerAction} objects which may
    // be referenced by action transitions in the ATN.
    atn.lexerActions = nil
    atn.modeToStartState = []

    return atn

}
	
// Compute the set of valid tokens that can occur starting in state {@code s}.
//  If {@code ctx} is nil, the set of tokens will not include what can follow
//  the rule surrounding {@code s}. In other words, the set will be
//  restricted to tokens reachable staying within {@code s}'s rule.
func (this *ATN) nextTokensInContext(s, ctx) {
    var anal = NewLL1Analyzer(this)
    return anal.LOOK(s, nil, ctx)
}

// Compute the set of valid tokens that can occur starting in {@code s} and
// staying in same rule. {@link Token//EPSILON} is in set if we reach end of
// rule.
func (this *ATN) nextTokensNoContext(s) {
    if (s.nextTokenWithinRule != nil ) {
        return s.nextTokenWithinRule
    }
    s.nextTokenWithinRule = this.nextTokensInContext(s, nil)
    s.nextTokenWithinRule.readOnly = true
    return s.nextTokenWithinRule
}

func (this *ATN) nextTokens(s, ctx) {
    if ( ctx==nil ) {
        return this.nextTokensNoContext(s)
    } else {
        return this.nextTokensInContext(s, ctx)
    }
}

func (this *ATN) addState( state) {
    if ( state != nil ) {
        state.atn = this
        state.stateNumber = this.states.length
    }
    this.states.push(state)
}

func (this *ATN) removeState( state ) {
    this.states[state.stateNumber] = nil // just free mem, don't shift states in list
}

func (this *ATN) defineDecisionState( s) {
    this.decisionToState.push(s)
    s.decision = this.decisionToState.length-1
    return s.decision
}

func (this *ATN) getDecisionState( decision) {
    if (this.decisionToState.length==0) {
        return nil
    } else {
        return this.decisionToState[decision]
    }
}

// Computes the set of input symbols which could follow ATN state number
// {@code stateNumber} in the specified full {@code context}. This method
// considers the complete parser context, but does not evaluate semantic
// predicates (i.e. all predicates encountered during the calculation are
// assumed true). If a path in the ATN exists from the starting state to the
// {@link RuleStopState} of the outermost context without matching any
// symbols, {@link Token//EOF} is added to the returned set.
//
// <p>If {@code context} is {@code nil}, it is treated as
// {@link ParserRuleContext//EMPTY}.</p>
//
// @param stateNumber the ATN state number
// @param context the full parse context
// @return The set of potentially valid input symbols which could follow the
// specified state in the specified context.
// @panics IllegalArgumentException if the ATN does not contain a state with
// number {@code stateNumber}

//var Token = require('./../Token').Token

func (this *ATN) getExpectedTokens( stateNumber, ctx ) {
    if ( stateNumber < 0 || stateNumber >= this.states.length ) {
        panic("Invalid state number.")
    }
    var s = this.states[stateNumber]
    var following = this.nextTokens(s)
    if (!following.contains(TokenEpsilon)) {
        return following
    }
    var expected = NewIntervalSet()
    expected.addSet(following)
    expected.removeOne(TokenEpsilon)
    for (ctx != nil && ctx.invokingState >= 0 && following.contains(TokenEpsilon)) {
        var invokingState = this.states[ctx.invokingState]
        var rt = invokingState.transitions[0]
        following = this.nextTokens(rt.followState)
        expected.addSet(following)
        expected.removeOne(TokenEpsilon)
        ctx = ctx.parentCtx
    }
    if (following.contains(TokenEpsilon)) {
        expected.addOne(TokenEOF)
    }
    return expected
}

var ATNINVALID_ALT_NUMBER = 0

