package cases

import (
	"io/fs"
	"io/ioutil"
	"log"
)

type Output struct {
	OutFile string
}

func GetDir(root string) []fs.FileInfo {
	res, err := ioutil.ReadDir(root)
	if err != nil {
		log.Printf("Error in Get dir ：%v", err)
		panic(err)
	}
	return res
}

func NewOutput(path string) Output {
	var out Output
	dir := GetDir(path)
	out.OutFile = path + "/" + dir[0].Name()

	return out
}
