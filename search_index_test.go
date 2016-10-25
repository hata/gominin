package gominin

import (
	"strings"
	"testing"
)

func TestSearchIndex(t *testing.T) {
	si := newSearchIndex()
	doc, _ := si.Add(strings.NewReader("foo"))
	if doc == nil {
		t.Error("Add Reader failed.")
	}

	docIDs, err := si.Search("foo")

	if len(docIDs) != 1 {
		t.Error("Found document is not correct.", docIDs, err)
	}
}

func TestSearchIndexForSomeDocuments(t *testing.T) {
	si := newSearchIndex()
	doc1, _ := si.Add(strings.NewReader("foo"))
	doc2, _ := si.Add(strings.NewReader("bar"))
	doc3, _ := si.Add(strings.NewReader("foobar"))

	docIDs, err := si.Search("hoge")
	if len(docIDs) != 0 {
		t.Error("Unknown document found.")
	}

	docIDs, err = si.Search("foobar")
	if len(docIDs) != 1 || docIDs[0] != doc3.GetID() {
		t.Error("Found document should be one.", docIDs, err)
	}

	docIDs, err = si.Search("foo")
	if len(docIDs) != 2 || !(docIDs[0] == doc1.GetID() || docIDs[1] == doc1.GetID()) {
		t.Error("Found document is not correct.", docIDs, err)
	}

	docIDs, err = si.Search("bar")
	if len(docIDs) != 2 || !(docIDs[0] == doc2.GetID() || docIDs[1] == doc2.GetID()) {
		t.Error("Found document should be 2 documents.", doc1.GetID(), doc2.GetID(), doc3.GetID(), docIDs, err)
	}
}

func TestSearchNotFound(t *testing.T) {
	si := newSearchIndex()
	si.Add(strings.NewReader("foo"))
	docIDs, _ := si.Search("hoge")
	if len(docIDs) != 0 {
		t.Error("Not found should be return.")
	}
}

func TestSearchIndexUnicode(t *testing.T) {
	si := newSearchIndex()
	doc, _ := si.Add(strings.NewReader("foo\u3042\bar"))
	if doc == nil {
		t.Error("Add Reader failed.")
	}

	docIDs, err := si.Search("\u3042")

	if len(docIDs) != 1 {
		t.Error("Found document is not correct.", docIDs, err)
	}
}

func TestSearchIndexUnicodeNotFound(t *testing.T) {
	si := newSearchIndex()
	doc, _ := si.Add(strings.NewReader("foo\u3042\bar"))
	if doc == nil {
		t.Error("Add Reader failed.")
	}

	docIDs, err := si.Search("\u3043")

	if len(docIDs) != 0 {
		t.Error("Found document is not correct.", docIDs, err)
	}
}
