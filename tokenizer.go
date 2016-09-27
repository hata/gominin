package gominin

import "io"

type Token interface {
	Text() string
	Offset() int
}

type Tokenizer interface {
	Init(textBytes []byte)
	Next() (Token, error)
}

type simpleToken struct {
	text   string
	offset int
}

type simpleTokenizer struct {
	textBytes []byte
	pos       int
	charPos   int
	length    int
}

func NewCharTokenizer() (tokenizer Tokenizer) {
	tokenizer = newCharTokenizer()
	return
}

func newToken(text string, pos int) Token {
	t := new(simpleToken)
	t.text = text
	t.offset = pos
	return t
}

func newCharTokenizer() (tokenizer *simpleTokenizer) {
	tokenizer = new(simpleTokenizer)
	return
}

func (tokenizer *simpleTokenizer) Init(textBytes []byte) {
	tokenizer.textBytes = textBytes
	tokenizer.pos = 0
	tokenizer.charPos = 0
	tokenizer.length = len(textBytes)
}

func (tokenizer *simpleTokenizer) Next() (t Token, err error) {
	var lastPos int
	tokenCharPos := tokenizer.charPos
	startPos := tokenizer.pos

	if startPos >= tokenizer.length {
		return nil, io.EOF
	}

	for lastPos = startPos + 1; lastPos < tokenizer.length; lastPos++ {
		if (tokenizer.textBytes[lastPos] & 0xC0) != 0x80 {
			break
		}
	}
	tokenizer.pos = lastPos
	tokenizer.charPos++

	return newToken(string(tokenizer.textBytes[startPos:lastPos]), tokenCharPos), nil
}

func (t *simpleToken) Text() string {
	return t.text
}

func (t *simpleToken) Offset() int {
	return t.offset
}
