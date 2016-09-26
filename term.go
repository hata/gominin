package gominin

type TermID uint64

const (
	NotFound             = 0
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
	var newId TermID
	tid, ok := tt.terms[str]

	if ok {
		return tid
	} else if modify {
		newId = tt.nextID
		tt.nextID++
		tt.terms[str] = newId
		return newId
	} else {
		return NotFound
	}
}
