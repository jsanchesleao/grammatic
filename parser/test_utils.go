package parser

import "grammatic/model"

var keyword_token = model.Token{Type: "TOKEN_KEYWORD", Value: "test", Line: 1, Col: 1}
var string_token = model.Token{Type: "TOKEN_STRING", Value: "\"test\"", Line: 1, Col: 1}
var int_token = model.Token{Type: "TOKEN_INT", Value: "1", Line: 1, Col: 1}
var eof_token = model.Token{Type: "TOKEN_EOF", Value: "", Line: 1, Col: 1}

var comma_token = model.Token{Type: "TOKEN_COMMA", Value: ",", Line: 1, Col: 3}

var lparen_token = model.Token{Type: "TOKEN_LPAREN", Value: "(", Line: 1, Col: 1}
var rparen_token = model.Token{Type: "TOKEN_RPAREN", Value: ")", Line: 1, Col: 5}
var bool_token = model.Token{Type: "TOKEN_BOOL", Value: "true", Line: 1, Col: 5}
