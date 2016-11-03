package gominin

import (
	"errors"
	"io"
	"sort"
)

type SearchIndex interface {
	Add(in io.Reader) (Document, error)
	Search(query string) ([]DocID, error)
}

type searchIndex struct {
	store          DocumentStore
	tokenizer      Tokenizer
	termTable      TermTable
	term2positions InvertedIndex
}

type termIDPosition struct {
	id   TermID
	pos  LocalPosition
	size int
}

type termIDPositions []*termIDPosition

func NewSearchIndex() SearchIndex {
	return newSearchIndex()
}

func newSearchIndex() (si *searchIndex) {
	si = new(searchIndex)
	si.store = NewMemoryDocumentStore()
	si.tokenizer = NewCharTokenizer()
	si.termTable = NewTermTable()
	si.term2positions = NewMemoryInvertedIndex()
	return
}

func newTermIDPosition(id TermID, pos LocalPosition, size int) (tp *termIDPosition) {
	tp = new(termIDPosition)
	tp.id = id
	tp.pos = pos
	tp.size = size
	return
}

func (si *searchIndex) Add(in io.Reader) (Document, error) {
	doc := si.store.AddDoc(in, emptyAttrs())
	parsed, err := si.parse(doc.GetBytes(), true)
	if err != nil || len(parsed) == 0 {
		return nil, err
	}

	sort.Sort(parsed)

	for _, parsedPos := range parsed {
		si.term2positions.AppendPosition(parsedPos.id, doc.GetGlobalPosition(parsedPos.pos))
	}

	return doc, nil
}

func (si *searchIndex) Search(query string) ([]DocID, error) {
	var positions GlobalPositions

	parsed, err := si.parse([]byte(query), false)
	if err != nil {
		return nil, err
	}

	for _, termIDPos := range parsed {
		positions = si.searchPositions(termIDPos, positions)
	}

	if len(positions) == 0 {
		// si.termTable.dump(os.Stderr)
		// si.term2positions.dump(os.Stderr)
		return nil, nil
	}

	return si.decodeDoc(positions), nil
}

// termIDPos is query terms.
// e.g. [0, 1, 2, 5, ...]
// candPositions are candidates from term2positions.
// The positions are first searched term location.
// Then we check relative positions for each terms.
func (si *searchIndex) searchPositions(termIDPos *termIDPosition, candPositions GlobalPositions) GlobalPositions {
	positions := si.term2positions.FetchPositions(termIDPos.id)
	nextPositions := make(GlobalPositions, 0)

	if candPositions == nil {
		candPositions = make(GlobalPositions, len(positions))
		copy(candPositions, positions)
		for i := range candPositions {
			candPositions[i] -= GlobalPosition(termIDPos.pos)
		}
	}

	for _, candPos := range candPositions {
		foundIndex := BinarySearch(positions,
			func(index int) int { return int(positions[index] - (candPos + GlobalPosition(termIDPos.pos))) })
		if foundIndex > -1 {
			nextPositions = append(nextPositions, candPos)
		}
	}

	return nextPositions
}

func (si *searchIndex) parse(textBytes []byte, modify bool) (parsed termIDPositions, err error) {
	parsed = make(termIDPositions, 0)

	si.tokenizer.Init(textBytes)
	token, err := si.tokenizer.Next()
	for err == nil {
		id := si.termTable.GetID(token.Text(), modify)
		if id != NotFound {
			tp := newTermIDPosition(id, LocalPosition(token.Offset()), len(token.Text()))
			parsed = append(parsed, tp)
		} else {
			return nil, errors.New("NotFound")
		}
		token, err = si.tokenizer.Next()
	}

	return parsed, nil
}

func (si *searchIndex) decodeDoc(positions GlobalPositions) []DocID {
	var curID, prevID DocID
	var docIDs []DocID

	docIDs = make([]DocID, 0)
	prevID = InvalidDocID

	// fmt.Println("decodeDoc positions", positions)

	for _, globalPos := range positions {
		curID = si.store.DecodeDocID(globalPos)
		// fmt.Println("decodeDoc prevID, curID are ", prevID, curID)
		if prevID != curID {
			if prevID != InvalidDocID {
				docIDs = append(docIDs, prevID)
			}
			prevID = curID
		}
	}

	if prevID != InvalidDocID {
		docIDs = append(docIDs, prevID)
	}

	return docIDs
}

// Sort Interface

func (tps termIDPositions) Len() int {
	return len(tps)
}

func (tps termIDPositions) Less(i, j int) bool {
	if tps[i].id == tps[j].id {
		return tps[i].pos < tps[j].pos
	}
	return tps[i].id < tps[j].id
}

func (tps termIDPositions) Swap(i, j int) {
	tps[i], tps[j] = tps[j], tps[i]
}

// End of Sort Interface

// For BinarySearchList

func (list GlobalPositions) Len() int {
	return len(list)
}
