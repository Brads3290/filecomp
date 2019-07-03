package filelist

import (
	"encoding/hex"
	"hash"
	"io"
	"os"
)

type File struct {
	Root     string
	Path     string
	FileHash string
}

func (f *File) calculate(h hash.Hash) (err error) {
	fo, err := os.Open(f.Path)
	if err != nil {
		return err
	}

	defer func() {
		foErr := fo.Close()

		if err == nil {
			err = foErr
		}
	}()

	if _, err := io.Copy(h, fo); err != nil {
		return err
	}

	buf := h.Sum(nil)

	f.FileHash = hex.EncodeToString(buf)

	return nil
}

func (f *File) GetSubPath() string {
	//Assume the path always begins with the root, so we can perform this
	//operation with a simple substring

	return f.Path[len(f.Root):]
}
