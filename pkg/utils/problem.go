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

// CodeProblem describes error/warning happened to be in code.
type CodeProblem struct {
	// If true, it is error, if false, it is warning.
	critical bool

	// Location of the problem.
	location *CodeBlockLocation

	// If error is global (global = true) it is not attached to concrete location and file.
	global bool

	// Code of the problem.
	code int

	// Context. (vargs)
	ctx []interface{}
}

func NewLocalProblem(global bool, critical bool, location *CodeBlockLocation,
	code int, ctx []interface{}) *CodeProblem {
	if global {
		return &CodeProblem{
			critical: critical,
			code:     code,
			ctx:      ctx,
			global:   true,
		}
	} else {
		return &CodeProblem{
			critical: critical,
			location: location,
			code:     code,
			ctx:      ctx,
			global:   false,
		}
	}
}

func NewLocalWarning(location *CodeBlockLocation, code int, ctx []interface{}) *CodeProblem {
	return NewLocalProblem(false, false, location, code, ctx)
}

func NewLocalError(location *CodeBlockLocation, code int, ctx []interface{}) *CodeProblem {
	return NewLocalProblem(false, true, location, code, ctx)
}

func NewGlobalWarning(code int, ctx []interface{}) *CodeProblem {
	return NewLocalProblem(true, false, nil, code, ctx)
}

func NewGlobalError(code int, ctx []interface{}) *CodeProblem {
	return NewLocalProblem(true, true, nil, code, ctx)
}
