package gominin

import (
    "io"
)

type debugDump interface {
    dump(w io.Writer)
}

