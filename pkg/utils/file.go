package utils

import "os"

const PathSep = string(os.PathSeparator)

func CloseFileAndIgnoreErrors(f *os.File) {
	_ = f.Close()
}
