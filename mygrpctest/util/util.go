package util

import (
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
)

type TestCase struct {
	Name       string
	Expectfile string
	Initfile   string
	Wavfile    string
}

func GetDir(root string) []fs.FileInfo {
	res, err := ioutil.ReadDir(root)
	if err != nil {
		log.Printf("Error in Get dir ：%v", err)
		panic(err)
	}
	return res
}

func NewTestCase(path string) TestCase {
	var cases TestCase
	dir := GetDir(path)
	cases.Name = path
	cases.Expectfile = path + "/" + dir[0].Name()
	cases.Initfile = path + "/" + dir[1].Name()
	cases.Wavfile = path + "/" + dir[2].Name()

	return cases
}

func GetAllFiles(dirPth string) (files []string, err error) {
	fis, err := ioutil.ReadDir(filepath.Clean(filepath.ToSlash(dirPth)))
	if err != nil {
		return nil, err
	}

	for _, f := range fis {
		_path := filepath.Join(dirPth, f.Name())

		if f.IsDir() {
			fs, _ := GetAllFiles(_path)
			files = append(files, fs...)
			continue
		}

		// 指定格式
		switch filepath.Ext(f.Name()) {
		case ".txt", ".json", ".wav":
			files = append(files, _path)
		}
	}

	return files, nil
}
