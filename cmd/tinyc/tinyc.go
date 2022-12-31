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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tinylang-org/tiny/pkg/lexer"
	"github.com/tinylang-org/tiny/pkg/parser"
	"github.com/tinylang-org/tiny/pkg/repr"
	"github.com/tinylang-org/tiny/pkg/utils"
)

var parserPromptCmd = &cobra.Command{
	Use:   "parserprompt",
	Short: "Prompt for testing parser",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimRight(line, "\r\n")

			lineBytes := []byte(line)

			lineBytes = lineBytes[:len(lineBytes)-1]

			ph := utils.NewCodeProblemHandler()
			ph.SetSource(lineBytes)

			p := parser.NewParser("<repl>", lineBytes, ph)
			unit := p.ParseProgramUnit()
			if unit != nil {
				repr.Println(unit)
			}

			ph.SetLineStartOffsets(p.Lexer.LineStartOffsets)
			ph.SetLineEndOffsets(p.Lexer.LineEndOffsets)
			ph.SetColorfulOutput()
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
			line, _ := reader.ReadString('\n')

			line = strings.TrimRight(line, "\r\n")
			lineBytes := []byte(line)

			ph := utils.NewCodeProblemHandler()
			ph.SetSource(lineBytes)

			l := lexer.NewLexer("<repl>", lineBytes, ph)
			for {
				token := l.NextToken()

				fmt.Println(token.Dump())

				if token.Kind == lexer.EOFTokenKind {
					break
				}
			}

			ph.SetLineStartOffsets(l.LineStartOffsets)
			ph.SetLineEndOffsets(l.LineEndOffsets)
			ph.SetColorfulOutput()
			ph.PrintProblems()
		}
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: fmt.Sprintf(`To load completions:

Bash:

  $ source <(%[1]s completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
  # macOS:
  $ %[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ %[1]s completion zsh > "${fpath[1]}/_%[1]s"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ %[1]s completion fish | source

  # To load completions for each session, execute once:
  $ %[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

PowerShell:

  PS> %[1]s completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> %[1]s completion powershell > %[1]s.ps1
  # and source this file from your PowerShell profile.
`, "completion"),
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

var lexCmd = &cobra.Command{
	Use:   "lex",
	Short: "Lexer",
	Run: func(cmd *cobra.Command, args []string) {
		gh := utils.NewCodeProblemHandler()

		if len(args) != 1 {
			fmt.Println("required format: tinyc lex <filename>")
			os.Exit(1)
		}

		fileContent, err := ioutil.ReadFile(args[0])
		if err != nil {
			gh.AddCodeProblem(utils.NewGlobalError(utils.UnableToReadFileErr, []interface{}{args[0]}))
			gh.PrintDiagnostics()
			os.Exit(1)
		}

		ph := utils.NewCodeProblemHandler()
		ph.SetSource(fileContent)
		l := lexer.NewLexer(args[0], fileContent, ph)
		for {
			token := l.NextToken()

			fmt.Println(token.Dump())

			if token.Kind == lexer.EOFTokenKind {
				break
			}
		}

		ph.SetLineStartOffsets(l.LineStartOffsets)
		ph.SetLineEndOffsets(l.LineEndOffsets)
		ph.SetColorfulOutput()
		ph.PrintDiagnostics()

		if !ph.Ok {
			os.Exit(1)
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "tinyc",
	Short: "Compiler for tiny programming language",
}

func main() {
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(lexPromptCmd)
	rootCmd.AddCommand(lexCmd)
	rootCmd.AddCommand(parserPromptCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
