# Identifiers
Identifiers name entities such as variables, structs, types, .... Identifier is a sequence of letters and digits. The first character of identifier must be a letter. Unicode letters and digits are supported as well:
```
identifier = unicode_letter {unicode_letter | unicode_digit } .
```
Some of this sequences are predeclared as [keywords](./keywords.md) and [boolean literals](./boolean_literals.md).
Examples of identifier tokens:
```tiny
a
_bruuh2008
CoolVariable
αβGreek
中文
```

# Wrapped identifier syntax
If you can also name a variable with sequence of any unicode characters (including whitespaces), via '`' (grave) symbol:
```
wrapped_identifier = '`' {anything_except_grave} '`' .
```
Example:
```tiny
var `cool variable`: string;
```