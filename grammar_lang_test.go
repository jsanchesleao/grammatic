package grammatic

import (
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
Space := $EmptySpaceFormat (ignore)`

func TestJSONParsing(t *testing.T) {
	grammar := Compile(JSONGrammar)
	jsonText := `
{
  "name": "grammatic",
  "awesome": [true]
}`

	node, err := grammar.Parse("Value", jsonText)
	if err != nil {
		panic(err)
	}

	expectedSyntaxTree := `Root
  ├─Value
  │ └─Object
  │   ├─LeftBraces • {
  │   ├─ObjectBody
  │   │ ├─ObjectEntry
  │   │ │ ├─String • "name"
  │   │ │ ├─Colon • :
  │   │ │ └─Value
  │   │ │   └─String • "grammatic"
  │   │ ├─Comma • ,
  │   │ └─ObjectEntry
  │   │   ├─String • "awesome"
  │   │   ├─Colon • :
  │   │   └─Value
  │   │     └─Array
  │   │       ├─LeftBrackets • [
  │   │       ├─ArrayBody
  │   │       │ └─Value
  │   │       │   └─Bool • true
  │   │       └─RightBrackets • ]
  │   └─RightBraces • }
  └─EOF • 

`

	if expectedSyntaxTree != node.PrettyPrint() {
		t.Fatalf("Unexpected syntax tree\n%s", node.PrettyPrint())
	}

}
