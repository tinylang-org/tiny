Prism.languages.tiny = Prism.languages.extend("clike", {
	keyword:
	  /\b(?:handle|enum|result|break|case|const|continue|string|default|else|for|fun|return|pub|if|import|bool|int16|int32|int64|int8|namespace|struct|switch|uint16|uint32|uint64|uint8|var)\b/,
	builtin: /\b(?:null|self)\b/,
	boolean: /\b(?:true|false)\b/,
	operator:
	  /(==|!=|<=|>=|<|>|&&|\|\||!|=|\+\=|\-\=|\*\=|\/\=|\+|\-|\*|\/|%|\^|\.\.|\-\-|\+\+|\/\%|\/\%=|\^=|<=>|\||\&)/,
	// number: /(?:\b\d+(\.\d+)?\b)|(\b([0-9]+|\?)[gbci]\b)/,
	string: /[a-z]?"(?:\\.|[^\\"])*"|'(?:\\.|[^\\'])*'/,
	tag: /@([a-zA-Z_][a-zA-Z0-9_]*)/,
  })

delete Prism.languages.tiny["class-name"]
