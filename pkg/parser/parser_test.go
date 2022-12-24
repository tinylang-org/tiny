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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vertexgmd/tinylang/pkg/ast"
	"github.com/vertexgmd/tinylang/pkg/lexer"
	"github.com/vertexgmd/tinylang/pkg/utils"
)

func TestBooleanLiteral(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("true"), p)
	statement := parser.parseStatement()
	assert.Equal(t, true, statement.(*ast.BooleanLiteral).Value)
}

func TestStringLiteral(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("\"hello\""), p)
	statement := parser.parseStatement()
	assert.Equal(t, "hello", statement.(*ast.StringLiteral).Value)
}

func TestReturnStatement(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("return \"hello\";"), p)
	statement := parser.parseStatement()
	assert.Equal(t, true,
		statement.(*ast.ReturnStatement).HasReturnValue)
	assert.Equal(t, "hello",
		statement.(*ast.ReturnStatement).ReturnValue.(*ast.StringLiteral).Value)
}

func TestPrimaryType(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("int32"), p)
	tp := parser.parseType()
	assert.Equal(t, lexer.Int32KeywordTokenKind,
		tp.(*ast.PrimaryType).Token.Kind)
	p.PrintProblems()
}

func TestPointerType(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("*int32"), p)
	tp := parser.parseType()
	assert.Equal(t, lexer.Int32KeywordTokenKind,
		tp.(*ast.PointerType).Type.(*ast.PrimaryType).Token.Kind)
	p.PrintProblems()
}

func TestPointerType2(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("**int8"), p)
	tp := parser.parseType()
	assert.Equal(t, lexer.Int8KeywordTokenKind,
		tp.(*ast.PointerType).Type.(*ast.PointerType).Type.(*ast.PrimaryType).Token.Kind)
	p.PrintProblems()
}

func TestArrayType(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("[]int8"), p)
	tp := parser.parseType()
	assert.Equal(t, lexer.Int8KeywordTokenKind,
		tp.(*ast.ArrayType).Type.(*ast.PrimaryType).Token.Kind)
	p.PrintProblems()
}

func TestArrayType2(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("*[]*int8"), p)
	tp := parser.parseType()
	assert.Equal(t, lexer.Int8KeywordTokenKind,
		tp.(*ast.PointerType).Type.(*ast.ArrayType).Type.(*ast.PointerType).Type.(*ast.PrimaryType).Token.Kind)
	p.PrintProblems()
}

func TestCustomType(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("custom"), p)
	tp := parser.parseType()
	assert.Equal(t, "custom",
		tp.(*ast.CustomType).Name)
	p.PrintProblems()
}

func TestCustomType2(t *testing.T) {
	p := utils.NewCodeProblemHandler()
	parser := NewParser("", []byte("mylib.custom"), p)
	tp := parser.parseType()
	assert.Equal(t, "mylib.custom",
		tp.(*ast.CustomType).Name)
	p.PrintProblems()
}
