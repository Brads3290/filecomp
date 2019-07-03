package main

import (
	"filecomp/filelist"
	"fmt"
)

func processFiles(si SessionInfo) error {

	f1 := filelist.File{
		Root: si.iptOne,
		Path: si.iptOne,
	}

	f2 := filelist.File{
		Root: si.iptTwo,
		Path: si.iptTwo,
	}

	//Just comparing two files
	p := filelist.NewPair(f1, f2)

	err := p.Calculate(si.hashFunc)
	if err != nil {
		return err
	}

	//Send to output
	si.opt.Append(p)

	//Output
	err = si.opt.Write()
	if err != nil {
		fmt.Println("Error writing to output.")
		return err
	}

	return nil
}
