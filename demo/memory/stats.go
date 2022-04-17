package memory

import (
	"fmt"
	"runtime"
)

func main() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	fmt.Println(stats.HeapAlloc)
}
