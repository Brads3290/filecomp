package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"filecomp/output"
	"flag"
	"fmt"
	"hash"
	"os"
)

const (
	DEFAULT_HASH = "md5"
)

type SessionInfo struct {
	iptOne   string
	iptTwo   string
	hashFunc func() hash.Hash
	opt      output.Output
	threads  int
}

func main() {
	iptOne := flag.String("s1", "", "Source file or directory #1")
	iptTwo := flag.String("s2", "", "Source file or directory #2")
	hashFlag := flag.String("h", DEFAULT_HASH, "The hash algorithm to use: md5, sha1, sha256, sha512")
	xmlOutput := flag.String("o", "", "XML file to output the results to")
	threads := flag.Int("t", 1, "Number of threads to use in hash calculation")
	flag.Parse()

	//Check existance
	s1, err := os.Stat(*iptOne)
	if os.IsNotExist(err) {
		fmt.Println("ERROR: Input #1 does not exist")
		return
	}

	s2, err := os.Stat(*iptTwo)
	if os.IsNotExist(err) {
		fmt.Println("ERROR: Input #2 does not exist")
		return
	}

	//Check same type
	if s2.IsDir() != s1.IsDir() {
		fmt.Println("ERROR: Inputs must be of the same type")
		fmt.Println("Input #1 is directory: ", s1.IsDir())
		fmt.Println("Input #2 is directory: ", s2.IsDir())
		return
	}

	//Get hash func
	hf, err := GetHashFunc(*hashFlag)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	//Get the output func
	var opt output.Output
	if *xmlOutput != "" {
		opt = output.NewXmlOutput(*xmlOutput, *hashFlag)
	} else {
		opt = output.NewConsoleOutput()
	}

	//Create session info struct
	si := SessionInfo{
		iptOne:   *iptOne,
		iptTwo:   *iptTwo,
		hashFunc: hf,
		opt:      opt,
		threads:  *threads,
	}

	//Process accordingly
	if s1.IsDir() {
		err := processDirs(si)
		if err != nil {
			fmt.Println("ERROR processing directories: ", err)
		}
	} else {
		err := processFiles(si)
		if err != nil {
			fmt.Println("ERROR processing files: ", err)
		}
	}
}

func GetHashFunc(hashFlag string) (func() hash.Hash, error) {
	switch hashFlag {
	case "md5":
		return md5.New, nil
	case "sha1":
		return sha1.New, nil
	case "sha256":
		return sha256.New, nil
	case "sha512":
		return sha512.New, nil
	default:
		return nil, errors.New("invalid hash type: " + hashFlag)
	}
}
