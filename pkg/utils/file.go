package utils

import "os"

const PathSep = string(os.PathSeparator)

func CloseFile(f *os.File) {
	_ = f.Close()
}
