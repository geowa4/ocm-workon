package cluster

import (
	"os"
)

const pathSep = string(os.PathSeparator)

func closeFile(f *os.File) {
	_ = f.Close()
}
