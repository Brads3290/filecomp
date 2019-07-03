package output

import (
	"filecomp/filelist"
)

type Output interface {
	Append(pairs ...filelist.Pair)
	Write() error
}
