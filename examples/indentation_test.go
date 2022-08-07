package examples

import "testing"

func TestParseIndents(t *testing.T) {

	list := ParseIndents(`
Head
  Sub
  OtherSub
    Leaf
  YetAnother
Body
  Item
`)

	assertListEquals(t, IndentList{
		Header: "Root",
		Content: []IndentList{
			{
				Header: "Head",
				Content: []IndentList{
					{
						Header:  "Sub",
						Content: []IndentList{},
					},
					{
						Header: "OtherSub",
						Content: []IndentList{
							{
								Header:  "Leaf",
								Content: []IndentList{},
							},
						},
					},
					{
						Header:  "YetAnother",
						Content: []IndentList{},
					},
				},
			},
			{
				Header: "Body",
				Content: []IndentList{
					{
						Header:  "Item",
						Content: []IndentList{},
					},
				},
			},
		},
	}, list)

}

func assertListEquals(t *testing.T, expected, actual IndentList) {

	if expected.Header != actual.Header {
		t.Fatalf("Expected header to equal %s but it was %s", expected.Header, actual.Header)
	}

	if len(expected.Content) != len(actual.Content) {
		t.Fatalf("Expected header %s to have %d subitems, but it had %d", actual.Header, len(expected.Content), len(actual.Content))
	}

	for index := range expected.Content {
		assertListEquals(t, expected.Content[index], actual.Content[index])
	}

}
