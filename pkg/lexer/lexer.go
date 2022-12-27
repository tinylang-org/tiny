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
	"unicode"
	"unicode/utf8"

	"github.com/tinylang-org/tiny/pkg/utils"
)

type Lexer struct {
	filepath     string
	source       []byte
	sourceLength int

	currentCodePoint     rune
	currentCodePointSize int
	currentLocation      *utils.CodePointLocation

	problemHandler *utils.CodeProblemHandler
}

func NewLexer(filepath string, source []byte, problemHandler *utils.CodeProblemHandler) *Lexer {
	l := &Lexer{
		filepath:       filepath,
		source:         source,
		sourceLength:   len(source),
		problemHandler: problemHandler,
		currentLocation: &utils.CodePointLocation{
			Filepath: filepath,
			Index:    0,
			Line:     1,
			Column:   0,
		},
	}

	l.decodeRune()

	return l
}

// Decode unicode codepoint from utf8 source code text
func (l *Lexer) decodeRune() {
	offset := 1

	if l.currentLocation.Index < l.sourceLength {
		r := rune(l.source[l.currentLocation.Index])

		switch {
		case r == 0:
			l.problemHandler.AddCodeProblem(
				utils.NewError(
					utils.NewOneCodePointBlockLocation(l.currentLocation),
					utils.IllegalNullCharacterErr, []interface{}{}))
		case r >= utf8.RuneSelf:
			r, offset = utf8.DecodeRune(l.source[l.currentLocation.Index:])
			if r == utf8.RuneError && offset == 1 {
				l.problemHandler.AddCodeProblem(
					utils.NewError(
						utils.NewOneCodePointBlockLocation(l.currentLocation),
						utils.IllegalUTF8EncodingErr, []interface{}{}))
			}
		}

		l.currentCodePoint, l.currentCodePointSize = r, offset
	} else {
		l.currentCodePoint, l.currentCodePointSize = -1, 1
	}
}

// Advance lexer state
func (l *Lexer) advance() {
	if l.currentCodePoint == '\n' {
		l.currentLocation.Line++
		l.currentLocation.Column = 0
	} else {
		l.currentLocation.Column++
	}

	l.currentLocation.Index += l.currentCodePointSize

	l.decodeRune()
}

func (l *Lexer) skipWhitespaces() {
	for l.currentCodePoint == ' ' ||
		l.currentCodePoint == '\t' ||
		l.currentCodePoint == '\n' ||
		l.currentCodePoint == '\r' {
		l.advance()
	}
}

func (l *Lexer) characterToken(kind int, literal string) *Token {
	return &Token{Kind: kind, Literal: literal,
		Location: utils.NewOneCodePointBlockLocation(l.currentLocation)}
}

func (l *Lexer) doubleCharacterToken(kind int, literal string) *Token {
	return &Token{Kind: kind, Literal: literal,
		Location: utils.NewTwoCodePointsBlockLocation(l.currentLocation)}
}

func lower(cp rune) rune     { return ('a' - 'A') | cp }
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }

func isIdentifierContinue(cp rune) bool {
	return 'a' <= lower(cp) && lower(cp) <= 'z' || cp == '_' || '0' <= cp && cp <= '9' ||
		cp >= utf8.RuneSelf && unicode.IsLetter(cp)
}

func isIdentifierStart(cp rune) bool {
	return 'a' <= lower(cp) && lower(cp) <= 'z' || cp == '_' ||
		cp >= utf8.RuneSelf && unicode.IsLetter(cp)
}

func binarySearchKeyword(keyword string) int {
	leftBound := 0
	rightBound := keywordsAmount - 1

	for leftBound <= rightBound {
		midPoint := leftBound + (rightBound-leftBound)/2

		if keywords[midPoint] == keyword {
			return midPoint
		}

		if keywords[midPoint] > keyword {
			rightBound = midPoint - 1
		} else {
			leftBound = midPoint + 1
		}
	}

	return -1
}

func (l *Lexer) peekByte() byte {
	if l.currentLocation.Index+1 < l.sourceLength {
		return l.source[l.currentLocation.Index+1]
	}
	return 0
}

// Number scanning implementation is mostly
// from https://github.com/golang/go/blob/78fc81070a/src/go/scanner/scanner.go
func invalidSeparator(x string) int {
	x1 := ' ' // prefix char, we only care if it's 'x'
	d := '.'  // digit, one of '_', '0' (a digit), or '.' (anything else)
	i := 0

	// a prefix counts as a digit
	if len(x) >= 2 && x[0] == '0' {
		x1 = lower(rune(x[1]))
		if x1 == 'x' || x1 == 'o' || x1 == 'b' {
			d = '0'
			i = 2
		}
	}

	// mantissa and exponent
	for ; i < len(x); i++ {
		p := d // previous digit
		d = rune(x[i])
		switch {
		case d == '_':
			if p != '0' {
				return i
			}
		case isDecimal(d) || x1 == 'x' && isHex(d):
			d = '0'
		default:
			if p == '_' {
				return i - 1
			}
			d = '.'
		}
	}
	if d == '_' {
		return len(x) - 1
	}

	return -1
}

func (l *Lexer) digits(base int, invalid *int) (digitSeparator int) {
	if base <= 10 {
		max := rune('0' + base)
		for isDecimal(l.currentCodePoint) || l.currentCodePoint == '_' {
			ds := 1
			if l.currentCodePoint == '_' {
				ds = 2
			} else if l.currentCodePoint >= max && *invalid < 0 {
				*invalid = l.currentLocation.Index // record invalid rune offset
			}
			digitSeparator |= ds
			l.advance()
		}
	} else {
		for isHex(l.currentCodePoint) || l.currentCodePoint == '_' {
			ds := 1
			if l.currentCodePoint == '_' {
				ds = 2
			}
			digitSeparator |= ds
			l.advance()
		}
	}
	return
}

func (l *Lexer) nextNumberToken() *Token {
	startLocation := l.currentLocation.Copy()
	tokenKind := InvalidTokenKind

	base := 10          // number base
	prefix := rune(0)   // one of 0 (decimal), '0' (0-octal), 'x', 'o', or 'b'
	digitSeparator := 0 // bit 0: digit present, bit 1: '_' present

	// location of invalid digit in literal, or < 0
	invalid := -1

	// integer part
	if l.currentCodePoint != '.' {
		tokenKind = IntTokenKind
		if l.currentCodePoint == '0' {
			l.advance()
			switch lower(l.currentCodePoint) {
			case 'x':
				l.advance()
				base, prefix = 16, 'x'
			case 'o':
				l.advance()
				base, prefix = 8, 'o'
			case 'b':
				l.advance()
				base, prefix = 2, 'b'
			default:
				base, prefix = 8, '0'
				digitSeparator = 1 // leading 0
			}
		}
		digitSeparator |= l.digits(base, &invalid)
	}

	// fractional part
	if l.currentCodePoint == '.' {
		tokenKind = FloatTokenKind
		if prefix == 'o' || prefix == 'b' {
			l.problemHandler.AddCodeProblem(
				utils.NewError(
					&utils.CodeBlockLocation{StartLocation: startLocation,
						EndLocation: l.currentLocation.Copy()},
					utils.InvalidRadixPointErr, []interface{}{numberLiteralName(prefix)}))
		}
		l.advance()
		digitSeparator |= l.digits(base, &invalid)
	}

	if digitSeparator&1 == 0 {
		l.problemHandler.AddCodeProblem(
			utils.NewError(
				&utils.CodeBlockLocation{StartLocation: startLocation,
					EndLocation: l.currentLocation.Copy()},
				utils.HasNoDigitsErr, []interface{}{numberLiteralName(prefix)}))
	}

	// exponent
	if e := lower(l.currentCodePoint); e == 'e' || e == 'p' {
		switch {
		case e == 'e' && prefix != 0 && prefix != '0':
			l.problemHandler.AddCodeProblem(
				utils.NewError(
					&utils.CodeBlockLocation{StartLocation: startLocation,
						EndLocation: l.currentLocation.Copy()},
					utils.ExponentRequiresDecimalMantissaErr,
					[]interface{}{l.currentCodePoint}))
		case e == 'p' && prefix != 'x':
			l.problemHandler.AddCodeProblem(
				utils.NewError(
					&utils.CodeBlockLocation{StartLocation: startLocation,
						EndLocation: l.currentLocation.Copy()},
					utils.ExponentRequiresHexadecimalMantissaErr,
					[]interface{}{l.currentCodePoint}))
		}
		l.advance()
		tokenKind = FloatTokenKind
		if l.currentCodePoint == '+' || l.currentCodePoint == '-' {
			l.advance()
		}
		ds := l.digits(10, nil)
		digitSeparator |= ds
		if ds&1 == 0 {
			l.problemHandler.AddCodeProblem(
				utils.NewError(
					&utils.CodeBlockLocation{StartLocation: startLocation,
						EndLocation: l.currentLocation.Copy()},
					utils.ExponentHasNoDigitsErr,
					[]interface{}{}))
		}
	} else if prefix == 'x' && tokenKind == FloatTokenKind {
		l.problemHandler.AddCodeProblem(
			utils.NewError(
				&utils.CodeBlockLocation{StartLocation: startLocation,
					EndLocation: l.currentLocation.Copy()},
				utils.HexadecimalMantissaRequiresPExponentErr,
				[]interface{}{}))
	}

	// suffix 'i'
	if l.currentCodePoint == 'i' {
		tokenKind = ImaginaryTokenKind
		l.advance()
	}

	buffer := string(l.source[startLocation.Index:l.currentLocation.Index])
	if tokenKind == IntTokenKind && invalid >= 0 {
		l.problemHandler.AddCodeProblem(
			utils.NewError(
				&utils.CodeBlockLocation{StartLocation: startLocation,
					EndLocation: l.currentLocation.Copy()},
				utils.InvalidDigitErr,
				[]interface{}{buffer[invalid-startLocation.Index], numberLiteralName(prefix)}))
	}
	if digitSeparator&2 != 0 {
		if i := invalidSeparator(buffer); i >= 0 {
			l.problemHandler.AddCodeProblem(
				utils.NewError(
					&utils.CodeBlockLocation{StartLocation: startLocation,
						EndLocation: l.currentLocation.Copy()},
					utils.UnderscoreMustSeparateSuccessiveDigitsErr,
					[]interface{}{}))
		}
	}

	return &Token{Kind: tokenKind, Literal: buffer,
		Location: &utils.CodeBlockLocation{StartLocation: startLocation,
			EndLocation: l.currentLocation.Copy()}}
}

func numberLiteralName(prefix rune) string {
	switch prefix {
	case 'x':
		return "hexadecimal literal"
	case 'o', '0':
		return "octal literal"
	case 'b':
		return "binary literal"
	}
	return "decimal literal"
}

func (l *Lexer) nextCommentToken() *Token {
	startLocation := l.currentLocation.Copy()

	l.advance()
	l.advance() // skip '/' twice
	for l.currentCodePoint != '\n' && l.currentCodePoint != -1 {
		l.advance()
	}

	buffer := string(l.source[startLocation.Index+2 : l.currentLocation.Index])
	return &Token{Kind: CommentTokenKind, Literal: buffer,
		Location: &utils.CodeBlockLocation{StartLocation: startLocation,
			EndLocation: l.currentLocation.Copy()}}
}

func (l *Lexer) nextMultiLineCommentToken() *Token {
	startLocation := l.currentLocation.Copy()

	l.advance()
	l.advance() // skip '/' and '*'
	for {
		if l.currentCodePoint == '*' && l.peekByte() == '/' {
			break
		}

		if l.currentCodePoint == -1 {
			l.problemHandler.AddCodeProblem(
				utils.NewError(
					utils.NewOneCodePointBlockLocation(l.currentLocation),
					utils.NotClosedMultiLineCommentErr, []interface{}{l.currentCodePoint}))
			return &Token{Kind: CommentTokenKind,
				Literal: string(l.source[startLocation.Index+2 : l.currentLocation.Index]),
				Location: &utils.CodeBlockLocation{StartLocation: startLocation,
					EndLocation: l.currentLocation.Copy()}}
		}

		l.advance()
	}

	l.advance()
	l.advance() // skip '*' and '/'

	buffer := string(l.source[startLocation.Index+2 : l.currentLocation.Index-2])
	return &Token{Kind: CommentTokenKind, Literal: buffer,
		Location: &utils.CodeBlockLocation{StartLocation: startLocation,
			EndLocation: l.currentLocation.Copy()}}
}

func (l *Lexer) nextNameToken() *Token {
	startLocation := l.currentLocation.Copy()
	for isIdentifierContinue(l.currentCodePoint) {
		l.advance()
	}

	buffer := string(l.source[startLocation.Index:l.currentLocation.Index])

	if buffer == "true" || buffer == "false" {
		return &Token{Kind: BooleanTokenKind, Literal: buffer,
			Location: &utils.CodeBlockLocation{StartLocation: startLocation,
				EndLocation: l.currentLocation.Copy()}}
	}

	keywordIndex := binarySearchKeyword(buffer)
	if keywordIndex != -1 {
		return &Token{Kind: BreakKeywordTokenKind + keywordIndex, Literal: buffer,
			Location: &utils.CodeBlockLocation{StartLocation: startLocation,
				EndLocation: l.currentLocation.Copy()}}
	}

	return &Token{Kind: IdentifierTokenKind, Literal: buffer,
		Location: &utils.CodeBlockLocation{StartLocation: startLocation,
			EndLocation: l.currentLocation.Copy()}}
}

func (l *Lexer) nextWrappedIdentifierToken() *Token {
	startLocation := l.currentLocation.Copy()
	l.advance() // '`'
	for l.currentCodePoint != '`' {
		if l.currentCodePoint == '\n' || l.currentCodePoint == -1 {
			location := &utils.CodeBlockLocation{StartLocation: startLocation,
				EndLocation: l.currentLocation.Copy()}

			l.problemHandler.AddCodeProblem(
				utils.NewError(
					location,
					utils.NotClosedWrappedIdentifierErr, []interface{}{}))

			return &Token{Kind: IdentifierTokenKind,
				Literal:  string(l.source[startLocation.Index:l.currentLocation.Index]),
				Location: location}
		}
		l.advance()
	}

	l.advance() // '`' again

	buffer := string(l.source[startLocation.Index:l.currentLocation.Index])
	return &Token{Kind: IdentifierTokenKind, Literal: buffer,
		Location: &utils.CodeBlockLocation{StartLocation: startLocation,
			EndLocation: l.currentLocation.Copy()}}
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16 // larger than any legal digit val
}

func (l *Lexer) scanEscape(quote rune) bool {
	var n int
	var base, max uint32

	switch l.currentCodePoint {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		l.advance()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		l.advance()
		n, base, max = 2, 16, 255
	case 'u':
		l.advance()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		l.advance()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		err := utils.UnknownEscapeSequenceErr
		if l.currentCodePoint == -1 {
			err = utils.EscapeSequenceNotTerminatedErr
		}

		l.problemHandler.AddCodeProblem(
			utils.NewError(
				utils.NewOneCodePointBlockLocation(l.currentLocation),
				err, []interface{}{}))

		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(l.currentCodePoint))
		if d >= base {
			if l.currentCodePoint == -1 {
				l.problemHandler.AddCodeProblem(
					utils.NewError(
						utils.NewOneCodePointBlockLocation(l.currentLocation),
						utils.EscapeSequenceNotTerminatedErr, []interface{}{}))
			} else {
				l.problemHandler.AddCodeProblem(
					utils.NewError(
						utils.NewOneCodePointBlockLocation(l.currentLocation),
						utils.IllegalCharacterInEscapeSequenceErr,
						[]interface{}{l.currentCodePoint}))
			}
			return false
		}

		x = x*base + d
		l.advance()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		l.problemHandler.AddCodeProblem(
			utils.NewError(
				utils.NewOneCodePointBlockLocation(l.currentLocation),
				utils.EscapeSequenceIsInvalidUTF8CodePointErr,
				[]interface{}{}))
		return false
	}

	return true
}

func (l *Lexer) nextStringToken() *Token {
	startLocation := l.currentLocation.Copy()

	l.advance() // '"'

	for {
		if l.currentCodePoint == '\n' || l.currentCodePoint == -1 {
			location := &utils.CodeBlockLocation{StartLocation: startLocation,
				EndLocation: l.currentLocation.Copy()}

			l.problemHandler.AddCodeProblem(
				utils.NewError(
					location, utils.NotClosedStringErr, []interface{}{}))

			return &Token{Kind: StringTokenKind,
				Literal:  string(l.source[startLocation.Index+1 : l.currentLocation.Index-1]),
				Location: location}
		}

		if l.currentCodePoint == '"' {
			break
		}

		if l.currentCodePoint == '\\' {
			l.scanEscape('"')
		}

		l.advance()
	}

	l.advance() // '"' again

	buffer := string(l.source[startLocation.Index+1 : l.currentLocation.Index-1])
	return &Token{Kind: StringTokenKind, Literal: buffer,
		Location: &utils.CodeBlockLocation{StartLocation: startLocation,
			EndLocation: l.currentLocation.Copy()}}
}

func (l *Lexer) NextToken() *Token {
	l.skipWhitespaces()

	var result *Token

	switch l.currentCodePoint {
	case -1:
		return l.characterToken(EOFTokenKind, "\\0")
	case '`':
		return l.nextWrappedIdentifierToken()
	case '"':
		return l.nextStringToken()

	case '+':
		if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(PlusEqOpTokenKind, "+=")
		} else if l.peekByte() == '+' {
			l.advance()
			result = l.doubleCharacterToken(PlusPlusOpTokenKind, "++")
		} else {
			result = l.characterToken(PlusOpTokenKind, "+")
		}
	case '-':
		if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(MinusEqOpTokenKind, "-=")
		} else if l.peekByte() == '-' {
			l.advance()
			result = l.doubleCharacterToken(MinusMinusOpTokenKind, "--")
		} else {
			result = l.characterToken(MinusOpTokenKind, "-")
		}
	case '*':
		if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(MulEqOpTokenKind, "*=")
		} else {
			result = l.characterToken(MulOpTokenKind, "*")
		}
	case '/':
		if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(DivEqOpTokenKind, "/=")
		} else if l.peekByte() == '/' {
			result = l.nextCommentToken()
		} else if l.peekByte() == '*' {
			result = l.nextMultiLineCommentToken()
		} else {
			result = l.characterToken(DivOpTokenKind, "/")
		}

	case '^':
		if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(XOREqOpTokenKind, "^=")
		} else {
			result = l.characterToken(XOROpTokenKind, "^")
		}

	case '~':
		result = l.characterToken(NOTOpTokenKind, "~")

	case '(':
		result = l.characterToken(OpenParentTokenKind, "(")
	case ')':
		result = l.characterToken(CloseParentTokenKind, ")")
	case '[':
		result = l.characterToken(OpenBracketTokenKind, "[")
	case ']':
		result = l.characterToken(CloseBracketTokenKind, "]")
	case '{':
		result = l.characterToken(OpenBraceTokenKind, "{")
	case '}':
		result = l.characterToken(CloseBraceTokenKind, "}")

	case ',':
		result = l.characterToken(CommaTokenKind, ",")
	case ';':
		result = l.characterToken(SemiColTokenKind, ";")

	case '>':
		if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(GTEOpTokenKind, ">=")
		} else if l.peekByte() == '>' {
			l.advance()
			result = l.doubleCharacterToken(RShiftOpTokenKind, ">>")
		} else {
			result = l.characterToken(GTOpTokenKind, ">")
		}

	case '<':
		if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(LTEOpTokenKind, "<=")
		} else if l.peekByte() == '<' {
			l.advance()
			result = l.doubleCharacterToken(LShiftOpTokenKind, "<<")
		} else {
			result = l.characterToken(LTOpTokenKind, "<")
		}

	case '=':
		if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(EQOpTokenKind, "==")
		} else {
			result = l.characterToken(AssignOpTokenKind, "=")
		}

	case '!':
		if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(NEQOpTokenKind, "!=")
		} else {
			result = l.characterToken(BangOpTokenKind, "!")
		}

	case '&':
		if l.peekByte() == '&' {
			l.advance()
			result = l.doubleCharacterToken(ANDANDOpTokenKind, "&&")
		} else {
			result = l.characterToken(ANDOpTokenKind, "&")
		}

	case '|':
		if l.peekByte() == '|' {
			l.advance()
			result = l.doubleCharacterToken(OROROpTokenKind, "||")
		} else if l.peekByte() == '=' {
			l.advance()
			result = l.doubleCharacterToken(OREqOpTokenKind, "|=")
		} else {
			result = l.characterToken(OROpTokenKind, "|")
		}

	default:
		if isIdentifierStart(l.currentCodePoint) {
			return l.nextNameToken()
		} else if isDecimal(l.currentCodePoint) || l.currentCodePoint == '.' && isDecimal(rune(l.peekByte())) {
			result = l.nextNumberToken()
		} else if l.currentCodePoint == '.' {
			result = l.characterToken(DotTokenKind, ".")
		} else {
			l.problemHandler.AddCodeProblem(
				utils.NewError(
					utils.NewOneCodePointBlockLocation(l.currentLocation),
					utils.UnexpectedCharacterErr, []interface{}{l.currentCodePoint}))
			result = l.characterToken(InvalidTokenKind, string(l.currentCodePoint))
		}
	}

	l.advance()

	return result
}
