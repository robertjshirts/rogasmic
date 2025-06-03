# Assembly... Parser? Assembler? Compiler? IDK Man 

ROGASMIC
Robbie's Original Golang ASM Interpreter/Compiler

TODO
- Read file
- Write lexer types
- Write lexer tests
- Write lexer
  - Basic instructions
  - Labels, syntax, and comments
  - Split mnemonic from condition and suffixes

Lexer
- Parses instruction (with condition and suffix in the literal)
- Parses operands
  - Immediate
  - Register
  - Identifier (label)
- Parses label

Lexer stores
- Input string
- Current position
- Current line
- Current column
- Set of instructions
- Map of label to instruction #

Lex
  Loop through each character in the input string
    while token, err := l.NextToken(); token != EOF
      if token == INSTRUCTION
        

```
var suffixToBits = map[string]uint32{
	"eq": 0b0000,
	"ne": 0b0001,
	"cs": 0b0010,
	"cc": 0b0011,
	"mi": 0b0100,
	"pl": 0b0101,
	"vs": 0b0110,
	"vc": 0b0111,
	"hi": 0b1000,
	"ls": 0b1001,
	"ge": 0b1010,
	"lt": 0b1011,
	"gt": 0b1100,
	"le": 0b1101,
	"al": 0b1110,
}
```