# Integer literal
Integers can be decimal, binary, octal or hexadecimal ones.

Decimal integers start with and contain unicode characters `0`, `1`, `2`, `3`, `4`, `5`, ..., `9`. Examples of ones:

```tiny
0
12
033
30
```

If the token starts with `0b` or `0B`, then it is binary integer and can only contain `0` and `1`-es:
```tiny
0b011001
0b0
0b1
0b32 // incorrect
```

If the token starts with `0o` or `0O`, then it is octal integer and can only contain `0`, `1`, `2`, ...`7` unicode characters in it:
```tiny
0o32 // correct now
0o10023523
0o7
0o8 // incorrect
0o989 // incorrect
```
If the token starts with `0x` or `0X`, then it is octal integer and can only contain  `0`, `1`, `2`, ...`9`, `a`, `b`, `c`, ..., `f`, `A`, `B`, `C`, ..., `F` unicode characters in it:
```tiny
0x987
0xA
0xa
0xZ // incorrect
0xabcdefg // incorrect
0xaBcDef // correct
```

`_` separator for digits can be used (but only for successive ones):
```tiny
0x1_2 // ok
1_829_304_233 // ok
0_x1 // incorrect
```


EBNF grammar:
```ebnf
integer_literal                = decimal_literal | binary_literal | octal_literal | hexadecimal_literal .

decimal_literal                = "0" | ( "1" … "9" ) [ [ "_" ] decimal_digits_sequence ] .
binary_literal                 = "0" ( "b" | "B" ) [ "_" ] binary_digits_sequence .
octal_literal                  = "0" [ "o" | "O" ] [ "_" ] octal_digits_sequence .
hexadecimal_literal            = "0" ( "x" | "X" ) [ "_" ] hexadecimal_digits_sequence .

decimal_digits_sequence         = decimal_digit { [ "_" ] decimal_digit } .
binary_digits_sequence          = binary_digit { [ "_" ] binary_digit } .
octal_digits_sequence           = octal_digit { [ "_" ] octal_digit } .
hexadecimal_digits_sequence     = hexadecimal_digit { [ "_" ] hexadecimal_digit } .
```