package gominin

import (
	"io"
	"testing"
)

func TestNewCharTokenizer(t *testing.T) {
	tokenizer := newCharTokenizer()
	tokenizer.Init([]byte("foo"))
	if string(tokenizer.textBytes) != "foo" {
		t.Error("Init should initialize tokenizer fields.")
	}
	if tokenizer.pos != 0 {
		t.Error("Init should initialize pos field.")
	}
	if tokenizer.length != len("foo") {
		t.Error("Init should initialize length.")
	}
}

func TestNext(t *testing.T) {
	tokenizer := newCharTokenizer()
	tokenizer.Init([]byte("foo"))
	token, err := tokenizer.Next()
	if token.Text() != "f" {
		t.Error("Failed to return a token.")
	}
	if token.Offset() != 0 {
		t.Error("Failed to return an offset.")
	}
	if err != nil {
		t.Error("Failed to return no error.")
	}
}

func TestSomeNext(t *testing.T) {
	tokenizer := newCharTokenizer()
	tokenizer.Init([]byte("foo"))
	token, err := tokenizer.Next()
	var str string

	for err != io.EOF {
		str += token.Text()
		token, err = tokenizer.Next()
	}
	if str != "foo" {
		t.Error("Failed to return all tokens.")
	}
}

func TestNextOffset(t *testing.T) {
	tokenizer := newCharTokenizer()
	tokenizer.Init([]byte("foo"))
	for i := 0; i < len("foo"); i++ {
		token, _ := tokenizer.Next()
		if token.Offset() != i {
			t.Error("Offset should return a correct position.")
		}
	}
}

func TestNextForMultiBytes(t *testing.T) {
	tokenizer := newCharTokenizer()
	tokenizer.Init([]byte("f\u3042\u3044"))
	token, err := tokenizer.Next()
	if token.Text() != "f" {
		t.Error("Check a single byte returns correctly.")
	}
	token, err = tokenizer.Next()
	if token.Text() != "\u3042" {
		t.Error("Multibyte does not return correctly.")
	}
	if token.Offset() != 1 {
		t.Error("1st multibyte char pos should be 1.")
	}
	token, err = tokenizer.Next()
	if token.Text() != "\u3044" {
		t.Error("Multibyte should be able to handle.")
	}
	if token.Offset() != 2 {
		t.Error("Token offset should be the location of char instead of bytes.")
	}
	token, err = tokenizer.Next()
	if err != io.EOF {
		t.Error("Last error should be io.EOF")
	}
}
