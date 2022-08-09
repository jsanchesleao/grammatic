# Grammatic

Grammatic is a parser library written in pure golang, that provides both a grammar language and a programmable API, that is powerful enough to parse common programming and data languages.

## Installation

`go get github.com/jsanchesleao/grammatic@v1.0.2` 

## Quick Start

```go
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

grammar := grammatic.Compile(JSONGrammar)

tree, err := grammar.Parse("Value", `
{
  "name": "grammatic",
  "awesome: [true],
}
`)

if err != nil {
  panic(err)
}

fmt.Println(tree.PrettyPrint())

// Root
//   ├─Value
//   │ └─Object
//   │   ├─LeftBraces • {
//   │   ├─ObjectBody
//   │   │ ├─ObjectEntry
//   │   │ │ ├─String • "name"
//   │   │ │ ├─Colon • :
//   │   │ │ └─Value
//   │   │ │   └─String • "grammatic"
//   │   │ ├─Comma • ,
//   │   │ └─ObjectEntry
//   │   │   ├─String • "awesome"
//   │   │   ├─Colon • :
//   │   │   └─Value
//   │   │     └─Array
//   │   │       ├─LeftBrackets • [
//   │   │       ├─ArrayBody
//   │   │       │ └─Value
//   │   │       │   └─Bool • true
//   │   │       └─RightBrackets • ]
//   │   └─RightBraces • }
//   └─EOF • 

```

## Production Rules

When creating a grammar, you must define Production Rules, which in the syntax is defined with the `:=` operator.

### Token Rules

The most basic production rule is a rule that matches a single token. In the syntax, it's defined by giving a rule a regex value, like so:

```
#Token Rules

Identifier := /\w+/
Space := /\s+/ (ignore)
```

If a token should be identified, but not used in subsequent rules, you can add `(ignore)` to it.

There is also a number of convenience token formats, that are provided by the library, so you don't have to rewrite commonly used regular expressions. To use a convenience token, use the `$` syntax, like this: 

```
String := $DoubleQuotedStringFormat
```

Currently there are these convenience formats:

- DigitsFormat
- IntFormat
- FloatFormat
- NumberFormat
- KeywordFormat
- DoubleQuotedStringFormat
- EmptySpaceFormat
- OperandFormat
- OpenBracesFormat
- CloseBracesFormat
- PunctuationFormat

### Repeating Rules

You can create a rule that is based on another rule, being repeatedly applied zero, one or multiple times.
The syntax in the grammar is similar to the one with regular expressions, with `*`, `+`, and `?`.

```
#this will match zero, one or more numbers
ManyNumbers := Number* 

#this will match one or more numbers, but fails when zero occurs
OneOrManyNumbers := Number+ 

#this will match one or zero numbers
MaybeNumber := Number?

Number := /\d+/
```

Just be aware that the `*` and `?` rules will always match something, even if it is an empty match, so this can lead to infinite loops, if you combine them.
In the example above, this rule would cause an infinite loop:

```
#INFINITE LOOP! DANGEROUS!
Dangerous := MaybeNumber+
```

### Repeating Rules With Separator

A common use case for repeating items is to have them separated by some other thing. For instance, the arrays in the JSON example are values separated by commas. This is an extension to the `*` and `+` rules, by adding the separator rule in square brackets:

```
# This will match things like "a", "b", "c"
ListOfStrings := String[Comma]+ 

Comma := /,/
String := $DoubleQuotedStringFormat
Space := $EmptySpaceFormat (ignore)
```

Just like with regular repeating rules, the `*` rule will always pass, as it can pass with no accepted tokens, and can also cause an infinite loop.

### Sequence Rules

A rule can be defined as a sequence of other rules. This will produce results only if ALL items in the sequence produce a value.
In the JSON example above, the following rule is a sequence: 

```
ObjectEntry := String
               Colon
               Value
```

To define a sequence, all you have to do is to write rules separated by spaces or newlines.

With sequence rules, you need to add either a rule name with no modifiers, or an inline rule, which is basically any rule followed by `as RULENAME`:

```
Array := LeftBrackets
         Value[Comma]* as ArrayBody
         RightBrackets
```

It is a syntax error if you forget the `as ArrayBody` there, because every production rule needs a well defined name.
The `[]*` syntax in the `Value` rule generates another production rule, one that matches many value separated by commas, and this new rule needs also a name.

For an inline rule, you can also have it defined within parenthesis, so it makes it more readable:

```
List := LeftParen (ListItem[Separator]* as ListBody) RightParen
```

### Or Rules

A rule can also be derived from a list of possibilities. For that you can use the `or` operator, that is defined with the `pipe` operator, like this:

```
Value := Object
       | Array
       | Number
       | String
       | Bool
```

Or rules can also use inline rules, the same way as with the Sequences.

## Tree Api

A Grammar object can be created with the `Compile` function. This grammar provides a `Parse` method, which accepts a root rule and the input string.
This method returns a tree node and an error.

This tree node holds the whole produced data that came from the defined rules.
You can actually navigate this structure and visualise it:

### PrettyPrint
	
This method will generate a string that shows the internal structure of the tree.
It will generate an output like this:

```
`Root
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
```
	
It can be very useful to see this node while navigating the tree and also to verify if the structure is being generated as expected.
	
### Fetching Nodes

The node object provides some querying methods to retrieve child nodes and values:

- GetNodeWithType(string)
  Retrieves ONE SINGLE CHILD NODE with the specified type, or nil, if none was found
  
- GetNodesWithType(string)
  Retrieves ALL CHILD NODES with the specified type. (Only goes down one level) 

- GetNodeByIndex(int)
  Retrieves the nth child node. This will return nil if an invalid index is provided.
  
- GetAllNodes()
  Retrieves ALL DIRECT CHILD NODES.
  
If a node is a leaf node from the tree, it will have the field `Token` with a non nil value.

### The Token Object

The token object will have the following properties:

- Value
  The actual chunk of the original input matched by the token defined in the grammar (the regex rule)
- Type
  The type of the token, as defined in the grammar. It should be the same as the name of the rule.
- Col and Line
  The position of the token in the original input.
  When an error occurs, this is used to point the user where the syntax error occurred.


### Recommended Processing Method

When processing the tree, usually it's a good idea to use a recursive function, that holds a switch statement by the node type, that fetches the interesting child nodes, produces some values and/or calls itself again with the child nodes, so the values are recursively generated, until a leaf node is found.

You can check the `examples` folder for usages of this method, and for more advanced parsing techniques, like virtual tokens and indentation aware parse.


	
