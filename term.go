package gominin

import "math"

type TermID uint64

const (
	NotFound             = TermID(math.MaxUint64)
	firstID              = 1
	initialTermTableSize = 100
)

type TermTable interface {
	GetID(str string, modify bool) TermID
}

func NewTermTable() TermTable {
	return newTermTable()
}

type termTable struct {
	nextID TermID
	terms  map[string]TermID
}

func newTermTable() (tt *termTable) {
	tt = new(termTable)
	tt.nextID = firstID
	tt.terms = make(map[string]TermID)
	return
}

func (tt *termTable) GetID(str string, modify bool) TermID {
	tid, ok := tt.terms[str]

	if ok {
		return tid
	} else if modify {
		newID := tt.nextID
		tt.nextID++
		tt.terms[str] = newID
		return newID
	} else {
		return NotFound
	}
}
