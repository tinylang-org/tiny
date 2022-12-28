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
	"strings"

	"github.com/tinylang-org/tiny/pkg/ast"
	"github.com/tinylang-org/tiny/pkg/lexer"
	"github.com/tinylang-org/tiny/pkg/utils"
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
	filepath string

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
	p := &Parser{filepath: filepath,
		lexer: lexer.NewLexer(filepath, source, problem_handler)}
	p.problem_handler = problem_handler

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

	p.registerInfixFunction(lexer.DotTokenKind, p.parseMemberExpression)
	p.registerInfixFunction(lexer.OpenParentTokenKind, p.parseFunctionCallExpression)
	p.registerInfixFunction(lexer.OpenBracketTokenKind, p.parseIndexExpression)

	p.advance()
	p.advance()
	return p
}

func (p *Parser) ParseProgramUnit() *ast.ProgramUnit {
	namespace := p.parseNamespaceDecl()
	if namespace == nil {
		return nil
	}

	imports := p.parseImports()
	TLStatements := p.parseTopLevelStatementList()

	return &ast.ProgramUnit{
		Filepath:     p.filepath,
		Namespace:    namespace,
		Imports:      imports,
		TLStatements: TLStatements,
	}
}

func (p *Parser) parseNamespaceDecl() *ast.NamespaceDecl {
	if !p.expectCurrent(lexer.NamespaceKeywordTokenKind) {
		return nil
	}

	if !p.expectPeek(lexer.StringTokenKind) {
		return nil
	}

	location := &utils.CodeBlockLocation{
		StartLocation: p.currentToken.Location.StartLocation.Copy(),
		EndLocation:   p.peekToken.Location.EndLocation.Copy(),
	}

	name := p.currentToken.Literal

	if !p.expectPeek(lexer.SemiColTokenKind) {
		return nil
	}

	p.advance() // `;`

	return &ast.NamespaceDecl{NamespaceLocation: location, Name: name}
}

func (p *Parser) parseImports() []*ast.Import {
	imports := []*ast.Import{}

	for p.currentToken.Kind == lexer.ImportKeywordTokenKind {
		import_decl := p.parseImport()

		if import_decl != nil {
			imports = append(imports, import_decl)
		}

		p.advance() // ';'
	}

	return imports
}

func (p *Parser) parseImport() *ast.Import {
	if !p.expectPeek(lexer.StringTokenKind) {
		return nil
	}

	location := &utils.CodeBlockLocation{
		StartLocation: p.currentToken.Location.StartLocation.Copy().Copy(),
		EndLocation:   p.peekToken.Location.EndLocation.Copy(),
	}

	path := p.currentToken.Literal

	if !p.expectPeek(lexer.SemiColTokenKind) {
		return nil
	}

	return &ast.Import{ImportLocation: location, Path: path}
}

// top_level_statement = function_declaration |
//
//	struct_declaration .
func (p *Parser) parseTopLevelStatementList() []ast.TopLevelStatement {
	var list []ast.TopLevelStatement

	for p.currentToken.Kind != lexer.EOFTokenKind {
		stmt := p.parseTopLevelStatement()
		if stmt != nil {
			list = append(list, stmt)
		}

		p.advance() // '}'
	}

	return list
}

func (p *Parser) parseTopLevelStatement() ast.TopLevelStatement {
	switch p.currentToken.Kind {
	case lexer.PubKeywordTokenKind:
		switch p.peekToken.Kind {
		case lexer.FunKeywordTokenKind:
			return p.parseFunctionDeclaration(true)
		case lexer.StructKeywordTokenKind:
			return p.parseStructureDeclaration(true)
		default:
			return nil
		}

	case lexer.FunKeywordTokenKind:
		return p.parseFunctionDeclaration(false)
	case lexer.StructKeywordTokenKind:
		return p.parseStructureDeclaration(false)
	default:
		return nil
	}
}

func (p *Parser) parseFunctionDeclaration(public bool) ast.TopLevelStatement {
	startLocation := p.currentToken.Location.StartLocation

	if public {
		p.advance() // 'pub'
	}

	if !p.expectPeek(lexer.IdentifierTokenKind) {
		return nil
	}

	functionName := p.currentToken.Literal

	// TODO: generics

	if !p.expectPeek(lexer.OpenParentTokenKind) {
		return nil
	}

	var arguments []*ast.FunctionArgument

	p.advance()
	if p.currentTokenIs(lexer.CloseParentTokenKind) {
		arguments = []*ast.FunctionArgument{}
	} else {
		arguments = p.parseFunctionArguments()

		if !p.expectCurrent(lexer.CloseParentTokenKind) {
			return nil
		}
	}

	if !p.expectPeek(lexer.OpenBraceTokenKind) {
		return nil
	}

	statementsBlockStartLocation := p.currentToken.Location.StartLocation

	p.advance()

	var statements []ast.Statement
	if p.currentTokenIs(lexer.CloseBraceTokenKind) {
		statements = []ast.Statement{}
	} else {
		statements = p.parseStatementList()
	}

	if !p.expectCurrent(lexer.CloseBraceTokenKind) {
		return nil
	}

	endLocation := p.currentToken.Location.EndLocation

	return &ast.FunctionDeclaration{Public: public,
		Name:      functionName,
		Arguments: arguments,
		StatementsBlock: &ast.StatementsBlock{
			Statements:    statements,
			StartLocation: statementsBlockStartLocation,
		},
		BlockLocation: &utils.CodeBlockLocation{
			StartLocation: startLocation,
			EndLocation:   endLocation,
		},
	}
}

func (p *Parser) parseFunctionArguments() []*ast.FunctionArgument {
	var arguments []*ast.FunctionArgument

	for p.currentToken.Kind != lexer.CloseParentTokenKind {
		argument := p.parseFunctionArgument()
		arguments = append(arguments, argument)

		if p.currentToken.Kind == lexer.CommaTokenKind {
			p.advance() // skip comma
		}
	}

	return arguments
}

func (p *Parser) parseFunctionArgument() *ast.FunctionArgument {
	if !p.expectCurrent(lexer.IdentifierTokenKind) {
		return nil
	}

	startLocation := p.currentToken.Location.StartLocation.Copy()
	name := p.currentToken.Literal
	typeDef := p.parseType()
	p.advance()
	endLocation := p.currentToken.Location.EndLocation.Copy()

	return &ast.FunctionArgument{
		Name: name,
		Type: typeDef,
		BlockLocation: &utils.CodeBlockLocation{
			StartLocation: startLocation,
			EndLocation:   endLocation,
		},
	}
}

func (p *Parser) parseStructureDeclaration(public bool) ast.TopLevelStatement {
	return nil
}

func (p *Parser) parseType() ast.Type {
	switch p.currentToken.Kind {
	case lexer.Int8KeywordTokenKind:
		return &ast.PrimaryType{Token: p.currentToken}
	case lexer.Int16KeywordTokenKind:
		return &ast.PrimaryType{Token: p.currentToken}
	case lexer.Int32KeywordTokenKind:
		return &ast.PrimaryType{Token: p.currentToken}
	case lexer.Int64KeywordTokenKind:
		return &ast.PrimaryType{Token: p.currentToken}
	case lexer.Uint8KeywordTokenKind:
		return &ast.PrimaryType{Token: p.currentToken}
	case lexer.Uint16KeywordTokenKind:
		return &ast.PrimaryType{Token: p.currentToken}
	case lexer.Uint32KeywordTokenKind:
		return &ast.PrimaryType{Token: p.currentToken}
	case lexer.Uint64KeywordTokenKind:
		return &ast.PrimaryType{Token: p.currentToken}
	case lexer.MulOpTokenKind:
		return p.parsePointerType()
	case lexer.OpenBracketTokenKind:
		return p.parseArrayType()
	case lexer.IdentifierTokenKind:
		return p.parseCustomType()
	default:
		p.addUnexpectedCurrentTokenError()
	}

	return nil
}

func (p *Parser) parsePointerType() ast.Type {
	startLocation := p.currentToken.Location.StartLocation.Copy()
	p.advance()

	pointerType := p.parseType()
	return &ast.PointerType{Type: pointerType, StartLocation: startLocation}
}

func (p *Parser) parseArrayType() ast.Type {
	startLocation := p.currentToken.Location.StartLocation.Copy()

	if !p.expectPeek(lexer.CloseBracketTokenKind) {
		return nil
	}

	p.advance()

	arrayType := p.parseType()

	return &ast.ArrayType{Type: arrayType, StartLocation: startLocation}
}

func (p *Parser) parseCustomType() ast.Type {
	var name strings.Builder

	startLocation := p.currentToken.Location.StartLocation.Copy()

	for p.currentToken.Kind == lexer.IdentifierTokenKind {
		name.WriteString(p.currentToken.Literal)
		name.WriteByte('.')

		if !p.expectPeekNoErr(lexer.DotTokenKind) {
			buffer := name.String()
			buffer = buffer[0 : len(buffer)-1]
			return &ast.CustomType{Name: buffer, TypeLocation: &utils.CodeBlockLocation{
				StartLocation: startLocation,
				EndLocation:   p.currentToken.Location.StartLocation.Copy(),
			}}
		}

		p.advance() // '.'
	}

	buffer := name.String()
	buffer = buffer[0 : len(buffer)-1]
	return &ast.CustomType{Name: buffer, TypeLocation: &utils.CodeBlockLocation{
		StartLocation: startLocation,
		EndLocation:   p.currentToken.Location.StartLocation.Copy(),
	}}
}

func (p *Parser) parseStatementList() []ast.Statement {
	statements := []ast.Statement{}

	p.advance() // openbrace
	for p.currentToken.Kind != lexer.EOFTokenKind {
		statement := p.parseStatement()

		if statement != nil {
			statements = append(statements, statement)
		}

		p.advance() // ';'
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
	statement := &ast.ReturnStatement{TokenLocation: p.currentToken.Location.Copy()}

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
		StartLocation: p.currentToken.Location.StartLocation.Copy(),
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
		TokenLocation: p.currentToken.Location.Copy(),
		Value:         p.currentToken.Literal == "true"}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		TokenLocation: p.currentToken.Location.Copy(),
		Value:         p.currentToken.Literal}
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	return nil
}

func (p *Parser) parseFunctionCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Function: function}
	expression.Arguments, expression.EndLocation = p.parseExpressionList(lexer.CloseParentTokenKind)
	return expression
}

func (p *Parser) parseExpressionList(endTokenKind int) ([]ast.Expression, *utils.CodePointLocation) {
	var list []ast.Expression

	if p.peekTokenIs(endTokenKind) {
		p.advance()
		return list, p.currentToken.Location.EndLocation
	}

	p.advance()
	list = append(list, p.parseExpression(Lowest))

	for p.peekTokenIs(lexer.CommaTokenKind) {
		p.advance()
		p.advance()
		list = append(list, p.parseExpression(Lowest))
	}

	if !p.expectPeek(endTokenKind) {
		return nil, nil
	}

	return list, p.currentToken.Location.EndLocation
}

func (p *Parser) addUnexpectedCurrentTokenError() {
	p.problem_handler.AddCodeProblem(utils.NewLocalError(p.currentToken.Location.Copy(),
		utils.UnexpectedToken2Err, []interface{}{
			lexer.DumpTokenKind(p.currentToken.Kind)}))
}

func (p *Parser) addUnexpectedPeekTokenError() {
	p.problem_handler.AddCodeProblem(utils.NewLocalError(p.peekToken.Location.Copy(),
		utils.UnexpectedToken2Err, []interface{}{
			lexer.DumpTokenKind(p.peekToken.Kind)}))
}

func (p *Parser) expectCurrent(tokenKind int) bool {
	if p.currentTokenIs(tokenKind) {
		return true
	} else {
		p.problem_handler.AddCodeProblem(utils.NewLocalError(p.currentToken.Location.Copy(),
			utils.UnexpectedTokenErr, []interface{}{lexer.DumpTokenKind(tokenKind),
				lexer.DumpTokenKind(p.currentToken.Kind)}))
		return false
	}
}

func (p *Parser) expectPeekNoErr(tokenKind int) bool {
	if p.peekTokenIs(tokenKind) {
		p.advance()
		return true
	} else {
		return false
	}
}

func (p *Parser) expectPeek(tokenKind int) bool {
	if p.peekTokenIs(tokenKind) {
		p.advance()
		return true
	} else {
		p.problem_handler.AddCodeProblem(utils.NewLocalError(p.peekToken.Location.Copy(),
			utils.UnexpectedTokenErr, []interface{}{lexer.DumpTokenKind(tokenKind),
				lexer.DumpTokenKind(p.peekToken.Kind)}))
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
