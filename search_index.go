package gominin

import (
	"errors"
	"io"
	"sort"
)

type SearchIndex interface {
	Add(in io.Reader)
	Search(query string)
}

type searchIndex struct {
	store          DocumentStore
	tokenizer      Tokenizer
	termTable      TermTable
	term2positions map[TermID]globalPositions
}

type termIDPosition struct {
	id  TermID
	pos LocalPosition
}

type termIDPositions []*termIDPosition
type globalPositions []GlobalPosition

func newSearchIndex() (si *searchIndex) {
	si = new(searchIndex)
	si.store = NewMemoryDocumentStore()
	si.tokenizer = NewCharTokenizer()
	si.termTable = NewTermTable()
	si.term2positions = make(map[TermID]globalPositions)
	return
}

func newTermIDPosition(id TermID, pos LocalPosition) (tp *termIDPosition) {
	tp = new(termIDPosition)
	tp.id = id
	tp.pos = pos
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
		stored := si.term2positions[parsedPos.id]
		if stored == nil {
			stored = make(globalPositions, 0)
		}
		si.term2positions[parsedPos.id] = append(stored, doc.GetGlobalPosition(parsedPos.pos))
	}

	return doc, nil
}

func (si *searchIndex) Search(query string) ([]DocID, error) {
	var positions globalPositions

	parsed, err := si.parse([]byte(query), false)
	if err != nil {
		return nil, err
	}

	for _, termIDPos := range parsed {
		// fmt.Println("Search. termIDPos:", termIDPos)
		positions = si.intersect(termIDPos, positions)
	}

	// fmt.Println("Search. positions. ", positions)

	if len(positions) == 0 {
		return nil, nil
	}

	return si.decodeDoc(positions), nil
}

func (si *searchIndex) intersect(termIDPos *termIDPosition, candPositions globalPositions) globalPositions {
	positions := si.term2positions[termIDPos.id]
	nextPositions := make(globalPositions, 0)

	if candPositions == nil {
		candPositions = make(globalPositions, len(positions))
		copy(candPositions, positions)
		for i := range candPositions {
			candPositions[i] -= GlobalPosition(termIDPos.pos)
		}
	}

	// fmt.Println("intersect. candPositions.", candPositions)

	for _, candPos := range candPositions {
		foundIndex := BinarySearch(positions,
			func(index int) int { return int(positions[index] - (candPos + GlobalPosition(termIDPos.pos))) })
		if foundIndex > -1 {
			// fmt.Println("  intersect found term in termID, localPos globalPos ", termIDPos.id, termIDPos.pos, positions[foundIndex])
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
			tp := newTermIDPosition(id, LocalPosition(token.Offset()))
			parsed = append(parsed, tp)
		} else {
			return nil, errors.New("NotFound")
		}
		token, err = si.tokenizer.Next()
	}

	return parsed, nil
}

func (si *searchIndex) decodeDoc(positions globalPositions) []DocID {
	var curID, prevID DocID
	docIDs := make([]DocID, 0)
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
	return tps[i].id < tps[j].id
}

func (tps termIDPositions) Swap(i, j int) {
	tps[i], tps[j] = tps[j], tps[i]
}

// End of Sort Interface

// For BinarySearchList

func (list globalPositions) Len() int {
	return len(list)
}
