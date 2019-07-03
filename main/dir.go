package main

import (
	"filecomp/filelist"
	"fmt"
	"os"
	"path/filepath"
)

func processDirs(si SessionInfo) error {

	//Walk directories and list files
	fs1, err := walk(si.iptOne)
	if err != nil {
		return err
	}

	fs2, err := walk(si.iptTwo)
	if err != nil {
		return err
	}

	//Compare to find pairs
	pairs := make([]filelist.Pair, 0)
	for _, f1 := range fs1 {
		for _, f2 := range fs2 {
			if compare(f1, f2) {
				pairs = append(pairs, filelist.NewPair(f1, f2))
				break
			}
		}
	}

	//Calculate all
	filelist.CalculateAllPairs(&pairs, si.hashFunc, si.threads)

	//Append to the output
	si.opt.Append(pairs...)

	//Create the output
	err = si.opt.Write()
	if err != nil {
		fmt.Println("Error writing output: ", err)
		return err
	}

	return nil
}

func compare(f1 filelist.File, f2 filelist.File) bool {
	return f1.GetSubPath() == f2.GetSubPath()
}

func walk(filePath string) ([]filelist.File, error) {
	fs1 := make([]filelist.File, 0)
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error occurred on path: ", path)
			fmt.Println("\t", err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		fs1 = append(fs1, filelist.File{
			Path: path,
			Root: filePath,
		})

		return nil
	})

	if err != nil {
		fmt.Println("Error walking root path: ", filePath)
		fmt.Println("\t", err)
		return nil, err
	}

	return fs1, nil
}
