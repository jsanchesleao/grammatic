package examples

import (
	"grammatic"
	"grammatic/model"
)

const urlGrammar = `
Url := Protocol
       AfterProtocol
       Credentials? as MaybeCredentials
       Domain
       Path
       QueryString? as MaybeQueryString
       HashAnchor? as MaybeHashAnchor

Credentials := Word as AuthUser
               Colon
               Word as AuthPassword
               At

Domain := Word[Dot]+
Path := PathPart*

QueryString := QuestionMark
               QueryParam[Ampersand]* as QueryStringItems

QueryParam := Word as ParamKey
              Equals
              Word as ParamValue

HashAnchor := Hash Word

Protocol := /https?/
PathPart := /\/[^\/\?]+/

Word := /[a-zA-Z0-9_-]+/

QuestionMark := /\?/
Ampersand := /&/
Hash := /#/
Equals := /=/
Dot := /\./
At := /@/

AfterProtocol := /:\/\//
Colon := /:/
`

type ParsedUrl struct {
	Protocol string
	Auth     *Auth
	Domain   string
	Path     string
	Query    []QueryParam
	Anchor   string
}

type Auth struct {
	Username string
	Password string
}

type QueryParam struct {
	Key   string
	Value string
}

func UrlParse(input string) ParsedUrl {
	grammar := grammatic.Compile(urlGrammar)
	tree, err := grammar.Parse("Url", input)

	if err != nil {
		panic(err)
	}

	return reduceUrlTree(tree)
}

func reduceUrlTree(node *model.Node) ParsedUrl {
	url := ParsedUrl{}
	node = node.GetNodeWithType("Url")

	url.Protocol = node.GetNodeWithType("Protocol").Token.Value
	url.Domain = reduceChildNodesToString(node.GetNodeWithType("Domain"))
	url.Path = reduceChildNodesToString(node.GetNodeWithType("Path"))
	url.Query = []QueryParam{}

	anchor := node.GetNodeWithType("MaybeHashAnchor").GetNodeWithType("HashAnchor")
	if anchor != nil {
		url.Anchor = anchor.GetNodeWithType("Word").Token.Value
	}

	queryString := node.GetNodeWithType("MaybeQueryString").GetNodeWithType("QueryString")
	if queryString != nil {
		params := queryString.GetNodeWithType("QueryStringItems").GetNodesWithType("QueryParam")
		for _, queryParam := range params {
			queryStringItem := QueryParam{
				Key:   queryParam.GetNodeWithType("ParamKey").Token.Value,
				Value: queryParam.GetNodeWithType("ParamValue").Token.Value,
			}
			url.Query = append(url.Query, queryStringItem)
		}
	}

	auth := node.GetNodeWithType("MaybeCredentials").GetNodeWithType("Credentials")
	if auth != nil {
		authValue := Auth{
			Username: auth.GetNodeWithType("AuthUser").Token.Value,
			Password: auth.GetNodeWithType("AuthPassword").Token.Value,
		}
		url.Auth = &authValue
	}

	return url
}

func reduceChildNodesToString(node *model.Node) string {
	domain := ""
	for _, n := range node.GetAllNodes() {
		domain = domain + n.Token.Value
	}
	return domain
}
