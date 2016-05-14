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

#pragma once

#include "Vocabulary.h"

namespace org {
namespace antlr {
namespace v4 {
namespace runtime {
namespace dfa {

  /// This class provides a default implementation of the <seealso cref="Vocabulary"/>
  /// interface.
  class ANTLR4CPP_PUBLIC VocabularyImpl : public Vocabulary {
  public:
    virtual ~VocabularyImpl() {};
    
    /// Gets an empty <seealso cref="Vocabulary"/> instance.
    ///
    /// <para>
    /// No literal or symbol names are assigned to token types, so
    /// <seealso cref="#getDisplayName(int)"/> returns the numeric value for all tokens
    /// except <seealso cref="Token#EOF"/>.</para>
    static const Ref<Vocabulary> EMPTY_VOCABULARY;

    /// <summary>
    /// Constructs a new instance of <seealso cref="VocabularyImpl"/> from the specified
    /// literal and symbolic token names.
    /// </summary>
    /// <param name="literalNames"> The literal names assigned to tokens, or {@code null}
    /// if no literal names are assigned. </param>
    /// <param name="symbolicNames"> The symbolic names assigned to tokens, or
    /// {@code null} if no symbolic names are assigned.
    /// </param>
    /// <seealso cref= #getLiteralName(int) </seealso>
    /// <seealso cref= #getSymbolicName(int) </seealso>
    VocabularyImpl(const std::vector<std::string> &literalNames, const std::vector<std::string> &symbolicNames);

    /// <summary>
    /// Constructs a new instance of <seealso cref="VocabularyImpl"/> from the specified
    /// literal, symbolic, and display token names.
    /// </summary>
    /// <param name="literalNames"> The literal names assigned to tokens, or {@code null}
    /// if no literal names are assigned. </param>
    /// <param name="symbolicNames"> The symbolic names assigned to tokens, or
    /// {@code null} if no symbolic names are assigned. </param>
    /// <param name="displayNames"> The display names assigned to tokens, or {@code null}
    /// to use the values in {@code literalNames} and {@code symbolicNames} as
    /// the source of display names, as described in
    /// <seealso cref="#getDisplayName(int)"/>.
    /// </param>
    /// <seealso cref= #getLiteralName(int) </seealso>
    /// <seealso cref= #getSymbolicName(int) </seealso>
    /// <seealso cref= #getDisplayName(int) </seealso>
    VocabularyImpl(const std::vector<std::string> &literalNames, const std::vector<std::string> &symbolicNames,
                   const std::vector<std::string> &displayNames);

    /// <summary>
    /// Returns a <seealso cref="VocabularyImpl"/> instance from the specified set of token
    /// names. This method acts as a compatibility layer for the single
    /// {@code tokenNames} array generated by previous releases of ANTLR.
    ///
    /// <para>The resulting vocabulary instance returns {@code null} for
    /// <seealso cref="#getLiteralName(int)"/> and <seealso cref="#getSymbolicName(int)"/>, and the
    /// value from {@code tokenNames} for the display names.</para>
    /// </summary>
    /// <param name="tokenNames"> The token names, or {@code null} if no token names are
    /// available. </param>
    /// <returns> A <seealso cref="Vocabulary"/> instance which uses {@code tokenNames} for
    /// the display names of tokens. </returns>
    static Ref<Vocabulary> fromTokenNames(const std::vector<std::string> &tokenNames);

    virtual int getMaxTokenType() const override;
    virtual std::string getLiteralName(ssize_t tokenType) const override;
    virtual std::string getSymbolicName(ssize_t tokenType) const override;
    virtual std::string getDisplayName(ssize_t tokenType) const override;

  private:
    static std::vector<std::string> const EMPTY_NAMES;

    std::vector<std::string> const _literalNames;
    std::vector<std::string> const _symbolicNames;
    std::vector<std::string> const _displayNames;
    const int _maxTokenType;
  };
  
} // namespace atn
} // namespace runtime
} // namespace v4
} // namespace antlr
} // namespace org
