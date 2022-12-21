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

package utils

const (
	IllegalNullCharacterErr = iota
	IllegalUTF8EncodingErr
	UnexpectedCharacterErr
	NotClosedWrappedIdentifierErr
	NotClosedMultiLineCommentErr
	NotClosedStringErr
	UnknownEscapeSequenceErr
	EscapeSequenceNotTerminatedErr
	IllegalCharacterInEscapeSequenceErr
	EscapeSequenceIsInvalidUTF8CodePointErr
	InvalidRadixPointErr
	HasNoDigitsErr
	ExponentRequiresDecimalMantissaErr
	ExponentRequiresHexadecimalMantissaErr
	ExponentHasNoDigitsErr
	HexadecimalMantissaRequiresPExponentErr
	InvalidDigitErr
	UnderscoreMustSeparateSuccessiveDigitsErr
)

var error_messages = map[int]string{
	IllegalNullCharacterErr:                   "illegal null character",
	IllegalUTF8EncodingErr:                    "illegal UTF-8 encoding",
	UnexpectedCharacterErr:                    "unexpected character `%c`",
	NotClosedWrappedIdentifierErr:             "not closed wrapped identifier",
	NotClosedMultiLineCommentErr:              "not closed multiline comment",
	NotClosedStringErr:                        "not closed string literal",
	UnknownEscapeSequenceErr:                  "unknown escape sequence",
	EscapeSequenceNotTerminatedErr:            "escape sequence not terminated",
	IllegalCharacterInEscapeSequenceErr:       "illegal character %#U in escape sequence",
	EscapeSequenceIsInvalidUTF8CodePointErr:   "escape sequence is invalid Unicode code point",
	InvalidRadixPointErr:                      "invalid radix point in %s",
	HasNoDigitsErr:                            "%s has no digits",
	ExponentRequiresDecimalMantissaErr:        "%q exponent requires decimal mantissa",
	ExponentRequiresHexadecimalMantissaErr:    "%q exponent requires hexadecimal mantissa",
	ExponentHasNoDigitsErr:                    "exponent has no digits",
	HexadecimalMantissaRequiresPExponentErr:   "hexadecimal mantissa requires `p` exponent",
	InvalidDigitErr:                           "invalid digit %q in %s",
	UnderscoreMustSeparateSuccessiveDigitsErr: "`_` must separate successive digits",
}
