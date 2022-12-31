# The Tiny Programming Language Specification
### Version of December 28, 2022

<table>
<tr><td width=33% valign=top>

- [Introduction](#introduction)
- [Notation](#notation)
- [Basic source code elements](#basic-source-code-elements)
  - [Characters](#characters)
  - [Letters and digits](#letters-and-digits)
- [Tokens](#tokens)
  - [Comments](#comments)
  - [Semicolons](#semicolons)
  - [Identifiers](#identifiers)
  - [Keywords](#keywords)
  - [Punctuators and operators](#punctuators-and-operators)
  - [String literals]()
  - [Numbers]()
    - [Integer literals]()
    - [Floating-point literals]()
    - [Imaginary number literals]()

</td><td width=33% valign=top>

- [Variables]()
- [Type system]()
  - [Boolean type]()
  - [Numeric type]()
  - [Pointer type]()
  - [Struct type]()
  - [Function type]()
  - [Interface type]()
- [Basic syntax]()
  - [Top-level statements]()
    - [Imports]()
    - [Function declarations]()
    - [Structure declarations]()
  - [Statements]()
    - [Expression statements]()
    - [Variable declaration statement]()
    - [If-else statements]()
    - [Switch statements]()
    - [For statements]()
    - [Return statements]()
    - [Break statements]()
    - [Return statements]()
    - [Continue statements]()


</td><td valign=top>

- [Memory managment system]()
  - [Stack and heap]()
  - [Heap allocations `new` and `destroy`]()
  - [Memory managment and OOP]()
- [Error handling]()
  - [Default error handler `handle`]()
  - [Custom error handlers]()

</td></tr>
</table>

## Introduction

This is the reference manual for the Tiny programming language. 

Tiny is a language designed with embeded programming in mind. It is strongly typed and has manual memory managment system.

The syntax is compact and simple to parse, allowing for easy analysis by automatic tools such as integrated development environments and fast compilation times.

Tiny is something between C and C++. It has more less abstractions than in C++, but at the same time allows developers to write easy-to-read and scalable code.

Here is an example of hello world program in Tiny:

```tiny
pub fun main() {
    printf("hello world\n");
}
```

## Notation
The syntax is specified using a variant of Extended Backus-Naur Form(EBNF):
```ebnf
Syntax      = { Production } .
Production  = production_name "=" [ Expression ] "." .
Expression  = Term { "|" Term } .
Term        = Factor { Factor } .
Factor      = production_name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .
```

Productions are expressions constructed from terms and the following operators, in increasing precedence:

```
|   alternation
()  grouping
[]  option (0 or 1 times)
{}  repetition (0 to n times)
```

Lowercase production names are used to identify lexical (terminal) tokens. Non-terminals are in CamelCase. Lexical tokens are enclosed in double quotes "" or back quotes ``.

The form a … b represents the set of characters from a through b as alternatives. The horizontal ellipsis … is also used elsewhere in the spec to informally denote various enumerations or code snippets that are not further specified. The character … (as opposed to the three characters ...) is not a token of the Tiny language.

## Basic source code elements

Source code is Unicode text encoded in UTF-8. The text is not canonicalized, so a single accented code point is distinct from the same character constructed from combining an accent and a letter; those are treated as two code points. For simplicity, this document will use the unqualified term character to refer to a Unicode code point in the source text.

Each code point is distinct; for instance, uppercase and lowercase letters are different characters.

Implementation restriction: For compatibility with other tools, a compiler may disallow the NUL character (U+0000) in the source text.

### Characters

The following terms are used to denote specific Unicode character categories:

```ebnf
newline        = /* the Unicode code point U+000A */ .
unicode_char   = /* an arbitrary Unicode code point except newline */ .
unicode_letter = /* a Unicode code point categorized as "Letter" */ .
unicode_digit  = /* a Unicode code point categorized as "Number, decimal digit" */ .
```

In [The Unicode Standard 8.0](https://www.unicode.org/versions/Unicode8.0.0/), Section 4.5 "General Category" defines a set of character categories. Tiny treats all characters in any of the Letter categories Lu, Ll, Lt, Lm, or Lo as Unicode letters, and those in the Number category Nd as Unicode digits.

### Letters and digits

The underscore character `_` (U+005F) is considered a lowercase letter.

```ebnf
letter        = unicode_letter | "_" .
decimal_digit = "0" … "9" .
binary_digit  = "0" | "1" .
octal_digit   = "0" … "7" .
hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
```

## Tokens

### Comments

Comments serve as program documentation. There are two forms:

1. Line comments start with the character sequence `//` and stop at the end of the line. Example:
```tiny
// This is recursive implementation of factorial :3.
pub fun factorial(n: f64): f64 {
    if (n < 2) return 1;
    return factorial(n - 1) * n;
}
```

1. Multiline comments start with the character sequence `/*` and stop with the first subsequent character sequence `*/`. Example:

```tiny
/**
 * @param a first number
 * @param b second number
 *
 * @return maximum number of numbers a and b
 */
pub fun max(a: f64, b: f64): f64 {
    if (a > b) return a;
    return b;
}
```

> A comment **cannot start** inside **a character or string literal**, or **inside a comment**.

## Semicolons
The formal syntax uses semicolons ";" as terminators in a number of productions. Tiny programs are required to have ";" at the end of each statement. Compiler does **NOT** emit them automatically.

Example of wrong Tiny program:

```tiny
pub fun printNumber(n: f64) {
    printf("%f", n) // no semicolon here. syntax error
}
```

## Identifiers
Identifiers name program entities such as variables and types. An identifier is a sequence of one or more letters and digits. The first character in an identifier must be a letter.

```ebnf
identifier = unicode_letter { unicode_letter | unicode_digit } .
```

Here are some examples of valid identifiers:
```tiny
test_identifier
название22
_x15
a
someVariable
```

## Keywords
The following keywords are reserved and may not be used as identifiers.

```tiny
break       default     fun     interface       case
struct      else        switch  const           if
i8          i16         i32     i64
u8          u16         u32     u64
continue    for         import  return          var
```

## Punctuators and operators
The following character sequences represent operators and punctuators:
```tiny
+    &     +=    &=     &&    ==    !=    (    )
-    |     -=    |=     ||    <     <=    [    ]
/    <<    /=    <<=    ++    =     ,     ;    ~
%    >>    %=    >>=    --    !     ...   .    :
*    ^     *=    ^=     >     >=    {     }    ?
```

```ebnf
char_lit      = "'" (  )
```