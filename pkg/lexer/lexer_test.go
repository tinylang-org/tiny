// MIT License
//
// Copyright (c) 2022 Adi Salimgereev
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vertexgmd/tinylang/pkg/utils"
)

func TestEOF(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	l := NewLexer("", []byte(""), p)
	tok := l.NextToken()
	assert.Equal(t, tok.Kind, EOFTokenKind)
}

func TestEOF2(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	l := NewLexer(" ", []byte(""), p)
	tok := l.NextToken()
	assert.Equal(t, tok.Kind, EOFTokenKind)
}

func TestASCIIIdentifier(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	l := NewLexer(" ", []byte("test"), p)
	tok := l.NextToken()
	assert.Equal(t, tok.Kind, IdentifierTokenKind)
}

func TestUTF8Identifier(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	l := NewLexer(" ", []byte("привет"), p)
	tok := l.NextToken()
	assert.Equal(t, tok.Kind, IdentifierTokenKind)
}

func TestBoolean(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	l := NewLexer(" ", []byte("true false"), p)
	tok := l.NextToken()
	assert.Equal(t, tok.Kind, BooleanTokenKind)

	tok = l.NextToken()
	assert.Equal(t, tok.Kind, BooleanTokenKind)
}

func TestWrappedIdentifier(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	l := NewLexer(" ", []byte("`test`"), p)
	tok := l.NextToken()
	assert.Equal(t, tok.Kind, IdentifierTokenKind)
	assert.Equal(t, tok.Literal, "`test`")
}

func TestWrappedIdentifier2(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	l := NewLexer(" ", []byte("`hello world`"), p)
	tok := l.NextToken()
	assert.Equal(t, tok.Kind, IdentifierTokenKind)
	assert.Equal(t, tok.Literal, "`hello world`")
}

func TestComment(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	l := NewLexer(" ", []byte("//test"), p)
	tok := l.NextToken()
	assert.Equal(t, tok.Kind, CommentTokenKind)
	assert.Equal(t, tok.Literal, "test")
}

func TestMultiLineComment(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	l := NewLexer(" ", []byte("/*test*/"), p)
	tok := l.NextToken()
	assert.Equal(t, tok.Kind, CommentTokenKind)
	assert.Equal(t, tok.Literal, "test")
}

func OneCharacterTokenTests(t *testing.T) {
	tests := map[string]int{
		"+": PlusOpTokenKind,
		"-": MinusOpTokenKind,
		"*": MulOpTokenKind,
		"/": DivOpTokenKind,
		"!": BangOpTokenKind,

		">": GTOpTokenKind,
		"<": LTOpTokenKind,
		"=": AssignOpTokenKind,

		"|": OROpTokenKind,
		"&": ANDOpTokenKind,
		"^": XOROpTokenKind,
		"~": NOTOpTokenKind,

		",": CommaTokenKind,
		".": DotTokenKind,
		";": SemiColTokenKind,

		"(": OpenParentTokenKind,
		")": CloseParentTokenKind,
		"[": OpenBracketTokenKind,
		"]": CloseBracketTokenKind,
		"{": OpenBraceTokenKind,
		"}": CloseBraceTokenKind,

		"": EOFTokenKind,
	}

	for input, output := range tests {
		d := utils.NewCodeProblemHandler()
		l := NewLexer("", []byte(input), d)

		assert.Equal(t, output, l.NextToken().Kind)
	}
}

func DoubleCharacterTokenTests(t *testing.T) {
	tests := map[string]int{
		">=": GTEOpTokenKind,
		"<=": LTEOpTokenKind,
		"!=": NEQOpTokenKind,

		">>": RShiftOpTokenKind,
		"<<": LShiftOpTokenKind,
		"||": OROROpTokenKind,
		"&&": ANDANDOpTokenKind,

		"+=": PlusEqOpTokenKind,
		"-=": MinusEqOpTokenKind,
		"*=": MulEqOpTokenKind,
		"/=": DivEqOpTokenKind,
		"^=": XOREqOpTokenKind,
		"|=": OREqOpTokenKind,

		"++": PlusPlusOpTokenKind,
		"--": MinusMinusOpTokenKind,
	}

	for input, output := range tests {
		d := utils.NewCodeProblemHandler()
		l := NewLexer("", []byte(input), d)

		assert.Equal(t, output, l.NextToken().Kind)
	}
}

func TestKeywords(t *testing.T) {
	tests := map[string][][]interface{}{
		"fun struct break default case if else switch var const continue for namespace import": {
			{FunKeywordTokenKind, "fun"},
			{StructKeywordTokenKind, "struct"},
			{BreakKeywordTokenKind, "break"},
			{DefaultKeywordTokenKind, "default"},
			{CaseKeywordTokenKind, "case"},
			{IfKeywordTokenKind, "if"},
			{ElseKeywordTokenKind, "else"},
			{SwitchKeywordTokenKind, "switch"},
			{VarKeywordTokenKind, "var"},
			{ConstKeywordTokenKind, "const"},
			{ContinueKeywordTokenKind, "continue"},
			{ForKeywordTokenKind, "for"},
			{NamespaceKeywordTokenKind, "namespace"},
			{ImportKeywordTokenKind, "import"},
			{EOFTokenKind, "\\0"},
		},
		"int8 int16 int32 int64 uint8 uint16 uint32 uint64": {
			{Int8KeywordTokenKind, "int8"},
			{Int16KeywordTokenKind, "int16"},
			{Int32KeywordTokenKind, "int32"},
			{Int64KeywordTokenKind, "int64"},
			{UInt8KeywordTokenKind, "uint8"},
			{UInt16KeywordTokenKind, "uint16"},
			{UInt32KeywordTokenKind, "uint32"},
			{UInt64KeywordTokenKind, "uint64"},
			{EOFTokenKind, "\\0"},
		},
	}

	for input, output := range tests {
		d := utils.NewCodeProblemHandler()
		l := NewLexer("", []byte(input), d)

		var toks []*Token
		for {
			tok := l.NextToken()
			toks = append(toks, tok)

			if tok.Kind == EOFTokenKind {
				break
			}
		}

		for i, tok := range toks {
			assert.Equal(t, output[i][0], tok.Kind)
			assert.Equal(t, output[i][1], tok.Literal)
		}
	}
}

func TestString(t *testing.T) {

	p := utils.NewCodeProblemHandler()
	l := NewLexer(" ", []byte("\"hello world \\u2!\\n\""), p)
	tok := l.NextToken()

	assert.Equal(t, tok.Kind, StringTokenKind)
	p.PrintProblems()
}

func TestNumber(t *testing.T) {
	tests := map[string][][]interface{}{
		"1234567890 0xcafebabe 0. .0 3.14159265 1e+17 2.71828e-1000 0i 1i 123456789012345678890i": {
			{IntTokenKind, "1234567890"},
			{IntTokenKind, "0xcafebabe"},
			{FloatTokenKind, "0."},
			{FloatTokenKind, ".0"},
			{FloatTokenKind, "3.14159265"},
			{FloatTokenKind, "1e+17"},
			{FloatTokenKind, "2.71828e-1000"},
			{ImaginaryTokenKind, "0i"},
			{ImaginaryTokenKind, "1i"},
			{ImaginaryTokenKind, "123456789012345678890i"},
			{EOFTokenKind, "\\0"},
		},
		"0.i .0i 1e0i": {
			{ImaginaryTokenKind, "0.i"},
			{ImaginaryTokenKind, ".0i"},
			{ImaginaryTokenKind, "1e0i"},
			{EOFTokenKind, "\\0"},
		},
		"2.17i 2.71828e-33i": {
			{ImaginaryTokenKind, "2.17i"},
			{ImaginaryTokenKind, "2.71828e-33i"},
			{EOFTokenKind, "\\0"},
		},
	}

	for input, output := range tests {
		d := utils.NewCodeProblemHandler()
		l := NewLexer("", []byte(input), d)

		var toks []*Token
		for {
			tok := l.NextToken()
			toks = append(toks, tok)

			if tok.Kind == EOFTokenKind {
				break
			}
		}

		for i, tok := range toks {
			assert.Equal(t, output[i][0], tok.Kind)
			assert.Equal(t, output[i][1], tok.Literal)
		}
	}
}
