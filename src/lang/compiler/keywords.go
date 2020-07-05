package compiler

//list of keywords in json
var keywordJSON =
[]byte(
`
[
  {
    "name": "local",
    "remove": "local",
    "pattern": "(local(\\s*)(~?))",
    "type": "id"
  },
  {
    "name": "local",
    "remove": "lcl",
    "pattern": "(lcl(\\s*)(~?))",
    "type": "id"
  },
  {
    "name": "->",
    "remove": "->",
    "pattern": "(\\-\\>)",
    "type": "operation"
  },
  {
    "name": "=>",
    "remove": "=>",
    "pattern": "(\\=\\>)",
    "type": "operation"
  },
  {
    "name": "++",
    "remove": "++",
    "pattern": "(\\+\\+)",
    "type": "?operation"
  },
  {
    "name": "+=",
    "remove": "+=",
    "pattern": "(\\+\\=)",
    "type": "?operation"
  },
  {
    "name": "--",
    "remove": "--",
    "pattern": "(\\-\\-)",
    "type": "?operation"
  },
  {
    "name": "-=",
    "remove": "-=",
    "pattern": "(\\-\\-)",
    "type": "?operation"
  },
  {
    "name": "*=",
    "remove": "*=",
    "pattern": "(\\*\\=)",
    "type": "?operation"
  },
  {
    "name": "/=",
    "remove": "/=",
    "pattern": "(\\/\\=)",
    "type": "?operation"
  },
  {
    "name": "^=",
    "remove": "^=",
    "pattern": "(\\^\\=)",
    "type": "?operation"
  },
  {
    "name": "%=",
    "remove": "%=",
    "pattern": "(\\%\\=)",
    "type": "?operation"
  },
  {
    "name": "each",
    "remove": "each",
    "pattern": "(each(\\s*)\\()",
    "type": "cond"
  },
  {
    "name": "~~~",
    "remove": "~~~",
    "pattern": "((~~~))",
    "type": "operation"
  },
  {
    "name": "~~",
    "remove": "~~",
    "pattern": "(~~)",
    "type": "operation"
  },
  {
    "name": "alt",
    "remove": "alt",
    "pattern": "(alt\\()",
    "type": "cond"
  },
  {
    "name": "log",
    "remove": "log",
    "pattern": "(log(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "log",
    "remove": "log",
    "pattern": "(log(\\s+))",
    "type": "id"
  },
  {
    "name": "print",
    "remove": "print",
    "pattern": "(print(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "print",
    "remove": "print",
    "pattern": "(print(\\s+))",
    "type": "id"
  },
  {
    "name": "^",
    "remove": "**",
    "pattern": "(\\*\\*)",
    "type": "operation"
  },
  {
    "name": ".",
    "remove": ".>",
    "pattern": "(\\.\\>)",
    "type": "operation"
  },
  {
    "name": "newlineN",
    "remove": "\\n",
    "pattern": "\\n",
    "type": "newline"
  },
  {
    "name": "newlineN",
    "remove": "\\r\\n",
    "pattern": "(\\r\\n)",
    "type": "newline"
  },
  {
    "name": "~",
    "remove": "~",
    "pattern": "\\~",
    "type": "operation"
  },
  {
    "name": "[:",
    "remove": "[:",
    "pattern": "(\\[\\:)",
    "type": "?open_brace"
  },
  {
    "name": ":]",
    "remove": ":]",
    "pattern": "(\\:\\])",
    "type": "?close_brace"
  },
  {
    "name": "::",
    "remove": "::",
    "pattern": "(\\:\\:)",
    "type": "operation"
  },
  {
    "name": ":",
    "remove": ":=",
    "pattern": "(\\:\\=)",
    "type": "operation"
  },
  {
    "name": ":",
    "remove": ":",
    "pattern": "\\:",
    "type": "operation"
  },
  {
    "name": "+",
    "remove": "+",
    "pattern": "(\\+)",
    "type": "operation"
  },
  {
    "name": "-",
    "remove": "-",
    "pattern": "(\\-)",
    "type": "operation"
  },
  {
    "name": "*",
    "remove": "*",
    "pattern": "\\*",
    "type": "operation"
  },
  {
    "name": "/",
    "remove": "/",
    "pattern": "\\/",
    "type": "operation"
  },
  {
    "name": "%",
    "remove": "%",
    "pattern": "\\%",
    "type": "operation"
  },
  {
    "name": "^",
    "remove": "^",
    "pattern": "\\^",
    "type": "operation"
  },
  {
    "name": "(",
    "remove": "(",
    "pattern": "\\(",
    "type": "?open_brace"
  },
  {
    "name": ")",
    "remove": ")",
    "pattern": "\\)",
    "type": "?close_brace"
  },
  {
    "name": "{",
    "remove": "{",
    "pattern": "\\{",
    "type": "?open_brace"
  },
  {
    "name": "}",
    "remove": "}",
    "pattern": "\\}",
    "type": "?close_brace"
  },
  {
    "name": "[",
    "remove": "[",
    "pattern": "\\[",
    "type": "?open_brace"
  },
  {
    "name": "]",
    "remove": "]",
    "pattern": "\\]",
    "type": "?close_brace"
  },
  {
    "name": "async",
    "remove": "async",
    "pattern": "(async\\s*~)",
    "type": "id_non_tilde"
  },
  {
    "name": "async",
    "remove": "async",
    "pattern": "(async\\s+)",
    "type": "id"
  },
  {
    "name": "fargc",
    "remove": "fargc",
    "pattern": "(fargc\\s+)",
    "type": "id"
  },
  {
    "name": "fargc",
    "remove": "fargc",
    "pattern": "(fargc\\s*\\~)",
    "type": "id_non_tilde"
  },
  {
    "name": "function",
    "remove": "fn",
    "pattern": "(fn(\\s*)\\()",
    "type": "id_non_tilde"
  },
  {
    "name": "global",
    "remove": "gbl",
    "pattern": "(gbl(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "global",
    "remove": "gbl",
    "pattern": "(gbl(\\s+))",
    "type": "id"
  },
  {
    "name": ",",
    "remove": ",",
    "pattern": "\\,",
    "type": "operation"
  },
  {
    "name": "await",
    "remove": "await",
    "pattern": "(await(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "await",
    "remove": "await",
    "pattern": "(await(\\s+))",
    "type": "id"
  },
  {
    "name": "return",
    "remove": "return",
    "pattern": "(return(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "return",
    "remove": "return",
    "pattern": "(return(\\s+))",
    "type": "id"
  },
  {
    "name": "if",
    "remove": "if",
    "pattern": "(if(\\s*)\\()",
    "type": "cond"
  },
  {
    "name": "elseif",
    "remove": "elseif",
    "pattern": "(else(\\s*)if(\\s*)\\()",
    "type": "cond"
  },
  {
    "name": "else",
    "remove": "else",
    "pattern": "(else(\\s*)(\\{|\\~))",
    "type": "cond"
  },
  {
    "name": "else",
    "remove": "else",
    "pattern": "(else(\\s+))",
    "type": "id"
  },
  {
    "name": "loop",
    "remove": "loop",
    "pattern": "(loop(\\s*)\\()",
    "type": "cond"
  },
  {
    "name": "<=",
    "remove": "<=",
    "pattern": "(\\<\\=)",
    "type": "operation"
  },
  {
    "name": ">=",
    "remove": ">=",
    "pattern": "(\\>\\=)",
    "type": "operation"
  },
  {
    "name": ">",
    "remove": ">",
    "pattern": "\\>",
    "type": "operation"
  },
  {
    "name": "<",
    "remove": "<",
    "pattern": "\\<",
    "type": "operation"
  },
  {
    "name": "!=",
    "remove": "!=",
    "pattern": "(\\!\\=)",
    "type": "operation"
  },
  {
    "name": "=",
    "remove": "==",
    "pattern": "(\\=\\=)",
    "type": "operation"
  },
  {
    "name": "=",
    "remove": "=",
    "pattern": "\\=",
    "type": "operation"
  },
  {
    "name": "import",
    "remove": "import",
    "pattern": "(import(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "import",
    "remove": "import",
    "pattern": "(import(\\s+))",
    "type": "id"
  },
  {
    "name": "include",
    "remove": "include",
    "pattern": "(include(\\s+))",
    "type": "id_non_tilde"
  },
  {
    "name": "include",
    "remove": "include",
    "pattern": "(include(\\s*)(~))",
    "type": "id_non_tilde"
  },
  {
    "name": "&",
    "remove": "&&",
    "pattern": "(\\&\\&)",
    "type": "operation"
  },
  {
    "name": "!|",
    "remove": "!|",
    "pattern": "(\\!\\|)",
    "type": "operation"
  },
  {
    "name": "$|",
    "remove": "$|",
    "pattern": "(\\$\\|)",
    "type": "operation"
  },
  {
    "name": "!$|",
    "remove": "$!|",
    "pattern": "(\\$\\!\\|)",
    "type": "operation"
  },
  {
    "name": "!$|",
    "remove": "!$|",
    "pattern": "(\\!\\$\\|)",
    "type": "operation"
  },
  {
    "name": "!&",
    "remove": "!&",
    "pattern": "(\\!\\&)",
    "type": "operation"
  },
  {
    "name": "|",
    "remove": "||",
    "pattern": "(\\|\\|)",
    "type": "operation"
  },
  {
    "name": "!",
    "remove": "!",
    "pattern": "\\!",
    "type": "?operation",
    "comment": "?operation is just notation for the error detector to not detect this as an operation"
  },
  {
    "name": "&",
    "remove": "&",
    "pattern": "\\&",
    "type": "operation"
  },
  {
    "name": "|",
    "remove": "|",
    "pattern": "\\|",
    "type": "operation"
  },
  {
    "name": "break",
    "remove": "break",
    "pattern": "(break(\\s+))",
    "type": "id_non_tilde"
  },
  {
    "name": "skip",
    "remove": "skip",
    "pattern": "(skip(\\s+))",
    "type": "id_non_tilde"
  },
  {
    "name": "loop",
    "remove": "while",
    "pattern": "(while\\()",
    "type": "cond"
  },
  {
    "name": "number",
    "remove": "number",
    "pattern": "(number\\-\\>)",
    "type": "type"
  },
  {
    "name": "number",
    "remove": "num",
    "pattern": "(num\\-\\>)",
    "type": "type"
  },
  {
    "name": "string",
    "remove": "string",
    "pattern": "(string\\-\\>)",
    "type": "type"
  },
  {
    "name": "boolean",
    "remove": "bool",
    "pattern": "(bool\\-\\>)",
    "type": "type"
  },
  {
    "name": "falsey",
    "remove": "falsey",
    "pattern": "(falsey\\-\\>)",
    "type": "type"
  },
  {
    "name": "hash",
    "remove": "hash",
    "pattern": "(hash\\-\\>)",
    "type": "type"
  },
  {
    "name": "array",
    "remove": "array",
    "pattern": "(array\\-\\>)",
    "type": "type"
  }
]
`,
)
