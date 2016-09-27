package gominin

import (
	"strings"
	"testing"
)

func TestNewMemoryDocumentStore(t *testing.T) {
	ms := newMemoryDocumentStore()
	if len(ms.documents) != len(docHeader) {
		t.Error("Initialize documents.")
	}
	if len(ms.docInfos) != 0 {
		t.Error("Initialize document Information table")
	}
}

func TestNewMemoryDocument(t *testing.T) {
	doc := newMemoryDocument(nil, 1, 2, 3, map[string]string{"foo": "bar"})
	if doc.id != 1 {
		t.Error("memoryDocument.id should be initialized.")
	}
	if doc.offset != 2 {
		t.Error("memoryDocument.offset should be initialized.")
	}
	if doc.length != 3 {
		t.Error("memoryDocument.length should be initialized.")
	}
	if doc.attrs["foo"] != "bar" {
		t.Error("memoryDocument.attrs should be initialized.")
	}
}

func TestAddDoc(t *testing.T) {
	ms := newMemoryDocumentStore()
	doc := ms.AddDoc(strings.NewReader("foo"), map[string]string{})
	if doc.GetID() == 0 {
		t.Error("Doc should be created by AddDoc.")
	}
	if string(doc.GetBytes()) != "foo" {
		t.Error("stored document is not retrieved.")
	}
}

func TestAddDocString(t *testing.T) {
	ms := newMemoryDocumentStore()
	doc := ms.AddDocString("foo", map[string]string{})
	if doc.GetID() == 0 {
		t.Error("Doc should be created by AddDoc.")
	}
	if string(doc.GetBytes()) != "foo" {
		t.Error("stored document is not retrieved.")
	}
}

func TestGetDoc(t *testing.T) {
	ms := newMemoryDocumentStore()
	doc := ms.AddDoc(strings.NewReader("foo"), map[string]string{})
	id := doc.GetID()
	doc = ms.GetDoc(doc.GetID())
	if doc.GetID() != id {
		t.Error("Failed to get a document because of id.")
	}
	if string(doc.GetBytes()) != "foo" {
		t.Error("Failed to get a document because of text.")
	}
}

func TestGetDocNotFound(t *testing.T) {
	ms := newMemoryDocumentStore()
	doc := ms.GetDoc(0)
	if doc != nil {
		t.Error("Unknown document is found.")
	}
}

func TestGetDocNotFoundAfterAddDoc(t *testing.T) {
	ms := newMemoryDocumentStore()
	doc := ms.AddDoc(strings.NewReader("foo"), map[string]string{})
	doc = ms.GetDoc(0)
	if doc != nil {
		t.Error("Unknown document is found.")
	}
}

func TestSomeAddDoc(t *testing.T) {
	ms := newMemoryDocumentStore()
	doc1 := ms.AddDoc(strings.NewReader("foo"), emptyAttrs())
	doc2 := ms.AddDoc(strings.NewReader("bar"), emptyAttrs())
	doc3 := ms.AddDoc(strings.NewReader("hoge"), emptyAttrs())
	doc := ms.GetDoc(doc1.GetID())
	if doc.GetID() != doc1.GetID() || string(doc.GetBytes()) != "foo" {
		t.Error("Wrong document is returned.")
	}
	doc = ms.GetDoc(doc2.GetID())
	if doc.GetID() != doc2.GetID() || string(doc.GetBytes()) != "bar" {
		t.Error("Wrong document is returned.")
	}
	doc = ms.GetDoc(doc3.GetID())
	if doc.GetID() != doc3.GetID() || string(doc.GetBytes()) != "hoge" {
		t.Error("Wrong document is returned.")
	}
}
