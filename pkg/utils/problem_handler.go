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
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

type CodeProblemHandler struct {
	Ok             bool
	colorfulOutput bool

	source       []byte
	sourceLength int

	problems []*CodeProblem

	lineStartOffsets *[]int
	lineEndOffsets   *[]int
}

func NewCodeProblemHandler() *CodeProblemHandler {
	return &CodeProblemHandler{
		Ok:             true,
		colorfulOutput: false,
		source:         []byte(""),
		sourceLength:   0,
		problems:       []*CodeProblem{},

		lineStartOffsets: &[]int{},
		lineEndOffsets:   &[]int{},
	}
}

func (h *CodeProblemHandler) SetColorfulOutput() {
	h.colorfulOutput = true
}

func (h *CodeProblemHandler) SetLineStartOffsets(lineStartOffsets *[]int) {
	h.lineStartOffsets = lineStartOffsets
}

func (h *CodeProblemHandler) SetLineEndOffsets(lineEndOffsets *[]int) {
	h.lineEndOffsets = lineEndOffsets
}

func (h *CodeProblemHandler) AddCodeProblem(problem *CodeProblem) {
	if problem.critical {
		h.Ok = false
	}

	h.problems = append(h.problems, problem)
}

func (h *CodeProblemHandler) SetSource(source []byte) {
	h.source = source
	h.sourceLength = len(source)
}

func (h *CodeProblemHandler) printFormattedCodeBlock(
	start int, end int, c *color.Color) {
	s := strings.Replace(
		string(
			h.source[start:end],
		),
		"\t",
		" ",
		-1,
	)

	if h.colorfulOutput {
		c.Fprint(os.Stderr, s)
	} else {
		fmt.Fprint(os.Stderr, s)
	}
}

func (h *CodeProblemHandler) printCodeBlock(problem *CodeProblem) {
	lineNumberStr := fmt.Sprintf("%d", problem.location.StartLocation.Line)
	lineNumberStrLength := len(lineNumberStr)

	spacesBeforeBar := ""
	for i := 0; i < lineNumberStrLength+2; i++ {
		spacesBeforeBar += " "
	}

	fmt.Fprint(os.Stderr, spacesBeforeBar)
	fmt.Fprint(os.Stderr, "|\n")
	fmt.Fprintf(os.Stderr, " %d | ", problem.location.StartLocation.Line)

	w := color.New(color.FgWhite)

	var prc *color.Color
	if problem.critical {
		prc = color.New(color.FgRed, color.Bold)
	} else {
		prc = color.New(color.FgYellow)
	}

	h.printFormattedCodeBlock(
		(*h.lineStartOffsets)[problem.location.StartLocation.Line-1],
		problem.location.StartLocation.Index,
		w)

	h.printFormattedCodeBlock(
		problem.location.StartLocation.Index,
		problem.location.EndLocation.Index,
		prc)

	if problem.location.StartLocation.Index < h.sourceLength {
		h.printFormattedCodeBlock(
			problem.location.EndLocation.Index,
			(*h.lineEndOffsets)[problem.location.EndLocation.Line-1]+1,
			w)
	}

	fmt.Fprint(os.Stderr, "\n")

	spaceBeforeArrow := ""
	for i := 0; i < problem.location.StartLocation.Column+1; i++ {
		spaceBeforeArrow += " "
	}

	fmt.Fprint(os.Stderr, spacesBeforeBar)
	fmt.Fprint(os.Stderr, "|")
	fmt.Fprint(os.Stderr, spaceBeforeArrow)

	tildaSymbols := ""
	for i := 0; i < problem.location.EndLocation.Column-problem.location.StartLocation.Column-1; i++ {
		tildaSymbols += "~"
	}

	if !h.colorfulOutput {
		fmt.Fprint(os.Stderr, "^")
		fmt.Fprint(os.Stderr, tildaSymbols)
	} else {
		var c *color.Color
		if problem.critical {
			c = color.New(color.FgRed, color.Bold)
		} else {
			c = color.New(color.FgYellow)
		}

		c.Fprint(os.Stderr, "^")
		c.Fprint(os.Stderr, tildaSymbols)
	}

	fmt.Fprintf(os.Stderr, "\n\n")
}

func readLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

func (h *CodeProblemHandler) printError(problem *CodeProblem) {
	if problem.global {
		if h.colorfulOutput {
			color.New(color.FgRed, color.Bold).Fprint(os.Stderr, "error:")
		} else {
			fmt.Fprint(os.Stderr, "error:")
		}

		fmt.Fprintf(os.Stderr, " %s\n",
			fmt.Sprintf(error_messages[problem.code], problem.ctx...))
	} else {
		fmt.Fprintf(os.Stderr, "%s(%d:%d) ",
			problem.location.StartLocation.Filepath,
			problem.location.StartLocation.Line, problem.location.StartLocation.Column)

		if h.colorfulOutput {
			color.New(color.FgRed, color.Bold).Fprint(os.Stderr, "error:")
		} else {
			fmt.Fprint(os.Stderr, "error:")
		}

		fmt.Fprintf(os.Stderr, " %s\n",
			fmt.Sprintf(error_messages[problem.code], problem.ctx...))
		h.printCodeBlock(problem)
	}
}

func (h *CodeProblemHandler) printWarning(problem *CodeProblem) {
	if problem.global {
		fmt.Fprintf(os.Stderr, "warning: %s\n",
			fmt.Sprintf(warning_messages[problem.code], problem.ctx...))
	} else {
		fmt.Fprintf(os.Stderr, "%s(%d:%d) warning: %s\n",
			problem.location.StartLocation.Filepath,
			problem.location.StartLocation.Line, problem.location.StartLocation.Column,
			fmt.Sprintf(warning_messages[problem.code], problem.ctx...))
		h.printCodeBlock(problem)
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

	if !h.Ok {
		if h.colorfulOutput {
			color.New(color.FgRed, color.Bold).Fprintln(
				os.Stderr, "error: aborting due to previous error(-s)")
		} else {
			fmt.Fprintln(os.Stderr, "error: aborting due to previous error(-s)")
		}
	}
}
