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
Colon := /:/
Number := $NumberFormat
Bool   := /true|false/
String := $DoubleQuotedStringFormat
Space := $EmptySpaceFormat (ignore)

`

func TestJSONParsing(t *testing.T) {
	grammar := Compile(JSONGrammar)

	node, err := grammar.Parse("Value", `
{
  "name": "grammatic",
  "awesome": [true]
}`)
	if err != nil {
		panic(err)
	}

	fmt.Println(node.PrettyPrint())
}
