package examples

import "testing"

func assertUrlEquals(t *testing.T, expected, actual ParsedUrl) {
	t.Helper()

	if expected.Protocol != actual.Protocol {
		t.Fatalf("Expected protocol to be %s but it was %s", expected.Protocol, actual.Protocol)
	}
	if expected.Domain != actual.Domain {
		t.Fatalf("Expected domain to be %s but it was %s", expected.Domain, actual.Domain)
	}
	if expected.Path != actual.Path {
		t.Fatalf("Expected path to be %s but it was %s", expected.Path, actual.Path)
	}
	if expected.Anchor != actual.Anchor {
		t.Fatalf("Expected anchor to be %s but it was %s", expected.Anchor, actual.Anchor)
	}
	assertUrlQueryString(t, expected.Query, actual.Query)
	assertUrlAuth(t, expected.Auth, actual.Auth)
}

func assertUrlQueryString(t *testing.T, expected, actual []QueryParam) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Fatalf("Expected query string to contain %d items, but it had %d", len(expected), len(actual))
	}
	for index := range expected {
		if expected[index].Key != actual[index].Key {
			t.Fatalf("Expected key to be %s but it was %s", expected[index].Key, actual[index].Key)
		}
		if expected[index].Value != actual[index].Value {
			t.Fatalf("Expected value to be %s but it was %s", expected[index].Value, actual[index].Value)
		}
	}
}

func assertUrlAuth(t *testing.T, expected, actual *Auth) {
	t.Helper()
	if expected == nil && actual == nil {
		return
	}
	if expected == nil && actual != nil {
		t.Fatalf("Expected auth to be nil, but it was %v", actual)
	}
	if expected != nil && actual == nil {
		t.Fatalf("Expected auth to be %v, but it was nil", expected)
	}
	if expected.Username != actual.Username {
		t.Fatalf("Expected username to be %s but it was %s", expected.Username, expected.Password)
	}
	if expected.Password != actual.Password {
		t.Fatalf("Expected password to be %s but it was %s", expected.Password, expected.Password)
	}

}

func TestParseUrl(t *testing.T) {

	assertUrlEquals(t, ParsedUrl{
		Protocol: "https",
		Auth:     nil,
		Domain:   "github.com",
		Path:     "/jsanchesleao/grammatic",
		Query:    []QueryParam{},
		Anchor:   "",
	}, UrlParse("https://github.com/jsanchesleao/grammatic"))

	assertUrlEquals(t, ParsedUrl{
		Protocol: "https",
		Auth: &Auth{
			Username: "admin",
			Password: "passwd",
		},
		Domain: "some.domain.com",
		Path:   "/somepath",
		Query: []QueryParam{
			{
				Key:   "foo",
				Value: "1",
			},
			{
				Key:   "bar",
				Value: "2",
			},
		},
		Anchor: "link",
	}, UrlParse("https://admin:passwd@some.domain.com/somepath?foo=1&bar=2#link"))
}
