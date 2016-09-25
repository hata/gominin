package gominin

import "io"

type Token struct {
    Text string
    Offset int
}

type Tokenizer interface {
    Init(textBytes []byte)
    Next() (*Token, error)
}

type simpleTokenizer struct {
    textBytes []byte
    pos int
    length int
}


func NewSimpleTokenizer() (tokenizer Tokenizer) {
    tokenizer = newSimpleTokenizer()
    return
}

func newToken(text string, pos int) (t *Token) {
    t = new(Token)
    t.Text = text
    t.Offset = pos
    return
}

func newSimpleTokenizer() (tokenizer *simpleTokenizer) {
    tokenizer = new(simpleTokenizer)
    return
}


func (tokenizer *simpleTokenizer) Init(textBytes []byte) {
    tokenizer.textBytes = textBytes
    tokenizer.pos = 0
    tokenizer.length = len(textBytes)
}

func (tokenizer *simpleTokenizer) Next() (t *Token, err error) {
    var lastPos int
    startPos := tokenizer.pos

    if startPos >= tokenizer.length {
        return nil, io.EOF
    }

    for lastPos = startPos + 1;lastPos < tokenizer.length;lastPos++ {
        if (tokenizer.textBytes[lastPos] & 0xC0) != 0x80 {
            break
        }
    }
    tokenizer.pos = lastPos

    return newToken(string(tokenizer.textBytes[startPos:lastPos]), startPos), nil
}

