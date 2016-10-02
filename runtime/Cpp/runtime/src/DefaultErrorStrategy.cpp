﻿/*
 * [The "BSD license"]
 *  Copyright (c) 2016 Mike Lischke
 *  Copyright (c) 2013 Terence Parr
 *  Copyright (c) 2013 Dan McLaughlin
 *  All rights reserved.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions
 *  are met:
 *
 *  1. Redistributions of source code must retain the above copyright
 *     notice, this list of conditions and the following disclaimer.
 *  2. Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *  3. The name of the author may not be used to endorse or promote products
 *     derived from this software without specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
 *  IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
 *  OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
 *  IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,
 *  INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
 *  NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 *  DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 *  THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 *  (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
 *  THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

#include "NoViableAltException.h"
#include "misc/IntervalSet.h"
#include "atn/ParserATNSimulator.h"
#include "InputMismatchException.h"
#include "FailedPredicateException.h"
#include "ParserRuleContext.h"
#include "atn/RuleTransition.h"
#include "atn/ATN.h"
#include "atn/ATNState.h"
#include "Parser.h"
#include "CommonToken.h"
#include "Vocabulary.h"
#include "support/StringUtils.h"

#include "DefaultErrorStrategy.h"

using namespace antlr4;
using namespace antlr4::atn;

using namespace antlrcpp;

DefaultErrorStrategy::DefaultErrorStrategy() {
  InitializeInstanceFields();
}

DefaultErrorStrategy::~DefaultErrorStrategy() {
}

void DefaultErrorStrategy::reset(Parser *recognizer) {
  _errorSymbols.clear();
  endErrorCondition(recognizer);
}

void DefaultErrorStrategy::beginErrorCondition(Parser * /*recognizer*/) {
  errorRecoveryMode = true;
}

bool DefaultErrorStrategy::inErrorRecoveryMode(Parser * /*recognizer*/) {
  return errorRecoveryMode;
}

void DefaultErrorStrategy::endErrorCondition(Parser * /*recognizer*/) {
  errorRecoveryMode = false;
  lastErrorIndex = -1;
}

void DefaultErrorStrategy::reportMatch(Parser *recognizer) {
  endErrorCondition(recognizer);
}

void DefaultErrorStrategy::reportError(Parser *recognizer, const RecognitionException &e) {
  // If we've already reported an error and have not matched a token
  // yet successfully, don't report any errors.
  if (inErrorRecoveryMode(recognizer)) {
    return; // don't report spurious errors
  }

  beginErrorCondition(recognizer);
  if (is<const NoViableAltException *>(&e)) {
    reportNoViableAlternative(recognizer, (const NoViableAltException &)e);
  } else if (is<const InputMismatchException *>(&e)) {
    reportInputMismatch(recognizer, (const InputMismatchException &)e);
  } else if (is<const FailedPredicateException *>(&e)) {
    reportFailedPredicate(recognizer, (const FailedPredicateException &)e);
  } else if (is<const RecognitionException *>(&e)) {
    recognizer->notifyErrorListeners(e.getOffendingToken(), e.what(), std::current_exception());
  }
}

void DefaultErrorStrategy::recover(Parser *recognizer, std::exception_ptr /*e*/) {
  if (lastErrorIndex == (int)recognizer->getInputStream()->index() &&
      lastErrorStates.contains(recognizer->getState())) {

    // uh oh, another error at same token index and previously-visited
    // state in ATN; must be a case where LT(1) is in the recovery
    // token set so nothing got consumed. Consume a single token
    // at least to prevent an infinite loop; this is a failsafe.
    recognizer->consume();
  }
  lastErrorIndex = (int)recognizer->getInputStream()->index();
  lastErrorStates.add(recognizer->getState());
  misc::IntervalSet followSet = getErrorRecoverySet(recognizer);
  consumeUntil(recognizer, followSet);
}

void DefaultErrorStrategy::sync(Parser *recognizer) {
  atn::ATNState *s = recognizer->getInterpreter<atn::ATNSimulator>()->atn.states[recognizer->getState()];

  // If already recovering, don't try to sync
  if (inErrorRecoveryMode(recognizer)) {
    return;
  }

  TokenStream *tokens = recognizer->getTokenStream();
  size_t la = tokens->LA(1);

  // try cheaper subset first; might get lucky. seems to shave a wee bit off
  if (recognizer->getATN().nextTokens(s).contains(la) || la == Token::EOF) {
    return;
  }

  // Return but don't end recovery. only do that upon valid token match
  if (recognizer->isExpectedToken((int)la)) {
    return;
  }

  switch (s->getStateType()) {
    case atn::ATNState::BLOCK_START:
    case atn::ATNState::STAR_BLOCK_START:
    case atn::ATNState::PLUS_BLOCK_START:
    case atn::ATNState::STAR_LOOP_ENTRY:
      // report error and recover if possible
      if (singleTokenDeletion(recognizer) != nullptr) {
        return;
      }

      throw InputMismatchException(recognizer);

    case atn::ATNState::PLUS_LOOP_BACK:
    case atn::ATNState::STAR_LOOP_BACK: {
      reportUnwantedToken(recognizer);
      misc::IntervalSet expecting = recognizer->getExpectedTokens();
      misc::IntervalSet whatFollowsLoopIterationOrRule = expecting.Or(getErrorRecoverySet(recognizer));
      consumeUntil(recognizer, whatFollowsLoopIterationOrRule);
    }
      break;

    default:
      // do nothing if we can't identify the exact kind of ATN state
      break;
  }
}

void DefaultErrorStrategy::reportNoViableAlternative(Parser *recognizer, const NoViableAltException &e) {
  TokenStream *tokens = recognizer->getTokenStream();
  std::string input;
  if (tokens != nullptr) {
    if (e.getStartToken()->getType() == Token::EOF) {
      input = "<EOF>";
    } else {
      input = tokens->getText(e.getStartToken(), e.getOffendingToken());
    }
  } else {
    input = "<unknown input>";
  }
  std::string msg = "no viable alternative at input " + escapeWSAndQuote(input);
  recognizer->notifyErrorListeners(e.getOffendingToken(), msg, std::make_exception_ptr(e));
}

void DefaultErrorStrategy::reportInputMismatch(Parser *recognizer, const InputMismatchException &e) {
  std::string msg = "mismatched input " + getTokenErrorDisplay(e.getOffendingToken()) +
  " expecting " + e.getExpectedTokens().toString(recognizer->getVocabulary());
  recognizer->notifyErrorListeners(e.getOffendingToken(), msg, std::make_exception_ptr(e));
}

void DefaultErrorStrategy::reportFailedPredicate(Parser *recognizer, const FailedPredicateException &e) {
  const std::string& ruleName = recognizer->getRuleNames()[recognizer->getContext()->getRuleIndex()];
  std::string msg = "rule " + ruleName + " " + e.what();
  recognizer->notifyErrorListeners(e.getOffendingToken(), msg, std::make_exception_ptr(e));
}

void DefaultErrorStrategy::reportUnwantedToken(Parser *recognizer) {
  if (inErrorRecoveryMode(recognizer)) {
    return;
  }

  beginErrorCondition(recognizer);

  Token *t = recognizer->getCurrentToken();
  std::string tokenName = getTokenErrorDisplay(t);
  misc::IntervalSet expecting = getExpectedTokens(recognizer);

  std::string msg = "extraneous input " + tokenName + " expecting " + expecting.toString(recognizer->getVocabulary());
  recognizer->notifyErrorListeners(t, msg, nullptr);
}

void DefaultErrorStrategy::reportMissingToken(Parser *recognizer) {
  if (inErrorRecoveryMode(recognizer)) {
    return;
  }

  beginErrorCondition(recognizer);

  Token *t = recognizer->getCurrentToken();
  misc::IntervalSet expecting = getExpectedTokens(recognizer);
  std::string expectedText = expecting.toString(recognizer->getVocabulary());
  std::string msg = "missing " + expectedText + " at " + getTokenErrorDisplay(t);

  recognizer->notifyErrorListeners(t, msg, nullptr);
}

Token* DefaultErrorStrategy::recoverInline(Parser *recognizer) {
  // Single token deletion.
  Token *matchedSymbol = singleTokenDeletion(recognizer);
  if (matchedSymbol) {
    // We have deleted the extra token.
    // Now, move past ttype token as if all were ok.
    recognizer->consume();
    return matchedSymbol;
  }

  // Single token insertion.
  if (singleTokenInsertion(recognizer)) {
    return getMissingSymbol(recognizer);
  }

  // Even that didn't work; must throw the exception.
  throw InputMismatchException(recognizer);
}

bool DefaultErrorStrategy::singleTokenInsertion(Parser *recognizer) {
  ssize_t currentSymbolType = recognizer->getInputStream()->LA(1);

  // if current token is consistent with what could come after current
  // ATN state, then we know we're missing a token; error recovery
  // is free to conjure up and insert the missing token
  atn::ATNState *currentState = recognizer->getInterpreter<atn::ATNSimulator>()->atn.states[recognizer->getState()];
  atn::ATNState *next = currentState->transition(0)->target;
  const atn::ATN &atn = recognizer->getInterpreter<atn::ATNSimulator>()->atn;
  misc::IntervalSet expectingAtLL2 = atn.nextTokens(next, recognizer->getContext());
  if (expectingAtLL2.contains(currentSymbolType)) {
    reportMissingToken(recognizer);
    return true;
  }
  return false;
}

Token* DefaultErrorStrategy::singleTokenDeletion(Parser *recognizer) {
  size_t nextTokenType = recognizer->getInputStream()->LA(2);
  misc::IntervalSet expecting = getExpectedTokens(recognizer);
  if (expecting.contains(nextTokenType)) {
    reportUnwantedToken(recognizer);
    recognizer->consume(); // simply delete extra token
                           // we want to return the token we're actually matching
    Token *matchedSymbol = recognizer->getCurrentToken();
    reportMatch(recognizer); // we know current token is correct
    return matchedSymbol;
  }
  return nullptr;
}

Token* DefaultErrorStrategy::getMissingSymbol(Parser *recognizer) {
  Token *currentSymbol = recognizer->getCurrentToken();
  misc::IntervalSet expecting = getExpectedTokens(recognizer);
  size_t expectedTokenType = expecting.getMinElement(); // get any element
  std::string tokenText;
  if (expectedTokenType == Token::EOF) {
    tokenText = "<missing EOF>";
  } else {
    tokenText = "<missing " + recognizer->getVocabulary().getDisplayName(expectedTokenType) + ">";
  }
  Token *current = currentSymbol;
  Token *lookback = recognizer->getTokenStream()->LT(-1);
  if (current->getType() == Token::EOF && lookback != nullptr) {
    current = lookback;
  }

  _errorSymbols.push_back(recognizer->getTokenFactory()->create(
    { current->getTokenSource(), current->getTokenSource()->getInputStream() },
    (int)expectedTokenType, tokenText, Token::DEFAULT_CHANNEL, -1, -1,
    current->getLine(), current->getCharPositionInLine()));
  
  return _errorSymbols.back().get();
}

misc::IntervalSet DefaultErrorStrategy::getExpectedTokens(Parser *recognizer) {
  return recognizer->getExpectedTokens();
}

std::string DefaultErrorStrategy::getTokenErrorDisplay(Token *t) {
  if (t == nullptr) {
    return "<no Token>";
  }
  std::string s = getSymbolText(t);
  if (s == "") {
    if (getSymbolType(t) == Token::EOF) {
      s = "<EOF>";
    } else {
      s = "<" + std::to_string(getSymbolType(t)) + ">";
    }
  }
  return escapeWSAndQuote(s);
}

std::string DefaultErrorStrategy::getSymbolText(Token *symbol) {
  return symbol->getText();
}

size_t DefaultErrorStrategy::getSymbolType(Token *symbol) {
  return symbol->getType();
}

std::string DefaultErrorStrategy::escapeWSAndQuote(const std::string &s) const {
  std::string result = s;
  antlrcpp::replaceAll(result, "\n", "\\n");
  antlrcpp::replaceAll(result, "\r","\\r");
  antlrcpp::replaceAll(result, "\t","\\t");
  return "'" + result + "'";
}

misc::IntervalSet DefaultErrorStrategy::getErrorRecoverySet(Parser *recognizer) {
  const atn::ATN &atn = recognizer->getInterpreter<atn::ATNSimulator>()->atn;
  RuleContext *ctx = recognizer->getContext().get();
  misc::IntervalSet recoverSet;
  while (ctx->invokingState != ATNState::INVALID_STATE_NUMBER) {
    // compute what follows who invoked us
    atn::ATNState *invokingState = atn.states[ctx->invokingState];
    atn::RuleTransition *rt = dynamic_cast<atn::RuleTransition*>(invokingState->transition(0));
    misc::IntervalSet follow = atn.nextTokens(rt->followState);
    recoverSet.addAll(follow);

    if (ctx->parent.expired())
      break;
    ctx = (RuleContext *)ctx->parent.lock().get();
  }
  recoverSet.remove(Token::EPSILON);

  return recoverSet;
}

void DefaultErrorStrategy::consumeUntil(Parser *recognizer, const misc::IntervalSet &set) {
  size_t ttype = recognizer->getInputStream()->LA(1);
  while (ttype != Token::EOF && !set.contains(ttype)) {
    recognizer->consume();
    ttype = recognizer->getInputStream()->LA(1);
  }
}

void DefaultErrorStrategy::InitializeInstanceFields() {
  errorRecoveryMode = false;
  lastErrorIndex = -1;
}
