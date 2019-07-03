package output

import (
	"encoding/xml"
	"io/ioutil"
)

type XmlOutput struct {
	OutputBase
	fileName string
	hashType string
}

type FileListAll struct {
	XMLName xml.Name `xml:"fileList"`
	Same    FileList `xml:"same"`
	Diff    FileList `xml:"different"`
}

type FileList struct {
	Files []File `xml:"file"`
}

type File struct {
	XMLName  xml.Name `xml:"file"`
	RelPath  string   `xml:"relPath,attr"`
	HashType string   `xml:"hashType,attr"`
	Hash1    string   `xml:"hash1,attr"`
	Hash2    string   `xml:"hash2,attr"`
}

func NewXmlOutput(fileName string, hashType string) *XmlOutput {
	return &XmlOutput{
		fileName: fileName,
		hashType: hashType,
	}
}

func (x *XmlOutput) Write() error {
	pairs := x.OutputBase.pairs

	//Create Structure
	fla := FileListAll{
		Same: FileList{
			Files: make([]File, 0),
		},
		Diff: FileList{
			Files: make([]File, 0),
		},
	}

	for _, p := range pairs {

		//Don't output different files
		if !p.Compare() {
			continue
		}

		//Files are same. Output
		fla.Same.Files = append(fla.Same.Files, File{
			Hash1:    p.FileOne.FileHash,
			Hash2:    p.FileTwo.FileHash,
			HashType: x.hashType,
			RelPath:  p.FileOne.GetSubPath(),
		})
	}

	for _, p := range pairs {

		//Don't output same files
		if p.Compare() {
			continue
		}

		//Files are different. Output
		fla.Diff.Files = append(fla.Diff.Files, File{
			Hash1:    p.FileOne.FileHash,
			Hash2:    p.FileTwo.FileHash,
			HashType: x.hashType,
			RelPath:  p.FileOne.GetSubPath(),
		})
	}

	//Marshal
	s, err := xml.MarshalIndent(fla, "", "  ")
	if err != nil {
		return err
	}

	//Write
	err = ioutil.WriteFile(x.fileName, s, 0644)
	if err != nil {
		return err
	}

	return nil
}
