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
	"github.com/vertexgmd/tinylang/pkg/utils"
)

func TestProgramUnit1(t *testing.T) {
	ast := ProgramUnit{Filepath: "test.tl", Namespace: &NamespaceDecl{
		Name: "test",
		NamespaceLocation: &utils.CodeBlockLocation{
			StartLocation: &utils.CodePointLocation{Index: 0, Line: 1, Column: 0},
			EndLocation:   &utils.CodePointLocation{Index: 8, Line: 1, Column: 8},
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
			StartLocation: &utils.CodePointLocation{Index: 0, Line: 1, Column: 0},
			EndLocation:   &utils.CodePointLocation{Index: 8, Line: 1, Column: 8},
		},
	}, Imports: []*Import{{
		Path: "test", ImportLocation: &utils.CodeBlockLocation{
			StartLocation: &utils.CodePointLocation{Index: 1, Line: 1, Column: 1},
			EndLocation:   &utils.CodePointLocation{Index: 2, Line: 1, Column: 2},
		}}, {
		Path: "test2", ImportLocation: &utils.CodeBlockLocation{
			StartLocation: &utils.CodePointLocation{Index: 1, Line: 1, Column: 1},
			EndLocation:   &utils.CodePointLocation{Index: 2, Line: 1, Column: 2},
		}}}}

	assert.Equal(t, "ProgramUnit(\n"+
		"\tfilepath=\"test.tl\",\n"+
		"\tnamespace=Namespace(name=test),\n"+
		"\timports=[Import(\n"+
		")", ast.Dump(0))
}
