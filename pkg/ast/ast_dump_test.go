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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinylang-org/tiny/pkg/lexer"
	"github.com/tinylang-org/tiny/pkg/utils"
)

func TestProgramUnit1(t *testing.T) {
	ast := ProgramUnit{Filepath: "test.tl", Namespace: &NamespaceDecl{
		Name: "test",
		NamespaceLocation: &utils.CodeBlockLocation{
			StartLocation: &utils.CodePointLocation{Filepath: "test.tl", Index: 0, Line: 1, Column: 0},
			EndLocation:   &utils.CodePointLocation{Filepath: "test.tl", Index: 8, Line: 1, Column: 8},
		},
	}}
	assert.Equal(t, "ProgramUnit(\n"+
		"\tfilepath=\"test.tl\",\n"+
		"\tnamespace=Namespace(name=test),\n"+
		"\timports=[]\n"+
		")", ast.Dump(0))
}

func TestProgramUnit2(t *testing.T) {
	ast := ProgramUnit{Filepath: "test.tl", Namespace: &NamespaceDecl{
		Name: "test",
		NamespaceLocation: &utils.CodeBlockLocation{
			StartLocation: &utils.CodePointLocation{Filepath: "test.tl", Index: 0, Line: 1, Column: 0},
			EndLocation:   &utils.CodePointLocation{Filepath: "test.tl", Index: 8, Line: 1, Column: 8},
		},
	}, Imports: []*Import{{
		Path: "test", ImportLocation: &utils.CodeBlockLocation{
			StartLocation: &utils.CodePointLocation{Filepath: "test.tl", Index: 25, Line: 1, Column: 25},
			EndLocation:   &utils.CodePointLocation{Filepath: "test.tl", Index: 32, Line: 1, Column: 32},
		}}, {
		Path: "test2", ImportLocation: &utils.CodeBlockLocation{
			StartLocation: &utils.CodePointLocation{Filepath: "test.tl", Index: 40, Line: 1, Column: 40},
			EndLocation:   &utils.CodePointLocation{Filepath: "test.tl", Index: 48, Line: 1, Column: 48},
		}}}}

	assert.Equal(t, `ProgramUnit(
	filepath="test.tl",
	namespace=Namespace(name=test),
	imports=[
		Import(
			path="test",
			location=CBLocation(CPLocation(test.tl 25 1 25) CPLocation(test.tl 32 1 32))
		),
		Import(
			path="test2",
			location=CBLocation(CPLocation(test.tl 40 1 40) CPLocation(test.tl 48 1 48))
		)
	]
)`, ast.Dump(0))
}

func TestFunctionArgument(t *testing.T) {
	ast := &FunctionArgument{
		Name: "a",
		Type: &PrimaryType{
			Token: &lexer.Token{
				Kind:    lexer.Int16KeywordTokenKind,
				Literal: "int16",
				Location: &utils.CodeBlockLocation{
					StartLocation: &utils.CodePointLocation{
						Filepath: "test.tl",
						Index:    2,
						Line:     1,
						Column:   2,
					},
					EndLocation: &utils.CodePointLocation{
						Filepath: "test.tl",
						Index:    7,
						Line:     1,
						Column:   7,
					},
				},
			},
		},
		BlockLocation: &utils.CodeBlockLocation{
			StartLocation: &utils.CodePointLocation{
				Filepath: "test.tl",
				Index:    0,
				Line:     1,
				Column:   0,
			},
			EndLocation: &utils.CodePointLocation{
				Filepath: "test.tl",
				Index:    7,
				Line:     1,
				Column:   7,
			},
		},
	}

	assert.Equal(t, `FunctionArgument(
	name=a,
	type=PrimaryType(int16),
	location=CBLocation(CPLocation(test.tl 0 1 0) CPLocation(test.tl 7 1 7))
)`, ast.Dump(0))
}
