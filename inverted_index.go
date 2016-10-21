package gominin

// InvertedIndex is the interface of inverted index data structure.
type InvertedIndex interface {
	FetchPositions(id TermID) []GlobalPosition
	AppendPosition(id TermID, globalPos GlobalPosition) error
}

// NewMemoryInvertedIndex creates a new instance of InvertedIndex
// which is built in memory.
func NewMemoryInvertedIndex() InvertedIndex {
	return newMemoryInvertedIndex()
}

type memoryInvertedIndex struct {
	term2positions map[TermID]globalPositions
}

func newMemoryInvertedIndex() (ii *memoryInvertedIndex) {
	ii = new(memoryInvertedIndex)
	ii.term2positions = make(map[TermID]globalPositions)
	return
}

func (ii *memoryInvertedIndex) FetchPositions(id TermID) []GlobalPosition {
	return ii.term2positions[id]
}

func (ii *memoryInvertedIndex) AppendPosition(id TermID, globalPos GlobalPosition) (err error) {
	stored := ii.term2positions[id]
	if stored == nil {
		stored = make([]GlobalPosition, 0)
	}
	ii.term2positions[id] = append(stored, globalPos)
	return
}
