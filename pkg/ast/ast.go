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

package ast

import (
	"github.com/vertexgmd/tinylang/pkg/lexer"
	"github.com/vertexgmd/tinylang/pkg/utils"
)

type AST interface {
	Location() *utils.CodeBlockLocation
}

type Statement interface {
	AST
	statementNode()
}

type Type interface {
	AST
	typeNode()
}

type TopLevelStatement interface {
	AST
	topLevelStatement()
}

type Expression interface {
	AST
	Statement
	expressionNode()
}

type Imports struct {
	Imports  []string
	location *utils.CodeBlockLocation
}

func (i *Imports) Location() *utils.CodeBlockLocation { return i.location }

type NamespaceDecl struct {
	Name     string
	location *utils.CodeBlockLocation
}

func (n *NamespaceDecl) Location() *utils.CodeBlockLocation { return n.location }

type ProgramUnit struct {
	Namespace     *NamespaceDecl
	Imports       *Imports
	Tl_statements []TopLevelStatement
}

func (p *ProgramUnit) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: p.Namespace.Location().StartLocation,
		EndLocation: p.Tl_statements[len(p.Tl_statements)-1].Location().EndLocation}
}

type StatementsBlock struct {
	// location of '{'
	startLocation *utils.CodePointLocation
	Statements    []Statement
}

func (s *StatementsBlock) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: s.startLocation,
		EndLocation: s.Statements[len(s.Statements)-1].Location().EndLocation}
}

type VarStatement struct {
	// location of 'var'
	startLocation *utils.CodePointLocation
	name          *Name
	var_type      Type
	value         Expression
}

func (v *VarStatement) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: v.startLocation,
		EndLocation: v.value.Location().EndLocation}
}
func (v *VarStatement) statementNode()     {}
func (v *VarStatement) topLevelStatement() {}

type ReturnStatement struct {
	// location of 'return'
	TokenLocation *utils.CodeBlockLocation

	HasReturnValue bool
	ReturnValue    Expression
}

func (r *ReturnStatement) Location() *utils.CodeBlockLocation {
	if !r.HasReturnValue {
		return r.TokenLocation
	}

	return &utils.CodeBlockLocation{StartLocation: r.TokenLocation.StartLocation,
		EndLocation: r.ReturnValue.Location().EndLocation}
}

func (r *ReturnStatement) statementNode() {}

type PrefixExpression struct {
	// location of operator
	StartLocation *utils.CodePointLocation

	Operator   string
	Expression Expression
}

func (p *PrefixExpression) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: p.StartLocation,
		EndLocation: p.Expression.Location().EndLocation}
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) statementNode()  {}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: i.Left.Location().StartLocation,
		EndLocation: i.Right.Location().EndLocation}
}

func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) statementNode()  {}

type CallExpression struct {
	function  Expression
	arguments []Expression

	// location of ')'
	endLocation *utils.CodePointLocation
}

func (c *CallExpression) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: c.function.Location().StartLocation,
		EndLocation: c.endLocation}
}

func (c *CallExpression) expressionNode() {}

type Name struct {
	location *utils.CodeBlockLocation
	Name     string
}

func (n *Name) Location() *utils.CodeBlockLocation { return n.location }
func (n *Name) expressionNode()                    {}

type BooleanLiteral struct {
	TokenLocation *utils.CodeBlockLocation
	Value         bool
}

func (b *BooleanLiteral) Location() *utils.CodeBlockLocation { return b.TokenLocation }
func (b *BooleanLiteral) expressionNode()                    {}
func (b *BooleanLiteral) statementNode()                     {}

type StringLiteral struct {
	TokenLocation *utils.CodeBlockLocation
	Value         string
}

func (s *StringLiteral) Location() *utils.CodeBlockLocation { return s.TokenLocation }
func (s *StringLiteral) expressionNode()                    {}
func (s *StringLiteral) statementNode()                     {}

type ArrayLiteral struct {
	location *utils.CodeBlockLocation
	Elements []Expression
}

func (a *ArrayLiteral) Location() *utils.CodeBlockLocation { return a.location }
func (a *ArrayLiteral) expressionNode()                    {}

type IndexExpression struct {
	Left  Expression
	Index Expression

	// location of ']'
	EndLocation *utils.CodePointLocation
}

func (i *IndexExpression) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: i.Left.Location().StartLocation,
		EndLocation: i.EndLocation}
}

func (i *IndexExpression) expressionNode() {}
func (i *IndexExpression) statementNode()  {}

type MapLiteral struct {
	location *utils.CodeBlockLocation

	KeyType   Type
	ValueType Type
	Pairs     map[Expression]Expression
}

func (m *MapLiteral) Location() *utils.CodeBlockLocation { return m.location }
func (m *MapLiteral) expressionNode()                    {}

type PrimaryType struct {
	Token lexer.Token
}

func (p *PrimaryType) Location() *utils.CodeBlockLocation { return p.Token.Location }
func (p *PrimaryType) typeNode()                          {}

type PointerType struct {
	startLocation *utils.CodePointLocation
	Type          Type
}

func (p *PointerType) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: p.startLocation,
		EndLocation: p.Type.Location().EndLocation}
}

func (p *PointerType) typeNode() {}

type ArrayType struct {
	startLocation *utils.CodePointLocation
	Type          Type
}

func (a *ArrayType) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: a.startLocation,
		EndLocation: a.Type.Location().EndLocation}
}

func (a *ArrayType) typeNode() {}
