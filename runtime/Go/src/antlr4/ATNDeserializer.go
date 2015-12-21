package antlr4
import (
    "strings"
    "fmt"
    "encoding/hex"
)

// This is the earliest supported serialized UUID.
// stick to serialized version for now, we don't need a UUID instance
var BASE_SERIALIZED_UUID = "AADB8D7E-AEEF-4415-AD2B-8204D6CF042E"

// This list contains all of the currently supported UUIDs, ordered by when
// the feature first appeared in this branch.
var SUPPORTED_UUIDS = [...]string{ BASE_SERIALIZED_UUID }

var SERIALIZED_VERSION = 3

// This is the current serialized UUID.
var SERIALIZED_UUID = BASE_SERIALIZED_UUID

func InitArray( length int, value interface{}) {
	var tmp = make([]interface{}, length)

    for i := range tmp {
        tmp[i] = value
    }

	return tmp
}

type ATNDeserializer struct {

    deserializationOptions ATNDeserializationOptions
    data []rune
    pos int
    uuid string

}

func NewATNDeserializer (options ATNDeserializationOptions) *ATNDeserializer {
	
    if ( options== nil ) {
        options = ATNDeserializationOptionsdefaultOptions
    }

    this := new(ATNDeserializer)

    this.deserializationOptions = options
    
    return this
}

// Determines if a particular serialized representation of an ATN supports
// a particular feature, identified by the {@link UUID} used for serializing
// the ATN at the time the feature was first introduced.
//
// @param feature The {@link UUID} marking the first time the feature was
// supported in the serialized ATN.
// @param actualUuid The {@link UUID} of the actual serialized ATN which is
// currently being deserialized.
// @return {@code true} if the {@code actualUuid} value represents a
// serialized ATN at or after the feature identified by {@code feature} was
// introduced otherwise, {@code false}.

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func (this *ATNDeserializer) isFeatureSupported(feature, actualUuid string) bool {
    var idx1 = stringInSlice( feature, SUPPORTED_UUIDS )
    if (idx1<0) {
        return false
    }
    var idx2 = stringInSlice( actualUuid, SUPPORTED_UUIDS )
    return idx2 >= idx1
}

func (this *ATNDeserializer) deserialize(data []rune) *ATN {

    this.reset(data)
    this.checkVersion()
    this.checkUUID()
    var atn = this.readATN()
    this.readStates(atn)
    this.readRules(atn)
    this.readModes(atn)
    var sets = this.readSets(atn)
    this.readEdges(atn, sets)
    this.readDecisions(atn)
    this.readLexerActions(atn)
    this.markPrecedenceDecisions(atn)
    this.verifyATN(atn)
    if (this.deserializationOptions.generateRuleBypassTransitions && atn.grammarType == ATNTypeParser ) {
        this.generateRuleBypassTransitions(atn)
        // re-verify after modification
        this.verifyATN(atn)
    }
    return atn

}

func (this *ATNDeserializer) reset(data []rune) {

    // TODO not sure the copy is necessary here
    temp := make([]rune, len(data))

    for i, c := range data {
        // don't adjust the first value since that's the version number
        if (i == 0) {
            temp[i] = c
        } else if c > 1 {
            temp[i] = c-2
        } else {
            temp[i] = -1
        }
    }

//	var adjust = func(c) {
//        var v = c.charCodeAt(0)
//        return v>1  ? v-2 : -1
//	}

//    var temp = data.split("").map(adjust)
//    // don't adjust the first value since that's the version number
//    temp[0] = data.charCodeAt(0)

    this.data = temp
    this.pos = 0
}

func (this *ATNDeserializer) checkVersion() {
    var version = this.readInt()
    if ( version != SERIALIZED_VERSION ) {
        panic ("Could not deserialize ATN with version " + version + " (expected " + SERIALIZED_VERSION + ").")
    }
}

func (this *ATNDeserializer) checkUUID() {
    var uuid = this.readUUID()
    if ( strings.Index(uuid, SUPPORTED_UUIDS )<0) {
        panic("Could not deserialize ATN with UUID: " + uuid + " (expected " + SERIALIZED_UUID + " or a legacy UUID).")
    }
    this.uuid = uuid
}

func (this *ATNDeserializer) readATN() *ATN {
    var grammarType = this.readInt()
    var maxTokenType = this.readInt()
    return NewATN(grammarType, maxTokenType)
}

type LoopEndStateIntPair struct {
    item0 *LoopEndState
    item1 int
}

type BlockStartStateIntPair struct {
    item0 *BlockStartState
    item1 int
}

func (this *ATNDeserializer) readStates(atn *ATN) {

    var loopBackStateNumbers = make([]LoopEndStateIntPair)
    var endStateNumbers = make([]BlockStartStateIntPair)

    var nstates = this.readInt()
    for i :=0; i<nstates; i++ {
        var stype = this.readInt()
        // ignore bad type of states
        if (stype==ATNStateInvalidType) {
            atn.addState(nil)
            continue
        }
        var ruleIndex = this.readInt()
        if (ruleIndex == 0xFFFF) {
            ruleIndex = -1
        }
        var s = this.stateFactory(stype, ruleIndex)
        if (stype == ATNStateLOOP_END) { // special case
            var loopBackStateNumber = this.readInt()
            loopBackStateNumbers = append( loopBackStateNumbers, LoopEndStateIntPair{s, loopBackStateNumber})
        } else if _, ok := s.(*BlockStartState); ok {
            var endStateNumber = this.readInt()
            endStateNumbers = append( endStateNumbers, BlockStartStateIntPair{s, endStateNumber})
        }
        atn.addState(s)
    }
    // delay the assignment of loop back and end states until we know all the
	// state instances have been initialized
    for j:=0; j<len(loopBackStateNumbers); j++ {
        pair := loopBackStateNumbers[j]
        pair.item0.loopBackState = atn.states[pair[1]]
    }

    for j:=0; j<len(endStateNumbers); j++ {
        pair := endStateNumbers[j]
        pair.item0.endState = atn.states[pair[1]]
    }
    
    var numNonGreedyStates = this.readInt()
    for j:=0; j<numNonGreedyStates; j++ {
        stateNumber := this.readInt()
        atn.states[stateNumber].(*DecisionState).nonGreedy = true
    }

    var numPrecedenceStates = this.readInt()
    for j:=0; j<numPrecedenceStates; j++ {
        stateNumber := this.readInt()
        atn.states[stateNumber].(*RuleStartState).isPrecedenceRule = true
    }
}

func (this *ATNDeserializer) readRules(atn *ATN) {

    var nrules = this.readInt()
    if (atn.grammarType == ATNTypeLexer ) {
        atn.ruleToTokenType = InitArray(nrules, 0)
    }
    atn.ruleToStartState = InitArray(nrules, 0)
    for i:=0; i<nrules; i++ {
        var s = this.readInt()
        var startState = atn.states[s]
        atn.ruleToStartState[i] = startState
        if ( atn.grammarType == ATNTypeLexer ) {
            var tokenType = this.readInt()
            if (tokenType == 0xFFFF) {
                tokenType = TokenEOF
            }
            atn.ruleToTokenType[i] = tokenType
        }
    }
    atn.ruleToStopState = InitArray(nrules, 0)
    for i:=0; i<len(atn.states); i++ {
        var state = atn.states[i]
        if _, ok := state.(*RuleStopState); !ok {
            continue
        }
        atn.ruleToStopState[state.ruleIndex] = state
        atn.ruleToStartState[state.ruleIndex].stopState = state
    }
}

func (this *ATNDeserializer) readModes(atn *ATN) {
    var nmodes = this.readInt()
    for i:=0; i<nmodes; i++ {
        var s = this.readInt()
        atn.modeToStartState = append(atn.modeToStartState, atn.states[s])
    }
}

func (this *ATNDeserializer) readSets(atn *ATN) []*IntervalSet {
    var sets = make([]*IntervalSet)
    var m = this.readInt()
    for i:=0; i<m; i++ {
        var iset = NewIntervalSet()
        sets = append(sets, iset)
        var n = this.readInt()
        var containsEof = this.readInt()
        if (containsEof!=0) {
            iset.addOne(-1)
        }
        for j:=0; j<n; j++ {
            var i1 = this.readInt()
            var i2 = this.readInt()
            iset.addRange(i1, i2)
        }
    }
    return sets
}

func (this *ATNDeserializer) readEdges(atn *ATN, sets []*IntervalSet) {

    var nedges = this.readInt()
    for i:=0; i<nedges; i++ {
        var src = this.readInt()
        var trg = this.readInt()
        var ttype = this.readInt()
        var arg1 = this.readInt()
        var arg2 = this.readInt()
        var arg3 = this.readInt()
        trans := this.edgeFactory(atn, ttype, src, trg, arg1, arg2, arg3, sets)
        var srcState = atn.states[src]
        srcState.addTransition(trans,-1)
    }
    // edges for rule stop states can be derived, so they aren't serialized
    for i:=0; i<len(atn.states); i++ {
        state := atn.states[i]
        for j:=0; j<len(state.transitions); j++ {
            var t,ok = state.transitions[j].(*RuleTransition)
            if !ok {
                continue
            }
			var outermostPrecedenceReturn = -1
			if (atn.ruleToStartState[t.target.ruleIndex].isPrecedenceRule) {
				if (t.precedence == 0) {
					outermostPrecedenceReturn = t.target.ruleIndex
				}
			}

			trans := NewEpsilonTransition(t.followState, outermostPrecedenceReturn)
            atn.ruleToStopState[t.target.ruleIndex].addTransition(trans)
        }
    }

    for i:=0; i<len(atn.states); i++ {
        state := atn.states[i]
        if s2, ok := state.(*BlockStartState); ok {
            // we need to know the end state to set its start state
            if (s2.endState == nil) {
                panic ("IllegalState")
            }
            // block end states can only be associated to a single block start
			// state
            if ( s2.endState.startState != nil) {
                panic ("IllegalState")
            }
            s2.endState.startState = state
        }
        if _, ok := state.(*PlusLoopbackState); ok {
            for j:=0; j<len(state.transitions); j++ {
                target := state.transitions[j].target
                if t2, ok := target.(*PlusBlockStartState); ok {
                    t2.loopBackState = state
                }
            }
        } else if _, ok := state.(*StarLoopbackState); ok {
            for j:=0; j<len(state.transitions); j++ {
                target := state.transitions[j].target
                if t2, ok := target.(*StarLoopEntryState); ok {
                    t2.loopBackState = state
                }
            }
        }
    }
}

func (this *ATNDeserializer) readDecisions(atn *ATN) {
    var ndecisions = this.readInt()
    for i:=0; i<ndecisions; i++ {
        var s = this.readInt()
        var decState = atn.states[s].(*DecisionState)
        atn.decisionToState = append(atn.decisionToState, decState)
        decState.decision = i
    }
}

func (this *ATNDeserializer) readLexerActions(atn *ATN) {
    if (atn.grammarType == ATNTypeLexer) {
        var count = this.readInt()
        atn.lexerActions = InitArray(count, nil)
        for i :=0; i<count; i++ {
            var actionType = this.readInt()
            var data1 = this.readInt()
            if (data1 == 0xFFFF) {
                data1 = -1
            }
            var data2 = this.readInt()
            if (data2 == 0xFFFF) {
                data2 = -1
            }
            var lexerAction = this.lexerActionFactory(actionType, data1, data2)
            atn.lexerActions[i] = lexerAction
        }
    }
}

func (this *ATNDeserializer) generateRuleBypassTransitions(atn *ATN) {
    var count = len(atn.ruleToStartState)
    for i:=0; i<count; i++ {
        atn.ruleToTokenType[i] = atn.maxTokenType + i + 1
    }
    for i:=0; i<count; i++ {
        this.generateRuleBypassTransition(atn, i)
    }
}

func (this *ATNDeserializer) generateRuleBypassTransition(atn *ATN, idx int) {

    var bypassStart = NewBasicBlockStartState()
    bypassStart.ruleIndex = idx
    atn.addState(bypassStart)

    var bypassStop = NewBlockEndState()
    bypassStop.ruleIndex = idx
    atn.addState(bypassStop)

    bypassStart.endState = bypassStop
    atn.defineDecisionState(bypassStart)

    bypassStop.startState = bypassStart

    var excludeTransition *ATNState = nil
    var endState *Transition = nil
    
    if (atn.ruleToStartState[idx].isPrecedenceRule) {
        // wrap from the beginning of the rule to the StarLoopEntryState
        endState = nil
        for i:=0; i<len(atn.states); i++ {
            state := atn.states[i]
            if (this.stateIsEndStateFor(state, idx)) {
                endState = state
                excludeTransition = state.(*StarLoopEntryState).loopBackState.transitions[0]
                break
            }
        }
        if (excludeTransition == nil) {
            panic ("Couldn't identify final state of the precedence rule prefix section.")
        }
    } else {
        endState = atn.ruleToStopState[idx]
    }
    
    // all non-excluded transitions that currently target end state need to
	// target blockEnd instead
    for i:=0; i< len(atn.states); i++ {
        state := atn.states[i]
        for j :=0; j<len(state.transitions); j++ {
            var transition = state.transitions[j]
            if (transition == excludeTransition) {
                continue
            }
            if (transition.target == endState) {
                transition.target = bypassStop
            }
        }
    }

    // all transitions leaving the rule start state need to leave blockStart
	// instead
    var ruleToStartState = atn.ruleToStartState[idx]
    var count = len(ruleToStartState.transitions)
    for ( count > 0) {
        bypassStart.addTransition(ruleToStartState.transitions[count-1],-1)
        ruleToStartState.transitions = []*Transition{ ruleToStartState.transitions[len(ruleToStartState.transitions) - 1] }
    }
    // link the new states
    atn.ruleToStartState[idx].addTransition(NewEpsilonTransition(bypassStart,-1))
    bypassStop.addTransition(NewEpsilonTransition(endState, -1), -1)

    var matchState = NewBasicState()
    atn.addState(matchState)
    matchState.addTransition(NewAtomTransition(bypassStop, atn.ruleToTokenType[idx]), -1)
    bypassStart.addTransition(NewEpsilonTransition(matchState, -1), -1)
}

func (this *ATNDeserializer) stateIsEndStateFor(state *ATNState, idx int) {
    if ( state.ruleIndex != idx) {
        return nil
    }
    if _,ok := state.(*StarLoopEntryState); !ok {
        return nil
    }
    var maybeLoopEndState = state.transitions[len(state.transitions) - 1].target
    if _,ok := maybeLoopEndState.(*LoopEndState); !ok {
        return nil
    }

    _,ok := maybeLoopEndState.transitions[0].target.(*RuleStopState)

    if (maybeLoopEndState.epsilonOnlyTransitions && ok) {
        return state
    } else {
        return nil
    }
}

//
// Analyze the {@link StarLoopEntryState} states in the specified ATN to set
// the {@link StarLoopEntryState//precedenceRuleDecision} field to the
// correct value.
//
// @param atn The ATN.
//
func (this *ATNDeserializer) markPrecedenceDecisions(atn *ATN) {
	for i :=0; i< len(atn.states); i++ {
		var state = atn.states[i]
        if _,ok := state.(*StarLoopEntryState); !ok {
            continue
        }
        // We analyze the ATN to determine if this ATN decision state is the
        // decision for the closure block that determines whether a
        // precedence rule should continue or complete.
        //
        if ( atn.ruleToStartState[state.ruleIndex].isPrecedenceRule) {
            var maybeLoopEndState = state.transitions[len(state.transitions) - 1].target
            if _, ok := maybeLoopEndState.(*LoopEndState); ok {
                s2,ok2 := maybeLoopEndState.transitions[0].target.(*RuleStopState)
                if ( maybeLoopEndState.epsilonOnlyTransitions && ok2) {
                    s2.(*StarLoopEntryState).precedenceRuleDecision = true
                }
            }
        }
	}
}

func (this *ATNDeserializer) verifyATN(atn *ATN) {
    if (!this.deserializationOptions.verifyATN) {
        return
    }
    // verify assumptions
	for i:=0; i<len(atn.states); i++ {

        var state = atn.states[i]
        if (state == nil) {
            continue
        }
        this.checkCondition(state.epsilonOnlyTransitions || len(state.transitions) <= 1, nil)

        switch s2:= state.(type) {

            case *PlusBlockStartState:
                this.checkCondition(s2.loopBackState != nil,nil)
            case *StarLoopEntryState:

                this.checkCondition(s2.loopBackState != nil,nil)
                this.checkCondition(len(s2.transitions) == 2,nil)

                switch _ := s2.(type) {
                    case *StarBlockStartState:
                        _,ok2 := s2.transitions[1].target.(*LoopEndState)
                        this.checkCondition(ok2, nil)
                        this.checkCondition(!s2.nonGreedy, nil)
                    case *LoopEndState:
                        s3,ok2 := s2.transitions[1].target.(*StarBlockStartState)
                        this.checkCondition(ok2, nil)
                        this.checkCondition(s3.nonGreedy, nil)
                    default:
                        panic("IllegalState")
                }

            case *StarLoopbackState:
                this.checkCondition(len(state.transitions) == 1, nil)
                _,ok2 := state.transitions[0].target.(*StarLoopEntryState)
                this.checkCondition(ok2, nil)
            case *LoopEndState:
                this.checkCondition(s2.loopBackState != nil, nil)
            case *RuleStartState:
                this.checkCondition(s2.stopState != nil, nil)
            case *BlockStartState:
                this.checkCondition(s2.endState != nil, nil)
            case *BlockEndState:
                this.checkCondition(s2.startState != nil, nil)
            case *DecisionState:
                this.checkCondition(len(s2.transitions) <= 1 || s2.decision >= 0, nil)
            default:
                _, ok := s2.(*RuleStopState)
                this.checkCondition(len(s2.transitions) <= 1 || ok, nil)
        }
	}
}

func (this *ATNDeserializer) checkCondition(condition bool, message string) {
    if (!condition) {
        if (message==nil) {
            message = "IllegalState"
        }
        panic(message)
    }
}

func (this *ATNDeserializer) readInt() int {
    v := this.data[this.pos]
    this.pos += 1
    return v
}

func (this *ATNDeserializer) readInt32() int {
    var low = this.readInt()
    var high = this.readInt()
    return low | (high << 16)
}

func (this *ATNDeserializer) readLong() int64 {
    var low = this.readInt32()
    var high = this.readInt32()
    return (low & 0x00000000FFFFFFFF) | (high << 32)
}


func createByteToHex() []string {
	var bth = make([]string, 256)
    for i:= 0; i < 256; i++ {
        bth[i] = strings.ToUpper(hex.EncodeToString( []byte{ byte(i) } ))
    }
	return bth
}

var byteToHex = createByteToHex()
	
func (this *ATNDeserializer) readUUID() string {
	var bb = make([]int, 16)
	for  i:=7; i>=0 ;i-- {
		var integer = this.readInt()
		bb[(2*i)+1] = integer & 0xFF
		bb[2*i] = (integer >> 8) & 0xFF
	}
    return byteToHex[bb[0]] + byteToHex[bb[1]] +
        byteToHex[bb[2]] + byteToHex[bb[3]] + '-' +
        byteToHex[bb[4]] + byteToHex[bb[5]] + '-' +
        byteToHex[bb[6]] + byteToHex[bb[7]] + '-' +
        byteToHex[bb[8]] + byteToHex[bb[9]] + '-' +
        byteToHex[bb[10]] + byteToHex[bb[11]] +
        byteToHex[bb[12]] + byteToHex[bb[13]] +
        byteToHex[bb[14]] + byteToHex[bb[15]]
}


func (this *ATNDeserializer) edgeFactory(atn *ATN, typeIndex int, src, trg *ATNState, arg1, arg2, arg3 int, sets []*IntervalSet) *Transition {

    var target = atn.states[trg]

    switch (typeIndex) {
        case TransitionEPSILON :
            return NewEpsilonTransition(target, -1)
        case TransitionRANGE :
            if (arg3 != 0) {
                return NewRangeTransition(target, TokenEOF, arg2)
            } else {
                return NewRangeTransition(target, arg1, arg2)
            }
        case TransitionRULE :
            return NewRuleTransition(atn.states[arg1].(*RuleStartState), arg2, arg3, target)
        case TransitionPREDICATE :
            return NewPredicateTransition(target, arg1, arg2, arg3 != 0)
        case TransitionPRECEDENCE:
            return NewPrecedencePredicateTransition(target, arg1)
        case TransitionATOM :
            if (arg3 != 0) {
                return NewAtomTransition(target, TokenEOF)
            } else {
                return NewAtomTransition(target, arg1)
            }
        case TransitionACTION :
            return NewActionTransition(target, arg1, arg2, arg3 != 0)
        case TransitionSET :
            return NewSetTransition(target, sets[arg1])
        case TransitionNOT_SET :
            return NewNotSetTransition(target, sets[arg1])
        case TransitionWILDCARD :
            return NewWildcardTransition(target)
    }

    panic("The specified transition type is not valid.")
}

func (this *ATNDeserializer) stateFactory(typeIndex, ruleIndex int) *ATNState {
    var s *ATNState
    switch (typeIndex) {
        case ATNStateInvalidType:
            return nil;
        case ATNStateBASIC :
            s = NewBasicState()
        case ATNStateRULE_START :
            s = NewRuleStartState()
        case ATNStateBLOCK_START :
            s = NewBasicBlockStartState()
        case ATNStatePLUS_BLOCK_START :
            s = NewPlusBlockStartState()
        case ATNStateSTAR_BLOCK_START :
            s = NewStarBlockStartState()
        case ATNStateTOKEN_START :
            s = NewTokensStartState()
        case ATNStateRULE_STOP :
            s = NewRuleStopState()
        case ATNStateBLOCK_END :
            s = NewBlockEndState()
        case ATNStateSTAR_LOOP_BACK :
            s = NewStarLoopbackState()
        case ATNStateSTAR_LOOP_ENTRY :
            s = NewStarLoopEntryState()
        case ATNStatePLUS_LOOP_BACK :
            s = NewPlusLoopbackState()
        case ATNStateLOOP_END :
            s = NewLoopEndState()
        default :
            message := fmt.Sprintf("The specified state type %d is not valid.", typeIndex)
            panic(message)
    }

    s.ruleIndex = ruleIndex;
    return s;
}

func (this *ATNDeserializer) lexerActionFactory(typeIndex, data1, data2 int) *LexerAction {
    switch (typeIndex) {
        case LexerActionTypeCHANNEL:
            return NewLexerChannelAction(data1)
        case LexerActionTypeCUSTOM:
            return NewLexerCustomAction(data1, data2)
        case LexerActionTypeMODE:
            return NewLexerModeAction(data1)
        case LexerActionTypeMORE:
            return LexerMoreActionINSTANCE
        case LexerActionTypePOP_MODE:
            return LexerPopModeActionINSTANCE
        case LexerActionTypePUSH_MODE:
            return NewLexerPushModeAction(data1)
        case LexerActionTypeSKIP:
            return LexerSkipActionINSTANCE
        case LexerActionTypeTYPE:
            return NewLexerTypeAction(data1)
        default:
            message := fmt.Sprintf("The specified lexer action typeIndex%d is not valid.", typeIndex)
            panic(message)
    }
    return nil
}
