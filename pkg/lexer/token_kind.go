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

const (
	EOFTokenKind = iota // "\0"

	PlusOpTokenKind  // "+"
	MinusOpTokenKind // "-"
	MulOpTokenKind   // "*"
	DivOpTokenKind   // "/"
	BangOpTokenKind  // "!"

	GTOpTokenKind     // ">"
	GTEOpTokenKind    // ">="
	LTOpTokenKind     // "<"
	AssignOpTokenKind // "="
	LTEOpTokenKind    // "<="
	EQOpTokenKind     // "=="
	NEQOpTokenKind    // "!="

	RShiftOpTokenKind // ">>"
	LShiftOpTokenKind // "<<"
	OROpTokenKind     // "|"
	ANDOpTokenKind    // "&"
	XOROpTokenKind    // "^"
	NOTOpTokenKind    // "~"

	OROROpTokenKind   // "||"
	ANDANDOpTokenKind // "&&"

	PlusEqOpTokenKind  // "+="
	MinusEqOpTokenKind // "-="
	MulEqOpTokenKind   // "*="
	DivEqOpTokenKind   // "/="
	XOREqOpTokenKind   // "^="
	OREqOpTokenKind    // "|="

	OpenParentTokenKind   // "("
	CloseParentTokenKind  // ")"
	OpenBracketTokenKind  // "["
	CloseBracketTokenKind // "]"
	OpenBraceTokenKind    // "{"
	CloseBraceTokenKind   // "}"

	CommaTokenKind   // ""
	DotTokenKind     // "."
	SemiColTokenKind // ";"

	PlusPlusOpTokenKind   // "++"
	MinusMinusOpTokenKind // "--"

	IdentifierTokenKind // "id"
	IntTokenKind        // "number:int"
	FloatTokenKind      // "number:float"
	ImaginaryTokenKind  // "number:imag"
	StringTokenKind     // "string"
	BooleanTokenKind    // "true|false"

	CommentTokenKind // "// comment"

	// Keywords
	BreakKeywordTokenKind
	CaseKeywordTokenKind
	ConstKeywordTokenKind
	ContinueKeywordTokenKind
	DefaultKeywordTokenKind
	ElseKeywordTokenKind
	ForKeywordTokenKind
	FunKeywordTokenKind
	IfKeywordTokenKind
	ImportKeywordTokenKind
	Int16KeywordTokenKind
	Int32KeywordTokenKind
	Int64KeywordTokenKind
	Int8KeywordTokenKind
	NamespaceKeywordTokenKind
	StructKeywordTokenKind
	SwitchKeywordTokenKind
	UInt16KeywordTokenKind
	UInt32KeywordTokenKind
	UInt64KeywordTokenKind
	UInt8KeywordTokenKind
	VarKeywordTokenKind

	InvalidTokenKind
)

var dumpedKinds = map[int]string{
	EOFTokenKind:              "eof",
	PlusOpTokenKind:           "op:plus",
	MinusOpTokenKind:          "op:minus",
	MulOpTokenKind:            "op:mul",
	DivOpTokenKind:            "op:div",
	BangOpTokenKind:           "op:bang",
	GTOpTokenKind:             "op:gt",
	GTEOpTokenKind:            "op:gte",
	LTOpTokenKind:             "op:lt",
	AssignOpTokenKind:         "op:assign",
	LTEOpTokenKind:            "op:lte",
	EQOpTokenKind:             "op:eq",
	NEQOpTokenKind:            "op:neq",
	RShiftOpTokenKind:         "op:rshift",
	LShiftOpTokenKind:         "op:lshift",
	OROpTokenKind:             "op:or",
	ANDOpTokenKind:            "op:and",
	XOROpTokenKind:            "op:xor",
	NOTOpTokenKind:            "op:not",
	OROROpTokenKind:           "op:oror",
	ANDANDOpTokenKind:         "op:andand",
	PlusEqOpTokenKind:         "op:pluseq",
	MinusEqOpTokenKind:        "op:minuseq",
	MulEqOpTokenKind:          "op:muleq",
	DivEqOpTokenKind:          "op:diveq",
	XOREqOpTokenKind:          "op:xoreq",
	OREqOpTokenKind:           "op:oreq",
	OpenParentTokenKind:       "punct:open_parent",
	CloseParentTokenKind:      "punct:close_parent",
	OpenBracketTokenKind:      "punct:open_bracket",
	CloseBracketTokenKind:     "punct:close_bracket",
	OpenBraceTokenKind:        "punct:open_brace",
	CloseBraceTokenKind:       "punct:close_brace",
	CommaTokenKind:            "punct:comma",
	DotTokenKind:              "punct:dot",
	SemiColTokenKind:          "punct:semicol",
	PlusPlusOpTokenKind:       "punct:plusplus",
	MinusMinusOpTokenKind:     "punct:minusminus",
	IdentifierTokenKind:       "id",
	IntTokenKind:              "num.int",
	FloatTokenKind:            "num.float",
	ImaginaryTokenKind:        "num.imag",
	StringTokenKind:           "string",
	BooleanTokenKind:          "boolean",
	CommentTokenKind:          "comment",
	BreakKeywordTokenKind:     "keyword:break",
	CaseKeywordTokenKind:      "keyword:case",
	ConstKeywordTokenKind:     "keyword:const",
	ContinueKeywordTokenKind:  "keyword:continue",
	DefaultKeywordTokenKind:   "keyword:default",
	ElseKeywordTokenKind:      "keyword:else",
	ForKeywordTokenKind:       "keyword:for",
	FunKeywordTokenKind:       "keyword:fun",
	IfKeywordTokenKind:        "keyword:if",
	ImportKeywordTokenKind:    "keyword:import",
	Int16KeywordTokenKind:     "keyword:int16",
	Int32KeywordTokenKind:     "keyword:int32",
	Int64KeywordTokenKind:     "keyword:int64",
	Int8KeywordTokenKind:      "keyword:int8",
	NamespaceKeywordTokenKind: "keyword:namespace",
	StructKeywordTokenKind:    "keyword:struct",
	SwitchKeywordTokenKind:    "keyword:switch",
	UInt16KeywordTokenKind:    "keyword:uint16",
	UInt32KeywordTokenKind:    "keyword:uint32",
	UInt64KeywordTokenKind:    "keyword:uint64",
	UInt8KeywordTokenKind:     "keyword:uint8",
	VarKeywordTokenKind:       "keyword:var",
	InvalidTokenKind:          "invalid",
}

func DumpTokenKind(k int) string {
	return dumpedKinds[k]
}

// Sorted list of keywords (for binary searching them in scanning process)
var keywords = []string{
	"break", "case", "const", "continue", "default", "else",
	"for", "fun", "if", "import", "int16", "int32",
	"int64", "int8", "namespace", "struct", "switch", "uint16",
	"uint32", "uint64", "uint8", "var",
}

var keywordsAmount = len(keywords)
