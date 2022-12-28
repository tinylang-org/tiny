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

// Package ast declares the types used to represent syntax trees for Tiny source text.
package ast

import (
	"fmt"
	"strings"

	"github.com/tinylang-org/tiny/pkg/lexer"
	"github.com/tinylang-org/tiny/pkg/utils"
)

type AST interface {
	Location() *utils.CodeBlockLocation
	Dump(identationLevel int) string
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

type Import struct {
	Path           string
	ImportLocation *utils.CodeBlockLocation
}

func (i *Import) Location() *utils.CodeBlockLocation { return i.ImportLocation }
func (i *Import) Dump(identationLevel int) string {
	var sb strings.Builder

	addIdentation(identationLevel, &sb)
	sb.WriteString("Import(\n")

	identationLevel++
	addIdentation(identationLevel, &sb)
	sb.WriteString(fmt.Sprintf("path=\"%s\",\n", i.Path))
	addIdentation(identationLevel, &sb)
	sb.WriteString(fmt.Sprintf("location=%s\n", i.ImportLocation.Dump()))
	identationLevel--
	addIdentation(identationLevel, &sb)
	sb.WriteString(")")

	return sb.String()
}

type NamespaceDecl struct {
	Name              string
	NamespaceLocation *utils.CodeBlockLocation
}

func (n *NamespaceDecl) Location() *utils.CodeBlockLocation { return n.NamespaceLocation }
func (n *NamespaceDecl) Dump(identationLevel int) string {
	return fmt.Sprintf("Namespace(name=%s)", n.Name)
}

type ProgramUnit struct {
	Filepath     string
	Namespace    *NamespaceDecl
	Imports      []*Import
	TLStatements []TopLevelStatement
}

func (p *ProgramUnit) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: p.Namespace.Location().StartLocation,
		EndLocation: p.TLStatements[len(p.TLStatements)-1].Location().EndLocation}
}
func (p *ProgramUnit) Dump(identationLevel int) string {
	var sb strings.Builder
	sb.WriteString("ProgramUnit(\n")

	identationLevel++
	addIdentation(identationLevel, &sb)
	sb.WriteString(fmt.Sprintf("filepath=\"%s\",\n", p.Filepath))
	addIdentation(identationLevel, &sb)
	sb.WriteString(fmt.Sprintf("namespace=%s,\n", p.Namespace.Dump(identationLevel)))
	addIdentation(identationLevel, &sb)

	if len(p.Imports) != 0 {
		sb.WriteString(fmt.Sprintf("imports=[\n"))

		identationLevel++

		for i, imp := range p.Imports {
			sb.WriteString(imp.Dump(identationLevel))

			if i != len(p.Imports)-1 {
				sb.WriteString(",")
			}
			sb.WriteString("\n")
		}

		identationLevel--
		addIdentation(identationLevel, &sb)
		sb.WriteString("]\n")
	} else {
		sb.WriteString(fmt.Sprintf("imports=[]\n"))
	}

	identationLevel--
	addIdentation(identationLevel, &sb)
	sb.WriteString(")")

	return sb.String()
}

type FunctionDeclaration struct {
	BlockLocation   *utils.CodeBlockLocation
	Public          bool
	Name            string
	StatementsBlock *StatementsBlock
	Arguments       []*FunctionArgument
}

func (f *FunctionDeclaration) Location() *utils.CodeBlockLocation { return f.BlockLocation }
func (f *FunctionDeclaration) topLevelStatement()                 {}
func (f *FunctionDeclaration) Dump(identationLevel int) string    { return "" }

type FunctionArgument struct {
	BlockLocation *utils.CodeBlockLocation
	Name          string
	Type          Type
}

func (a *FunctionArgument) Location() *utils.CodeBlockLocation { return a.BlockLocation }
func (a *FunctionArgument) Dump(identationLevel int) string {
	var sb strings.Builder
	addIdentation(identationLevel, &sb)
	sb.WriteString("FunctionArgument(\n")

	identationLevel++
	addIdentation(identationLevel, &sb)
	sb.WriteString(fmt.Sprintf("name=%s,\n", a.Name))

	addIdentation(identationLevel, &sb)
	sb.WriteString(fmt.Sprintf("type=%s,\n", a.Type.Dump(identationLevel)))

	addIdentation(identationLevel, &sb)
	sb.WriteString(fmt.Sprintf("location=%s\n", a.BlockLocation.Dump()))

	identationLevel--

	addIdentation(identationLevel, &sb)
	sb.WriteString(")")

	return sb.String()
}

type StructureDeclaration struct {
	BlockLocation *utils.CodeBlockLocation
	Public        bool
	Name          string
	Functions     *FunctionDeclaration
	Members       []*StructureMember
}

func (s *StructureDeclaration) Location() *utils.CodeBlockLocation { return s.BlockLocation }
func (s *StructureDeclaration) topLevelStatement()                 {}
func (s *StructureDeclaration) Dump(identationLevel int) string    { return "" }

type StructureMember struct {
	BlockLocation *utils.CodeBlockLocation
	Public        bool
	Readonly      bool
	Name          string
	Type          Type
}

func (m *StructureMember) Location() *utils.CodeBlockLocation { return m.BlockLocation }

type StatementsBlock struct {
	// location of '{'
	StartLocation *utils.CodePointLocation
	Statements    []Statement
}

func (s *StatementsBlock) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: s.StartLocation,
		EndLocation: s.Statements[len(s.Statements)-1].Location().EndLocation}
}

func (s *StatementsBlock) Dump(identationLevel int) string { return "" }

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
func (v *VarStatement) statementNode()                  {}
func (v *VarStatement) topLevelStatement()              {}
func (v *VarStatement) Dump(identationLevel int) string { return "" }

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

func (r *ReturnStatement) statementNode()                  {}
func (r *ReturnStatement) Dump(identationLevel int) string { return "" }

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

func (p *PrefixExpression) expressionNode()                 {}
func (p *PrefixExpression) statementNode()                  {}
func (p *PrefixExpression) Dump(identationLevel int) string { return "" }

// infix_expression = left operator right .
type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: i.Left.Location().StartLocation,
		EndLocation: i.Right.Location().EndLocation}
}

func (i *InfixExpression) expressionNode()                 {}
func (i *InfixExpression) statementNode()                  {}
func (i *InfixExpression) Dump(identationLevel int) string { return "" }

type CallExpression struct {
	Function  Expression
	Arguments []Expression

	// location of ')'
	EndLocation *utils.CodePointLocation
}

func (c *CallExpression) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: c.Function.Location().StartLocation,
		EndLocation: c.EndLocation}
}

func (c *CallExpression) expressionNode()                 {}
func (c *CallExpression) statementNode()                  {}
func (c *CallExpression) Dump(identationLevel int) string { return "" }

type Name struct {
	location *utils.CodeBlockLocation
	Name     string
}

func (n *Name) Location() *utils.CodeBlockLocation { return n.location }
func (n *Name) expressionNode()                    {}
func (n *Name) Dump(identationLevel int) string    { return "" }

type BooleanLiteral struct {
	TokenLocation *utils.CodeBlockLocation
	Value         bool
}

func (b *BooleanLiteral) Location() *utils.CodeBlockLocation { return b.TokenLocation }
func (b *BooleanLiteral) expressionNode()                    {}
func (b *BooleanLiteral) statementNode()                     {}
func (b *BooleanLiteral) Dump(identationLevel int) string    { return "" }

type StringLiteral struct {
	TokenLocation *utils.CodeBlockLocation
	Value         string
}

func (s *StringLiteral) Location() *utils.CodeBlockLocation { return s.TokenLocation }
func (s *StringLiteral) expressionNode()                    {}
func (s *StringLiteral) statementNode()                     {}
func (s *StringLiteral) Dump(identationLevel int) string    { return "" }

type ArrayLiteral struct {
	location *utils.CodeBlockLocation
	Elements []Expression
}

func (a *ArrayLiteral) Location() *utils.CodeBlockLocation { return a.location }
func (a *ArrayLiteral) expressionNode()                    {}
func (a *ArrayLiteral) Dump(identationLevel int) string    { return "" }

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

func (i *IndexExpression) expressionNode()                 {}
func (i *IndexExpression) statementNode()                  {}
func (i *IndexExpression) Dump(identationLevel int) string { return "" }

type MapLiteral struct {
	location *utils.CodeBlockLocation

	KeyType   Type
	ValueType Type
	Pairs     map[Expression]Expression
}

func (m *MapLiteral) Location() *utils.CodeBlockLocation { return m.location }
func (m *MapLiteral) expressionNode()                    {}
func (m *MapLiteral) Dump(identationLevel int) string    { return "" }

type PrimaryType struct {
	Token *lexer.Token
}

func (p *PrimaryType) Location() *utils.CodeBlockLocation { return p.Token.Location }
func (p *PrimaryType) typeNode()                          {}
func (p *PrimaryType) Dump(identationLevel int) string {
	d := lexer.DumpTokenKind(p.Token.Kind)
	return fmt.Sprintf(
		"PrimaryType(%s)",
		d[:len(d)-8], // remove "... keyword" at the end
	)
}

type PointerType struct {
	StartLocation *utils.CodePointLocation
	Type          Type
}

func (p *PointerType) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: p.StartLocation,
		EndLocation: p.Type.Location().EndLocation}
}

func (p *PointerType) typeNode() {}

func (p *PointerType) Dump(identationLevel int) string {
	return fmt.Sprintf("PointerType(%s)", p.Type.Dump(identationLevel))
}

type ArrayType struct {
	StartLocation *utils.CodePointLocation
	Type          Type
}

func (a *ArrayType) Location() *utils.CodeBlockLocation {
	return &utils.CodeBlockLocation{StartLocation: a.StartLocation,
		EndLocation: a.Type.Location().EndLocation}
}

func (a *ArrayType) typeNode() {}

func (a *ArrayType) Dump(identationLevel int) string {
	return fmt.Sprintf("ArrayType(%s)", a.Type.Dump(identationLevel))
}

type CustomType struct {
	TypeLocation *utils.CodeBlockLocation
	Name         string
}

func (c *CustomType) Location() *utils.CodeBlockLocation { return c.TypeLocation }
func (c *CustomType) typeNode()                          {}
func (c *CustomType) Dump(identationLevel int) string {
	return fmt.Sprintf("CustomType(%s)", c.Name)
}

func addIdentation(identationLevel int, stringBuilder *strings.Builder) {
	for i := 0; i < identationLevel; i++ {
		stringBuilder.WriteString("\t")
	}
}
