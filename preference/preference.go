package preference

import (
	"runtime"
)

var (
	DatabasePath = "./albumin.db"
	ScanThread   = runtime.NumCPU()
)
