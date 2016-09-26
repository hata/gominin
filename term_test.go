package gominin

import (
	"strconv"
	"testing"
)

func TestNewTermTable(t *testing.T) {
	tt := newTermTable()
	if tt.nextID != 1 {
		t.Error("nextID should be set to first termID")
	}
}

func TestNewTermTableInterface(t *testing.T) {
	tt := NewTermTable()
	id := tt.GetID("a", true)
	if id != 1 {
		t.Error("TermTable instance returns wrong id first term.")
	}
}

func TestGetID(t *testing.T) {
	tt := newTermTable()
	id := tt.GetID("a", true)
	if id != 1 {
		t.Error("First id should be 1")
	}
}

func TestGetIDFor2ndTerm(t *testing.T) {
	tt := newTermTable()
	id := tt.GetID("a", true)
	id = tt.GetID("b", true)
	if id != 2 {
		t.Error("2nd id should be 2")
	}
}

func TestGetIDForNotModify(t *testing.T) {
	tt := newTermTable()
	id := tt.GetID("a", false)
	if id != NotFound {
		t.Error("NotFound should be return when modify is false")
	}
}

func BenchmarkAddTerms(b *testing.B) {
	tt := newTermTable()
	for i := 0; i < b.N; i++ {
		tt.GetID(strconv.Itoa(i), true)
	}
}

func BenchmarkGetTermsOnly(b *testing.B) {
	tt := newTermTable()
	for i := 0; i < b.N; i++ {
		tt.GetID(strconv.Itoa(i), true)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tt.GetID(strconv.Itoa(i), false)
	}
}
