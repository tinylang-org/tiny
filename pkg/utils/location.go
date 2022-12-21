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

import "fmt"

type CodePointLocation struct {
	Filepath string
	Index    int
	Line     int
	Column   int
}

type CodeBlockLocation struct {
	StartLocation *CodePointLocation
	EndLocation   *CodePointLocation
}

func (l *CodePointLocation) Copy() *CodePointLocation {
	return &CodePointLocation{l.Filepath, l.Index, l.Line, l.Column}
}

func (l *CodePointLocation) NextByteLocation() *CodePointLocation {
	return &CodePointLocation{l.Filepath, l.Index + 1, l.Line, l.Column + 1}
}

func (l *CodePointLocation) PreviousByteLocation() *CodePointLocation {
	return &CodePointLocation{l.Filepath, l.Index - 1, l.Line, l.Column - 1}
}

func NewOneCodePointBlockLocation(l *CodePointLocation) *CodeBlockLocation {
	return &CodeBlockLocation{
		l.Copy(),
		l.NextByteLocation(),
	}
}

func NewTwoCodePointsBlockLocation(l *CodePointLocation) *CodeBlockLocation {
	return &CodeBlockLocation{
		l.PreviousByteLocation(),
		l.NextByteLocation(),
	}
}

func (l *CodePointLocation) Dump() string {
	return fmt.Sprintf("CPLocation(%s %d %d %d)", l.Filepath, l.Index, l.Line, l.Column)
}

func (l *CodeBlockLocation) Dump() string {
	return fmt.Sprintf("CBLocation(%s %s)", l.StartLocation.Dump(),
		l.EndLocation.Dump())
}
