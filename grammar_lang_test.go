package grammatic

import "testing"

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

	Compile(`
Value := Object | Array
Other := A
`)

}
