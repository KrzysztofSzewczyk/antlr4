/*
 * [The "BSD license"]
 *  Copyright (c) 2016 Mike Lischke
 *  Copyright (c) 2014 Terence Parr
 *  Copyright (c) 2014 Dan McLaughlin
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

#include "VocabularyImpl.h"

using namespace org::antlr::v4::runtime::dfa;

const std::vector<std::wstring> VocabularyImpl::EMPTY_NAMES;
const Ref<Vocabulary> VocabularyImpl::EMPTY_VOCABULARY = std::make_shared<VocabularyImpl>(EMPTY_NAMES, EMPTY_NAMES, EMPTY_NAMES);

VocabularyImpl::VocabularyImpl(const std::vector<std::wstring> &literalNames, const std::vector<std::wstring> &symbolicNames)
: VocabularyImpl(literalNames, symbolicNames, {}) {
}

VocabularyImpl::VocabularyImpl(const std::vector<std::wstring> &literalNames,
  const std::vector<std::wstring> &symbolicNames, const std::vector<std::wstring> &displayNames)
  : _literalNames(!literalNames.empty() ? literalNames : EMPTY_NAMES),
    _symbolicNames(!symbolicNames.empty() ? symbolicNames : EMPTY_NAMES),
    _displayNames(!displayNames.empty() ? displayNames : EMPTY_NAMES),
    _maxTokenType(std::max((int)_displayNames.size(), std::max((int)_literalNames.size(), (int)_symbolicNames.size())) - 1) {
  // See note here on -1 part: https://github.com/antlr/antlr4/pull/1146
}

Ref<Vocabulary> VocabularyImpl::fromTokenNames(const std::vector<std::wstring> &tokenNames) {
  if (tokenNames.empty()) {
    return EMPTY_VOCABULARY;
  }

  std::vector<std::wstring> literalNames = tokenNames;
  std::vector<std::wstring> symbolicNames = tokenNames;
  for (size_t i = 0; i < tokenNames.size(); i++) {
    std::wstring tokenName = tokenNames[i];
    if (tokenName == L"") {
      continue;
    }

    if (!tokenName.empty()) {
      wchar_t firstChar = tokenName[0];
      if (firstChar == L'\'') {
        symbolicNames[i] = L"";
        continue;
      } else if (std::isupper(firstChar)) {
        literalNames[i] = L"";
        continue;
      }
    }

    // wasn't a literal or symbolic name
    literalNames[i] = L"";
    symbolicNames[i] = L"";
  }

  return std::make_shared<VocabularyImpl>(literalNames, symbolicNames, tokenNames);
}

int VocabularyImpl::getMaxTokenType() const {
  return _maxTokenType;
}

std::wstring VocabularyImpl::getLiteralName(ssize_t tokenType) const {
  if (tokenType >= 0 && tokenType < (int)_literalNames.size()) {
    return _literalNames[tokenType];
  }

  return L"";
}

std::wstring VocabularyImpl::getSymbolicName(ssize_t tokenType) const {
  if (tokenType >= 0 && tokenType < (int)_symbolicNames.size()) {
    return _symbolicNames[tokenType];
  }

  if (tokenType == EOF) {
    return L"EOF";
  }

  return L"";
}

std::wstring VocabularyImpl::getDisplayName(ssize_t tokenType) const {
  if (tokenType >= 0 && tokenType < (int)_displayNames.size()) {
    std::wstring displayName = _displayNames[tokenType];
    if (!displayName.empty()) {
      return displayName;
    }
  }

  std::wstring literalName = getLiteralName(tokenType);
  if (!literalName.empty()) {
    return literalName;
  }

  std::wstring symbolicName = getSymbolicName(tokenType);
  if (!symbolicName.empty()) {
    return symbolicName;
  }

  return std::to_wstring(tokenType);
}
