// Auto generated by token_kind_gen.py

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
	PubKeywordTokenKind
	ReturnKeywordTokenKind
	StructKeywordTokenKind
	SwitchKeywordTokenKind
	Uint16KeywordTokenKind
	Uint32KeywordTokenKind
	Uint64KeywordTokenKind
	Uint8KeywordTokenKind
	VarKeywordTokenKind

	InvalidTokenKind
)

var dumpedKinds = map[int]string{
	EOFTokenKind:              "eof",
	PlusOpTokenKind:           "plus",
	MinusOpTokenKind:          "minus",
	MulOpTokenKind:            "asterisk",
	DivOpTokenKind:            "slash",
	BangOpTokenKind:           "bang",
	GTOpTokenKind:             "greater than",
	GTEOpTokenKind:            "greater than or equal",
	LTOpTokenKind:             "less than",
	AssignOpTokenKind:         "assign",
	LTEOpTokenKind:            "less than or equal",
	EQOpTokenKind:             "equal",
	NEQOpTokenKind:            "not equal",
	RShiftOpTokenKind:         "right shift",
	LShiftOpTokenKind:         "left shift",
	OROpTokenKind:             "or",
	ANDOpTokenKind:            "and",
	XOROpTokenKind:            "xor",
	NOTOpTokenKind:            "not",
	OROROpTokenKind:           "or or",
	ANDANDOpTokenKind:         "and and",
	PlusEqOpTokenKind:         "plus equal",
	MinusEqOpTokenKind:        "minus equal",
	MulEqOpTokenKind:          "asterisk equal",
	DivEqOpTokenKind:          "slash equal",
	XOREqOpTokenKind:          "xor equal",
	OREqOpTokenKind:           "or equal",
	OpenParentTokenKind:       "open parent",
	CloseParentTokenKind:      "close parent",
	OpenBracketTokenKind:      "open bracket",
	CloseBracketTokenKind:     "close bracket",
	OpenBraceTokenKind:        "open brace",
	CloseBraceTokenKind:       "close brace",
	CommaTokenKind:            "comma",
	DotTokenKind:              "dot",
	SemiColTokenKind:          "semicolon",
	PlusPlusOpTokenKind:       "plus plus",
	MinusMinusOpTokenKind:     "minus minus",
	IdentifierTokenKind:       "identifier",
	IntTokenKind:              "integer",
	FloatTokenKind:            "float",
	ImaginaryTokenKind:        "imaginary number",
	StringTokenKind:           "string",
	BooleanTokenKind:          "boolean",
	CommentTokenKind:          "comment",
	BreakKeywordTokenKind:     "break keyword",
	CaseKeywordTokenKind:      "case keyword",
	ConstKeywordTokenKind:     "const keyword",
	ContinueKeywordTokenKind:  "continue keyword",
	DefaultKeywordTokenKind:   "default keyword",
	ElseKeywordTokenKind:      "else keyword",
	ForKeywordTokenKind:       "for keyword",
	FunKeywordTokenKind:       "fun keyword",
	IfKeywordTokenKind:        "if keyword",
	ImportKeywordTokenKind:    "import keyword",
	Int16KeywordTokenKind:     "int16 keyword",
	Int32KeywordTokenKind:     "int32 keyword",
	Int64KeywordTokenKind:     "int64 keyword",
	Int8KeywordTokenKind:      "int8 keyword",
	NamespaceKeywordTokenKind: "namespace keyword",
	PubKeywordTokenKind:       "pub keyword",
	ReturnKeywordTokenKind:    "return keyword",
	StructKeywordTokenKind:    "struct keyword",
	SwitchKeywordTokenKind:    "switch keyword",
	Uint16KeywordTokenKind:    "uint16 keyword",
	Uint32KeywordTokenKind:    "uint32 keyword",
	Uint64KeywordTokenKind:    "uint64 keyword",
	Uint8KeywordTokenKind:     "uint8 keyword",
	VarKeywordTokenKind:       "var keyword",

	InvalidTokenKind: "invalid token",
}

func DumpTokenKind(k int) string {
	return dumpedKinds[k]
}

var keywords = []string{
	"break", "case", "const", "continue", "default", "else", "for", "fun", "if", "import", "int16", "int32", "int64", "int8", "namespace", "pub", "return", "struct", "switch", "uint16", "uint32", "uint64", "uint8", "var",
}

var keywordsAmount = len(keywords)
