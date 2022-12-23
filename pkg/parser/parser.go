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
	"github.com/vertexgmd/tinylang/pkg/ast"
	"github.com/vertexgmd/tinylang/pkg/lexer"
	"github.com/vertexgmd/tinylang/pkg/utils"
)

const (
	_ int = iota
	Lowest
	Equals
	LessOrGreater
	Sum
	Product
	Prefix
	FunctionCall
	Index
)

var precedences = map[int]int{
	lexer.EQOpTokenKind:        Equals,
	lexer.NEQOpTokenKind:       Equals,
	lexer.LTOpTokenKind:        LessOrGreater,
	lexer.GTOpTokenKind:        LessOrGreater,
	lexer.LTEOpTokenKind:       LessOrGreater,
	lexer.GTEOpTokenKind:       LessOrGreater,
	lexer.PlusOpTokenKind:      Sum,
	lexer.MinusOpTokenKind:     Sum,
	lexer.MulOpTokenKind:       Product,
	lexer.DivOpTokenKind:       Product,
	lexer.OpenParentTokenKind:  FunctionCall,
	lexer.OpenBracketTokenKind: Index,
}

type Parser struct {
	problem_handler *utils.CodeProblemHandler
	lexer           *lexer.Lexer

	prefixParseFunctions map[int]prefixParseFunction
	infixParseFunctions  map[int]infixParseFunction

	currentToken *lexer.Token
	peekToken    *lexer.Token
}

type (
	prefixParseFunction func() ast.Expression
	infixParseFunction  func(ast.Expression) ast.Expression
)

func NewParser(filepath string, source []byte,
	problem_handler *utils.CodeProblemHandler) *Parser {
	p := &Parser{lexer: lexer.NewLexer(filepath, source, problem_handler)}

	p.prefixParseFunctions = make(map[int]prefixParseFunction)
	p.infixParseFunctions = make(map[int]infixParseFunction)

	p.registerPrefixFunction(lexer.BooleanTokenKind, p.parseBooleanLiteral)
	p.registerPrefixFunction(lexer.StringTokenKind, p.parseStringLiteral)

	p.registerInfixFunction(lexer.PlusOpTokenKind, p.parseInfixExpression)
	p.registerInfixFunction(lexer.MinusOpTokenKind, p.parseInfixExpression)

	p.registerInfixFunction(lexer.MulOpTokenKind, p.parseInfixExpression)
	p.registerInfixFunction(lexer.DivOpTokenKind, p.parseInfixExpression)

	p.registerInfixFunction(lexer.EQOpTokenKind, p.parseInfixExpression)
	p.registerInfixFunction(lexer.NEQOpTokenKind, p.parseInfixExpression)

	p.registerInfixFunction(lexer.LTOpTokenKind, p.parseInfixExpression)
	p.registerInfixFunction(lexer.GTOpTokenKind, p.parseInfixExpression)

	p.registerInfixFunction(lexer.LTEOpTokenKind, p.parseInfixExpression)
	p.registerInfixFunction(lexer.GTEOpTokenKind, p.parseInfixExpression)

	// p.registerInfixFunction(lexer.OpenParentTokenKind, p.parseFunctionCallExpression)
	p.registerInfixFunction(lexer.OpenBracketTokenKind, p.parseIndexExpression)

	p.advance()
	p.advance()
	return p
}

func (p *Parser) parseStatementList() []ast.Statement {
	statements := []ast.Statement{}

	for p.currentToken.Kind != lexer.EOFTokenKind {
		statement := p.parseStatement()

		if statement != nil {
			statements = append(statements, statement)
		}

		p.advance()
	}

	return statements
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Kind {
	case lexer.ReturnKeywordTokenKind:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	statement := &ast.ReturnStatement{TokenLocation: p.currentToken.Location}

	p.advance()
	if p.currentTokenIs(lexer.SemiColTokenKind) {
		statement.HasReturnValue = false
		p.advance()
	}

	statement.ReturnValue = p.parseExpression(Lowest)
	statement.HasReturnValue = true

	if p.expectPeek(lexer.SemiColTokenKind) {
		p.advance()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expression := p.parseExpression(Lowest)

	if p.peekTokenIs(lexer.SemiColTokenKind) {
		p.advance()
	}

	return expression
}

// func (p *Parser) parseVarStatement() ast.Statement {}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFunctions[p.currentToken.Kind]
	if prefix == nil {
		return nil
	}

	leftExpression := prefix()

	if !p.peekTokenIs(lexer.SemiColTokenKind) && precedence < p.peekPrecedence() {
		infix := p.infixParseFunctions[p.peekToken.Kind]
		if infix == nil {
			return leftExpression
		}

		p.advance()
		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Operator:      p.currentToken.Literal,
		StartLocation: p.currentToken.Location.StartLocation,
	}

	p.advance()
	expression.Expression = p.parseExpression(Prefix)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Left:     left,
		Operator: p.currentToken.Literal,
	}

	precedence := p.currentPrecedence()
	p.advance()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{
		Left: left,
	}

	p.advance()
	expression.Index = p.parseExpression(Lowest)

	if !p.expectPeek(lexer.CloseBracketTokenKind) {
		return nil
	}

	return expression
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		TokenLocation: p.currentToken.Location,
		Value:         p.currentToken.Literal == "true"}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		TokenLocation: p.currentToken.Location,
		Value:         p.currentToken.Literal}
}

func (p *Parser) expectPeek(tokenKind int) bool {
	if p.peekTokenIs(tokenKind) {
		p.advance()
		return true
	} else {
		return false
	}
}

func (p *Parser) currentTokenIs(tokenKind int) bool {
	return p.currentToken.Kind == tokenKind
}

func (p *Parser) peekTokenIs(tokenKind int) bool {
	return p.peekToken.Kind == tokenKind
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Kind]; ok {
		return p
	}

	return Lowest
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Kind]; ok {
		return p
	}

	return Lowest
}

func (p *Parser) registerPrefixFunction(tokenKind int, function prefixParseFunction) {
	p.prefixParseFunctions[tokenKind] = function
}

func (p *Parser) registerInfixFunction(tokenKind int, function infixParseFunction) {
	p.infixParseFunctions[tokenKind] = function
}

func (p *Parser) advance() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}
