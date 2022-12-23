# Import top level statement

> ```ebnf
> import_top_level_statement = import_keyword string_token .
> ```

Import top level statements are located in the start of each program unit:
```tiny
import "./test.tl"

pub fun main() {
	test.printHello()
}
```

If you want to include `$TINYPATH` you can use dollar sign:
```tiny
import "$std/io.tl"

pub fun main() ? {
	var fileContent = handle io.readFile("test.txt")
	print("%s\n", fileContent)
}
```