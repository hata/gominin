package gominin

import (
	"fmt"
	"io"
)

// InvertedIndex is the interface of inverted index data structure.
type InvertedIndex interface {
	debugDump
	FetchPositions(id TermID) GlobalPositions
	AppendPosition(id TermID, globalPos GlobalPosition) error
}

// NewMemoryInvertedIndex creates a new instance of InvertedIndex
// which is built in memory.
func NewMemoryInvertedIndex() InvertedIndex {
	return newMemoryInvertedIndex()
}

type memoryInvertedIndex struct {
	term2positions map[TermID]GlobalPositions
}

func newMemoryInvertedIndex() (ii *memoryInvertedIndex) {
	ii = new(memoryInvertedIndex)
	ii.term2positions = make(map[TermID]GlobalPositions)
	return
}

func (ii *memoryInvertedIndex) FetchPositions(id TermID) GlobalPositions {
	return ii.term2positions[id]
}

func (ii *memoryInvertedIndex) AppendPosition(id TermID, globalPos GlobalPosition) (err error) {
	stored := ii.term2positions[id]
	if stored == nil {
		stored = make(GlobalPositions, 0)
	}
	ii.term2positions[id] = append(stored, globalPos)
	return
}

func (ii *memoryInvertedIndex) dump(w io.Writer) {
    for k, v := range ii.term2positions {
        fmt.Fprintln(w, k, v)
    }
}
