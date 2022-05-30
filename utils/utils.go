package utils

import (
	"io/ioutil"
	"os"

	"github.com/laizy/web3/utils"
)

func AtomicWriteFile(filePath string, data string) {
	if FileExisted(filePath) {
		filename := filePath + "~"
		err := ioutil.WriteFile(filename, []byte(data), 0644)
		utils.Ensure(err)
		err = os.Rename(filename, filePath)
		utils.Ensure(err)
	} else {
		err := ioutil.WriteFile(filePath, []byte(data), 0644)
		utils.Ensure(err)
	}
}

// FileExisted checks whether filename exists in filesystem
func FileExisted(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
