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

import (
	"fmt"
	"os"
)

type CodeProblemHandler struct {
	Ok       bool
	problems []*CodeProblem
}

func NewCodeProblemHandler() *CodeProblemHandler {
	return &CodeProblemHandler{
		Ok:       true,
		problems: []*CodeProblem{},
	}
}

func (h *CodeProblemHandler) AddCodeProblem(problem *CodeProblem) {
	if problem.critical {
		h.Ok = false
	}

	h.problems = append(h.problems, problem)
}

func (h *CodeProblemHandler) printError(problem *CodeProblem) {
	if problem.global {
		fmt.Fprintf(os.Stderr, "💥 err[%d]: %s\n",
			problem.code,
			fmt.Sprintf(error_messages[problem.code], problem.ctx...))
	} else {
		fmt.Fprintf(os.Stderr, "💥 err[%d] at `%s` (%d:%d-%d:%d): %s\n",
			problem.code, problem.location.StartLocation.Filepath,
			problem.location.StartLocation.Line, problem.location.StartLocation.Column,
			problem.location.EndLocation.Line, problem.location.EndLocation.Column,
			fmt.Sprintf(error_messages[problem.code], problem.ctx...))
	}
}

func (h *CodeProblemHandler) printWarning(problem *CodeProblem) {
	if problem.global {
		fmt.Fprintf(os.Stderr, "⚠ warn[%d]: %s\n",
			problem.code,
			fmt.Sprintf(warning_messages[problem.code], problem.ctx...))
	} else {
		fmt.Fprintf(os.Stderr, "⚠ warn[%d] at `%s` (%d:%d-%d:%d): %s\n",
			problem.code, problem.location.StartLocation.Filepath,
			problem.location.StartLocation.Line, problem.location.StartLocation.Column,
			problem.location.EndLocation.Line, problem.location.EndLocation.Column,
			fmt.Sprintf(warning_messages[problem.code], problem.ctx...))
	}
}

func (h *CodeProblemHandler) printProblem(problem *CodeProblem) {
	if problem.critical {
		h.printError(problem)
	} else {
		h.printWarning(problem)
	}
}

func (h *CodeProblemHandler) PrintProblems() {
	for _, problem := range h.problems {
		h.printProblem(problem)
	}
}

func (h *CodeProblemHandler) PrintDiagnostics() {
	h.PrintProblems()

	fmt.Fprintf(os.Stderr, "\n")
	if !h.Ok {
		fmt.Fprintf(os.Stderr, "error: aborting due to previous error(-s)\n")
	}
}
