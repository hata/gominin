package gominin

import "testing"

func TestNewMemoryInvertedIndex(t *testing.T) {
	ii := newMemoryInvertedIndex()

	if ii == nil {
		t.Error("Failed to create new memoryInvertedIndex")
	}
}

func TestAccessPosition(t *testing.T) {
	ii := NewMemoryInvertedIndex()
	err := ii.AppendPosition(1, 2)
	if err != nil {
		t.Error("Failed to append position.", err)
	}

	posList := ii.FetchPositions(1)
	if posList[0] != 2 {
		t.Error("Stored position is wrong.")
	}
}

func TestMultiPositions(t *testing.T) {
	ii := NewMemoryInvertedIndex()
	err := ii.AppendPosition(1, 2)
	if err != nil {
		t.Error("Failed to add a new term.")
	}
	ii.AppendPosition(2, 1)
	ii.AppendPosition(1, 4)
	posList := ii.FetchPositions(1)
	if posList[0] != 2 || posList[1] != 4 {
		t.Error("")
	}
}
