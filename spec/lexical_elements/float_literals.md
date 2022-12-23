# Floating-point literals
Float-point literals can be either decimal, or hexadecimal.

EBNF grammar:
```ebnf
floating_point_literal        = decimal_float_literal | hexadecimal_float_literal .

decimal_float_literal         = decimal_digits_sequence "." [ decimal_digits_sequence ] [ decimal_exponent ] |
                                decimal_digits_sequence decimal_exponent |
                                "." decimal_digits_sequence [ decimal_exponent ] .
decimal_exponent              = ( "e" | "E" ) [ "+" | "-" ] decimal_digits_sequence .

hexadecimal_float_literal     = "0" ( "x" | "X" ) hexadecimal_mantissa hexadecimal_exponent .
hexadecimal_mantissa          = [ "_" ] hexadecimal_digits_sequence "." [         hexadecimal_digits_sequence ] |
                                [ "_" ] hexadecimal_digits_sequence |
                                "." hexadecimal_digits_sequence .
hexadecimal_exponent          = ( "p" | "P" ) [ "+" | "-" ] decimal_digits_sequence .
```