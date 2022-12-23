# Function declaration
> ```ebnf
> function_declaration = "pub" "fun" identifier "(" function_params ")" ":" type "{"
>                        { statement } "}" |
>                        "fun" identifier "(" function_params ")" ":" type "{"
>                        { statement } "}" .
>												 
> ```

Function declaration can mark function as public or private:

Public function (`pub` keyword is used):
```tiny
pub enum MyLibError { DivisionError }

// public function
pub fun division(a: int32, b: int32): result<MyLibError, int32> {
	if b == 0 {
		return Err(MyLibError.DivisionError)
	}
	return Ok(a / b)
}
```

Private function (`pub` keyword is not used):
```tiny
// private function
fun sum(a: int32, b: int32): int32 {
	return a + b
}
```

Then function name is written. Function name consists of one identifier token. So this one won't be compiled:
```tiny
pub fun $G() {}
pub fun bruh.test() {}
```

After that, function parameters:

> ```ebnf
> function_params      = { identifier ":" type "," } identifier ":" type .
> ```