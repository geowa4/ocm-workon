package cluster

import (
	_ "embed"
	"os"
)

//go:embed zshrc
var zshrcContents string

func Initialize(baseDir string) error {
	if err := os.MkdirAll(baseDir, 0744); err != nil {
		return err
	}
	if err := makeZshrc(baseDir); err != nil {
		return err
	}
	return nil
}

func makeZshrc(baseDir string) error {
	zshrcFilePath := baseDir + pathSep + ".zshrc"
	var zshrcFile *os.File
	if _, err := os.Lstat(zshrcFilePath); os.IsNotExist(err) {
		zshrcFile, err = os.OpenFile(zshrcFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0744)
		if err != nil {
			return err
		}
		defer closeFile(zshrcFile)
	}
	if _, err := zshrcFile.WriteString(zshrcContents); err != nil {
		return err
	}
	return nil
}
