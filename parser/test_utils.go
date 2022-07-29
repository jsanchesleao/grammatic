package parser

import "grammatic/model"

var keyword_token = model.Token{Type: "TOKEN_KEYWORD", Value: "test", Line: 1, Col: 1}
var int_token = model.Token{Type: "TOKEN_INT", Value: "1", Line: 1, Col: 1}
var eof_token = model.Token{Type: "TOKEN_EOF", Value: "", Line: 1, Col: 1}
