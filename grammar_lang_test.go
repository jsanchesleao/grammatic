package grammatic

import (
	"fmt"
	"testing"
)

const JSONGrammar = `
Value := Object
       | Array
       | Number
       | String
       | Bool

Object := LeftBraces
          ObjectEntry[Comma]* as ObjectBody
          RightBraces

ObjectEntry := String
               Colon
               Value

Array := LeftBrackets
         Value[Comma]* as ArrayBody
         RightBrackets

LeftBraces := /\{/

RightBraces := /\}/

LeftBrackets := /\[/

RightBrackets := /\]/

Comma := /,/

Number := $NumberFormat

Bool   := /true|false/

String := $DoubleQuotedStringFormat`

func TestJSONParsing(t *testing.T) {
	grammar := Compile(`
Line := LParen
        Number[Comma]+ as Items
        RParen
        (Semicolon? as LineEnd)

Number := $NumberFormat
LParen := /\(/
RParen := /\)/
Comma := /,/
Semicolon := /;/
Spaces := $EmptySpaceFormat (ignore)
`)

	fmt.Println(grammar.TokenDefs)

	node, err := grammar.Parse("Line", "( 5, 6, 7 );")
	if err != nil {
		panic(err)
	}

	fmt.Println(node.PrettyPrint())
}
