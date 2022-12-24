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

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertexgmd/tinylang/pkg/lexer"
	"github.com/vertexgmd/tinylang/pkg/parser"
	"github.com/vertexgmd/tinylang/pkg/utils"
)

var parserPromptCmd = &cobra.Command{
	Use:   "parserprompt",
	Short: "Prompt for testing parser",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}

			ph := utils.NewCodeProblemHandler()
			p := parser.NewParser("<repl>", []byte(line), ph)
			unit := p.ParseProgramUnit()
			if unit != nil {
				fmt.Println(unit.Dump(0))
			}

			ph.PrintProblems()
		}
	},
}

var lexPromptCmd = &cobra.Command{
	Use:   "lexprompt",
	Short: "Prompt for testing lexer",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}

			ph := utils.NewCodeProblemHandler()
			l := lexer.NewLexer("<repl>", []byte(line), ph)
			for {
				token := l.NextToken()

				fmt.Println(token.Dump())

				if token.Kind == lexer.EOFTokenKind {
					break
				}
			}

			ph.PrintProblems()
		}
	},
}

var lexCmd = &cobra.Command{
	Use:   "lex",
	Short: "Lexer",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("required format: spline lex <filename>")
			os.Exit(1)
		}

		fileContent, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}

		ph := utils.NewCodeProblemHandler()
		l := lexer.NewLexer(args[0], fileContent, ph)
		for {
			token := l.NextToken()

			fmt.Println(token.Dump())

			if token.Kind == lexer.EOFTokenKind {
				break
			}
		}

		ph.PrintProblems()
	},
}

var rootCmd = &cobra.Command{
	Use:   "tinyc",
	Short: "Compiler for tiny programming language",
}

func main() {
	rootCmd.AddCommand(lexPromptCmd)
	rootCmd.AddCommand(lexCmd)
	rootCmd.AddCommand(parserPromptCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
