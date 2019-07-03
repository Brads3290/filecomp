package output

import "filecomp/filelist"

type OutputBase struct {
	pairs []filelist.Pair
}

func (ob *OutputBase) Append(pairs ...filelist.Pair) {
	ob.pairs = append(ob.pairs, pairs...)
}
