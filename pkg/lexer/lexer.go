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

// Package scanner implements a scanner for Tinysource text.
// It takes a []byte as source which can then be tokenized
// through repeated calls to the NextToken method.
package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/tinylang-org/tiny/pkg/utils"
)

// Lexer struct is state machine, main goal of which is to scan source code and convert
// it to tokens.
// To get them we can use NextToken(), which returns the token in current "location" (function
// always returns different values).
// and increases location in code (or cursor) until it finds next token.
// Last token in the source file is EOF.
//
//  pub fun main() { ...
//  ^
//  |
//  "cursor" after lexer is initialized
//
// then we call NextToken() function and get Token(Literal="pub", Kind=PubKeywordTokenKind)
// while the function is executed cursor is moved to the next token:
//
//  pub fun main() { ..
//      ^
//      |
//     "cursor" after NextToken() is called
//
// then when we call NextToken() again, the return is Token(Literal="fun",
// Kind=FunKeywordTokenKind)
//
// when the cursor is at the end of the source code, NextToken() function will
// return Token(Literal="0", Kind=EOFTokenKind).
type Lexer struct {
	filepath     string
	source       []byte
	sourceLength int

	currentCodePoint     rune
	currentCodePointSize int
	currentLocation      *utils.CodePointLocation

	problemHandler *utils.CodeProblemHandler

	LineStartOffsets *[]int
	LineEndOffsets   *[]int
}

func NewLexer(filepath string, source []byte, problemHandler *utils.CodeProblemHandler) *Lexer {
	l := &Lexer{
		filepath:     filepath,
		source:       source,
		sourceLength: len(source),

		currentLocation: &utils.CodePointLocation{
			Filepath: filepath,
			Index:    0,
			Line:     1,
			Column:   0,
		},

		problemHandler: problemHandler,

		LineStartOffsets: &[]int{0},
		LineEndOffsets:   &[]int{},
	}

	l.decodeRune()

	// Check for the byte order mark in the beginning of the file
	if l.currentCodePoint == 0xFEFF {
		l.advance()
	}

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
				utils.NewLocalError(
					utils.NewOneCodePointBlockLocation(l.currentLocation),
					utils.IllegalNullCharacterErr, []interface{}{}))
		case r >= utf8.RuneSelf:
			r, offset = utf8.DecodeRune(l.source[l.currentLocation.Index:])
			if r == utf8.RuneError && offset == 1 {
				l.problemHandler.AddCodeProblem(
					utils.NewLocalError(utils.NewOneCodePointBlockLocation(l.currentLocation),
						utils.IllegalUTF8EncodingErr))
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
		*l.LineEndOffsets = append(*l.LineEndOffsets, l.currentLocation.Index-1)
		l.currentLocation.Line++
		l.currentLocation.Column = 0
		*l.LineStartOffsets = append(*l.LineStartOffsets, l.currentLocation.Index+1)
	} else {
		l.currentLocation.Column++
	}

	l.currentLocation.Index += l.currentCodePointSize

	l.decodeRune()
}

// Skip whitespaces block.
//   var  \t   a: i32 = 3;
//     ???      ???
//     ???       ???
//     ???       ending here
//     starting from here
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

// Binary search throw keyword list and try to find a keyword.
// If not found return -1.
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

// Return a byte that is located after the current one in source text.
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

// Scan digits sequence.
//   2938487198.34
//   ???         ???
//   ???          ???
//   ???         digits() is stopped
//   digits() is called
func (l *Lexer) digits(base int, invalid **utils.CodePointLocation) (digitSeparator int) {
	if base <= 10 {
		max := rune('0' + base)
		for isDecimal(l.currentCodePoint) || l.currentCodePoint == '_' {
			ds := 1
			if l.currentCodePoint == '_' {
				ds = 2
			} else if l.currentCodePoint >= max && *invalid == nil {
				*invalid = l.currentLocation.Copy() // record invalid rune location
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

	// location of invalid digit in literal, or nil
	var invalid *utils.CodePointLocation = nil

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
				utils.NewLocalError(
					utils.NewOneCodePointBlockLocation(l.currentLocation),
					utils.InvalidRadixPointErr, numberLiteralName(prefix)))
		}
		l.advance()
		digitSeparator |= l.digits(base, &invalid)
	}

	if digitSeparator&1 == 0 {
		l.problemHandler.AddCodeProblem(
			utils.NewLocalError(
				&utils.CodeBlockLocation{StartLocation: startLocation,
					EndLocation: l.currentLocation.Copy()},
				utils.HasNoDigitsErr, numberLiteralName(prefix)))
	}

	// exponent
	if e := lower(l.currentCodePoint); e == 'e' || e == 'p' {
		switch {
		case e == 'e' && prefix != 0 && prefix != '0':
			l.problemHandler.AddCodeProblem(
				utils.NewLocalError(
					utils.NewOneCodePointBlockLocation(l.currentLocation),
					utils.ExponentRequiresDecimalMantissaErr,
					l.currentCodePoint))
		case e == 'p' && prefix != 'x':
			l.problemHandler.AddCodeProblem(
				utils.NewLocalError(
					utils.NewOneCodePointBlockLocation(l.currentLocation),
					utils.ExponentRequiresHexadecimalMantissaErr,
					l.currentCodePoint))
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
				utils.NewLocalError(
					&utils.CodeBlockLocation{StartLocation: startLocation,
						EndLocation: l.currentLocation.Copy()},
					utils.ExponentHasNoDigitsErr))
		}
	} else if prefix == 'x' && tokenKind == FloatTokenKind {
		l.problemHandler.AddCodeProblem(
			utils.NewLocalError(
				&utils.CodeBlockLocation{StartLocation: startLocation,
					EndLocation: l.currentLocation.Copy()},
				utils.HexadecimalMantissaRequiresPExponentErr))
	}

	// suffix 'i'
	if l.currentCodePoint == 'i' {
		tokenKind = ImaginaryTokenKind
		l.advance()
	}

	buffer := string(l.source[startLocation.Index:l.currentLocation.Index])

	if tokenKind == IntTokenKind && invalid != nil {
		l.problemHandler.AddCodeProblem(
			utils.NewLocalError(
				utils.NewOneCodePointBlockLocation(invalid),
				utils.InvalidDigitErr,
				buffer[invalid.Index-startLocation.Index],
				numberLiteralName(prefix)))
	}

	if digitSeparator&2 != 0 {
		if i := invalidSeparator(buffer); i >= 0 {
			l.problemHandler.AddCodeProblem(
				utils.NewLocalError(
					&utils.CodeBlockLocation{StartLocation: startLocation,
						EndLocation: l.currentLocation.Copy()},
					utils.UnderscoreMustSeparateSuccessiveDigitsErr))
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
				utils.NewLocalError(
					utils.NewOneCodePointBlockLocation(startLocation),
					utils.NotClosedMultiLineCommentErr))
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

// From https://github.com/golang/go/blob/db36eca33c389871b132ffb1a84fd534a349e8d8/src/go/scanner/scanner.go#L663
func stripCR(b []byte, comment bool) []byte {
	c := make([]byte, len(b))
	i := 0
	for j, ch := range b {
		// In a /*-style comment, don't strip \r from *\r/ (incl.
		// sequences of \r from *\r\r...\r/) since the resulting
		// */ would terminate the comment too early unless the \r
		// is immediately following the opening /* in which case
		// it's ok because /*/ is not closed yet
		// (issue #11151 from https://github.com/golang/go).
		if ch != '\r' || comment && i > len("/*") && c[i-1] == '*' && j+1 < len(b) && b[j+1] == '/' {
			c[i] = ch
			i++
		}
	}
	return c[:i]
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

	hasCR := false

	for l.currentCodePoint != '`' {
		if l.currentCodePoint == '\r' {
			hasCR = true
		}

		if l.currentCodePoint == '\n' || l.currentCodePoint == -1 {
			var endLocation *utils.CodePointLocation

			if l.currentCodePoint == '\n' && hasCR {
				endLocation = l.currentLocation.PreviousByteLocation()
			} else {
				endLocation = l.currentLocation.Copy()
			}

			location := &utils.CodeBlockLocation{StartLocation: startLocation,
				EndLocation: endLocation}

			l.problemHandler.AddCodeProblem(
				utils.NewLocalError(
					location,
					utils.NotClosedWrappedIdentifierErr))

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
		if l.currentCodePoint < 0 {
			err = utils.EscapeSequenceNotTerminatedErr
		}

		l.problemHandler.AddCodeProblem(
			utils.NewLocalError(
				utils.NewOneCodePointBlockLocation(l.currentLocation),
				err))

		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(l.currentCodePoint))
		if d >= base {
			if l.currentCodePoint < 0 {
				l.problemHandler.AddCodeProblem(
					utils.NewLocalError(
						utils.NewOneCodePointBlockLocation(l.currentLocation),
						utils.EscapeSequenceNotTerminatedErr))
			} else {
				l.problemHandler.AddCodeProblem(
					utils.NewLocalError(
						utils.NewOneCodePointBlockLocation(l.currentLocation),
						utils.IllegalCharacterInEscapeSequenceErr,
						l.currentCodePoint))
			}
			return false
		}

		x = x*base + d
		l.advance()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		l.problemHandler.AddCodeProblem(
			utils.NewLocalError(
				utils.NewOneCodePointBlockLocation(l.currentLocation),
				utils.EscapeSequenceIsInvalidUTF8CodePointErr))
		return false
	}

	return true
}

func (l *Lexer) nextStringToken() *Token {
	startLocation := l.currentLocation.Copy()
	hasCR := false

	l.advance() // '"'

	for {
		if l.currentCodePoint == '\r' {
			hasCR = true
		}

		if l.currentCodePoint == '\n' || l.currentCodePoint == -1 {
			var endLocation *utils.CodePointLocation

			if l.currentCodePoint == '\n' && hasCR {
				endLocation = l.currentLocation.PreviousByteLocation()
			} else {
				endLocation = l.currentLocation.Copy()
			}

			location := &utils.CodeBlockLocation{StartLocation: startLocation,
				EndLocation: endLocation}

			l.problemHandler.AddCodeProblem(
				utils.NewLocalError(
					location, utils.NotClosedStringErr))

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

// NextToken scans the next token and return it. If there scanning error it will
// be added to l.diagnostics. If there's unexpected character InvalidTokenKind will
// be returned. If another type of error will be occured in the process
// of scanning, then lexer will try to return a valid token.
func (l *Lexer) NextToken() *Token {
	l.skipWhitespaces()

	var result *Token

	switch l.currentCodePoint {
	case -1:
		{
			*l.LineEndOffsets = append(*l.LineEndOffsets, l.currentLocation.Index-1)
			return l.characterToken(EOFTokenKind, "\\0")
		}
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
			return l.nextNumberToken()
		} else if l.currentCodePoint == '.' {
			result = l.characterToken(DotTokenKind, ".")
		} else {
			l.problemHandler.AddCodeProblem(
				utils.NewLocalError(
					utils.NewOneCodePointBlockLocation(l.currentLocation),
					utils.UnexpectedCharacterErr, l.currentCodePoint))
			result = l.characterToken(InvalidTokenKind, string(l.currentCodePoint))
		}
	}

	l.advance()

	return result
}
