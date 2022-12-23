# MIT License
#
# Copyright (c) 2022 Adi Salimgereev
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

import subprocess

keywords_list = ["break", "return", "case", "const", "continue", "default", "else",
                 "for", "fun", "if", "import", "int16", "int32",
                 "int64", "int8", "namespace", "struct", "switch", "uint16",
                 "uint32", "uint64", "uint8", "var"]
keywords_list.sort()

dumped_keywords_list = keywords_list.__str__(
)[1:][:-1].replace("'", "\"") + ","

keywords_list_code = "var keywords = []string{\n\t%s\n}\n\nvar keywordsAmount = len(keywords)" % (
    dumped_keywords_list)

token_kinds_list = """
const (
	EOFTokenKind = iota // "\\0"

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
"""

dumped_kinds_map = """
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
	IntTokenKind:              "num:int",
	FloatTokenKind:            "num:float",
	ImaginaryTokenKind:        "num:imag",
	StringTokenKind:           "string",
	BooleanTokenKind:          "boolean",
	CommentTokenKind:          "comment",
"""

for i in range(len(keywords_list)):
    keyword = keywords_list[i]

    kindName = keyword.title() + "KeywordTokenKind"
    token_kinds_list += "\t" + kindName + "\n"

    dumped_kinds_map_item = "\t" + keyword.title() + "KeywordTokenKind:\t" + \
        f"\"keyword:{keyword}\",\n"
    dumped_kinds_map += dumped_kinds_map_item

token_kinds_list += "\n\tInvalidTokenKind\n)\n"
dumped_kinds_map += """
	InvalidTokenKind:          "invalid",
}"""

obrace = '{'
cbrace = '}'

code = f"""// Auto generated by token_kind_gen.py

package lexer

{token_kinds_list}

{dumped_kinds_map}

func DumpTokenKind(k int) string {obrace}
	return dumpedKinds[k]
{cbrace}

{keywords_list_code}"""

with open("token_kind.go", "w") as f:
    f.write(code)

subprocess.call(["go", "fmt", "token_kind.go"])
