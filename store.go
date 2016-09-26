package gominin

import (
	"io"
	"math"
    "strings"
)

const (
	docHeader    = "gominin"
	InvalidDocID = DocID(math.MaxUint64)
)

type DocID uint64
type GlobalPosition uint64
type LocalPosition int

type Document interface {
	GetID() DocID
	GetBytes() []byte
	GetAttr(key string) string
	GetGlobalPosition(localPosition LocalPosition) GlobalPosition
	GetLocalPosition(globalPosition GlobalPosition) LocalPosition
}

type DocumentStore interface {
	AddDoc(in io.Reader, attrs map[string]string) Document
    AddDocString(s string, attrs map[string]string) Document
	GetDoc(id DocID) Document
	DecodeDocID(pos GlobalPosition) DocID
}

type memoryDocument struct {
	store  *memoryDocumentStore
	id     DocID
	offset GlobalPosition
	length LocalPosition
	attrs  map[string]string
}
type memoryDocumentPointers []*memoryDocument

type memoryDocumentStore struct {
	documents []byte
	docInfos  memoryDocumentPointers
}

func NewMemoryDocumentStore() DocumentStore {
	return newMemoryDocumentStore()
}

func newMemoryDocumentStore() (ds *memoryDocumentStore) {
	ds = new(memoryDocumentStore)
	ds.documents = make([]byte, 0)
	ds.documents = append(ds.documents, []byte(docHeader)...)
	ds.docInfos = make([]*memoryDocument, 0)
	return
}

func newMemoryDocument(store *memoryDocumentStore, id DocID,
	offset GlobalPosition, length int, attrs map[string]string) (doc *memoryDocument) {
	doc = new(memoryDocument)
	doc.store = store
	doc.id = id
	doc.length = LocalPosition(length)
	doc.offset = offset
	doc.attrs = attrs
	return
}

func (ds *memoryDocumentStore) AddDoc(in io.Reader, attrs map[string]string) Document {
	buf := make([]byte, 1024)
	ds.documents = append(ds.documents, '$') // Test dummy separator...

	id := DocID(len(ds.documents))
	offset := GlobalPosition(id)
	length := 0

	n, err := in.Read(buf)
	for err != io.EOF {
		length += n
		ds.documents = append(ds.documents, buf[:n]...)
		n, err = in.Read(buf)
	}

	doc := newMemoryDocument(ds, id, offset, length, attrs)
	ds.docInfos = append(ds.docInfos, doc)
	return doc
}

func (ds *memoryDocumentStore) AddDocString(s string, attrs map[string]string) Document {
    return ds.AddDoc(strings.NewReader(s), attrs)
}

func (ds *memoryDocumentStore) GetDoc(id DocID) Document {
	for _, doc := range ds.docInfos {
		if doc.id == id {
			return doc
		}
	}
	return nil
}

func (ds *memoryDocumentStore) DecodeDocID(pos GlobalPosition) DocID {
	foundIndex := LowerBound(ds.docInfos,
		func(index int) int { return int(ds.docInfos[index].offset - pos) })
	// fmt.Println("DecodeDocID pos foundIndex ", pos, foundIndex)
	if foundIndex == -1 {
		return InvalidDocID
	}

	if foundIndex == ds.docInfos.Len() { // The last document found.
		return ds.docInfos[foundIndex-1].id
	}

	if ds.docInfos[foundIndex].offset == pos {
		return ds.docInfos[foundIndex].id
	} else if foundIndex > 0 {
		return ds.docInfos[foundIndex-1].id
	} else {
		return InvalidDocID
	}
}

func (doc *memoryDocument) GetID() DocID {
	return doc.id
}

func (doc *memoryDocument) GetBytes() []byte {
	lastIndex := doc.offset + GlobalPosition(doc.length)
	return doc.store.documents[doc.offset:lastIndex]
}

func (doc *memoryDocument) GetAttr(key string) string {
	return doc.attrs[key]
}

func (doc *memoryDocument) GetGlobalPosition(localPosition LocalPosition) GlobalPosition {
	return doc.offset + GlobalPosition(localPosition)
}

func (doc *memoryDocument) GetLocalPosition(globalPosition GlobalPosition) LocalPosition {
	return LocalPosition(globalPosition - doc.offset)
}

func emptyAttrs() map[string]string {
	return map[string]string{}
}

func (md memoryDocumentPointers) Len() int {
	return len(md)
}
